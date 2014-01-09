package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"
)

const maxRecursion uint8 = 5

type GoWsdl struct {
	file, pkg             string
	wsdl                  *Wsdl
	resolvedXsdExternals  map[string]bool
	currentRecursionLevel uint8
}

var cacheDir = os.TempDir() + "gowsdl-cache"

func init() {
	err := os.MkdirAll(cacheDir, 0700)
	if err != nil {
		log.Fatalf("Unable to reate cache directory")
	}
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func downloadFile(url string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Dial: dialTimeout,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewGoWsdl(file, pkg string) (*GoWsdl, error) {
	file = strings.TrimSpace(file)
	if file == "" {
		log.Fatalln("WSDL file is required to generate Go proxy")
	}

	pkg = strings.TrimSpace(pkg)
	if pkg == "" {
		pkg = "main"
	}

	return &GoWsdl{
		file: file,
		pkg:  pkg,
	}, nil
}

func (g *GoWsdl) Start() (map[string][]byte, error) {
	gocode := make(map[string][]byte)

	err := g.unmarshal()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["types"], err = g.genTypes()
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["messages"], err = g.genMessages()
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["operations"], err = g.genOperations()
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error

		gocode["proxy"], err = g.genSoapProxy()
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return gocode, nil
}

func (g *GoWsdl) unmarshal() error {
	var data []byte

	parsedUrl, err := url.Parse(g.file)
	if parsedUrl.Scheme == "" {
		log.Printf("Reading file %s...\n", g.file)

		data, err = ioutil.ReadFile(g.file)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Downloading %s...\n", g.file)

		data, err = downloadFile(g.file)
		if err != nil {
			return err
		}
	}

	g.wsdl = &Wsdl{}
	err = xml.Unmarshal(data, g.wsdl)
	if err != nil {
		return err
	}

	for _, schema := range g.wsdl.Types.Schemas {
		err = g.resolveXsdExternals(schema, parsedUrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GoWsdl) resolveXsdExternals(schema *XsdSchema, url *url.URL) error {
	for _, incl := range schema.Includes {
		location, err := url.Parse(incl.SchemaLocation)
		if err != nil {
			return err
		}

		_, schemaName := filepath.Split(location.Path)
		if g.resolvedXsdExternals[schemaName] {
			continue
		}

		schemaLocation := location.String()
		if !location.IsAbs() {
			if !url.IsAbs() {
				return errors.New(fmt.Sprintf("Unable to resolve external schema %s through WSDL URL %s", schemaLocation, url))
			}
			schemaLocation = url.Scheme + "://" + url.Host + schemaLocation
		}

		log.Printf("Downloading external schema: %s\n", schemaLocation)

		data, err := downloadFile(schemaLocation)
		newschema := &XsdSchema{}

		err = xml.Unmarshal(data, newschema)
		if err != nil {
			return err
		}

		if len(newschema.Includes) > 0 &&
			maxRecursion > g.currentRecursionLevel {

			g.currentRecursionLevel++

			//log.Printf("Entering recursion %d\n", g.currentRecursionLevel)
			g.resolveXsdExternals(newschema, url)
		}

		g.wsdl.Types.Schemas = append(g.wsdl.Types.Schemas, newschema)

		if g.resolvedXsdExternals == nil {
			g.resolvedXsdExternals = make(map[string]bool, maxRecursion)
		}
		g.resolvedXsdExternals[schemaName] = true
	}

	return nil
}

func (g *GoWsdl) genTypes() ([]byte, error) {
	funcMap := template.FuncMap{
		"toGoType":             toGoType,
		"stripns":              stripns,
		"replaceReservedWords": replaceReservedWords,
		"makeFieldPublic":      makeFieldPublic,
	}

	//TODO resolve element refs in place.
	//g.resolveElementsRefs()

	data := new(bytes.Buffer)
	tmpl := template.Must(template.New("types").Funcs(funcMap).Parse(typesTmpl))
	err := tmpl.Execute(data, g.wsdl.Types)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

// func (g *GoWsdl) resolveElementsRefs() error {
// 	for _, schema := range g.wsdl.Types.Schemas {
// 		for _, globalEl := range schema.Elements {
// 			for _, localEl := range globalEl.ComplexType.Sequence.Elements {

// 			}
// 		}
// 	}
// }

func (g *GoWsdl) genOperations() ([]byte, error) {
	funcMap := template.FuncMap{
		"toGoType":             toGoType,
		"stripns":              stripns,
		"replaceReservedWords": replaceReservedWords,
		"makeFieldPublic":      makeFieldPublic,
		"findType":             g.findType,
	}

	data := new(bytes.Buffer)
	tmpl := template.Must(template.New("operations").Funcs(funcMap).Parse(opsTmpl))
	err := tmpl.Execute(data, g.wsdl.PortTypes)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

var reservedWords = map[string]string{
	"break":       "break_",
	"default":     "default_",
	"func":        "func_",
	"interface":   "interface_",
	"select":      "select_",
	"case":        "case_",
	"defer":       "defer_",
	"go":          "go_",
	"map":         "map_",
	"struct":      "struct_",
	"chan":        "chan_",
	"else":        "else_",
	"goto":        "goto_",
	"package":     "package_",
	"switch":      "switch_",
	"const":       "const_",
	"fallthrough": "fallthrough_",
	"if":          "if_",
	"range":       "range_",
	"type":        "type_",
	"continue":    "continue_",
	"for":         "for_",
	"import":      "import_",
	"return":      "return_",
	"var":         "var_",
}

func replaceReservedWords(identifier string) string {
	value := reservedWords[identifier]
	if value != "" {
		return value
	}
	return identifier
}

var xsd2GoTypes = map[string]string{
	"string":        "string",
	"token":         "string",
	"float":         "float32",
	"double":        "float64",
	"decimal":       "float64",
	"integer":       "int32",
	"int":           "int32",
	"short":         "int16",
	"byte":          "int8",
	"long":          "int64",
	"boolean":       "bool",
	"dateTime":      "time.Time",
	"date":          "time.Time",
	"time":          "time.Time",
	"base64Binary":  "[]byte",
	"hexBinary":     "[]byte",
	"unsignedInt":   "uint32",
	"unsignedShort": "uint16",
	"unsignedByte":  "byte",
	"unsignedLong":  "uint64",
	"anyType":       "interface{}",
}

func toGoType(xsdType string) string {
	//Handles name space, ie. xsd:string, xs:string
	r := strings.Split(xsdType, ":")

	type_ := r[0]

	if len(r) == 2 {
		type_ = r[1]
	}

	value := xsd2GoTypes[type_]

	if value != "" {
		return value
	}

	return "*" + type_
}

//I'm not very proud of this function but
//it works for now and performance doesn't
//seem critical at this point
func (g *GoWsdl) findType(message string) string {
	message = stripns(message)
	for _, msg := range g.wsdl.Messages {
		if msg.Name != message {
			continue
		}

		//Assumes document/literal wrapped WS-I
		part := msg.Parts[0]
		if part.Type != "" {
			return stripns(part.Type)
		}

		elRef := stripns(part.Element)
		for _, schema := range g.wsdl.Types.Schemas {
			for _, el := range schema.Elements {
				if elRef == el.Name {
					if el.Type != "" {
						return stripns(el.Type)
					}
					return el.Name
				}
			}
		}
	}
	return ""
}

//TODO: Add namespace support instead of stripping it
func stripns(xsdType string) string {
	r := strings.Split(xsdType, ":")
	type_ := r[0]

	if len(r) == 2 {
		type_ = r[1]
	}

	return type_
}

func makeFieldPublic(field_ string) string {
	field := []rune(field_)
	if len(field) == 0 {
		return field_
	}

	field[0] = unicode.ToUpper(field[0])
	return string(field)
}

func (g *GoWsdl) genMessages() ([]byte, error) {
	return nil, nil
}

func (g *GoWsdl) genSoapProxy() ([]byte, error) {
	return nil, nil
}
