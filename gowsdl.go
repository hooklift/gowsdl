package main

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func downloadFile(url string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
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

	totalAdts := 0
	for _, schema := range g.wsdl.Types.Schemas {
		for _, el := range schema.Elements {
			if el.Type == "" {
				log.Printf("%s -> %#v\n\n", strings.TrimSuffix(el.ComplexType.Name, "Type"), el.ComplexType.Sequence)
				totalAdts++
			}
		}

		for _, complexType := range schema.ComplexTypes {
			log.Printf("%s -> %#v\n\n", strings.TrimSuffix(complexType.Name, "Type"), complexType.Sequence)
			totalAdts++
		}
	}

	log.Printf("Abstract data types: %d\n", totalAdts)
	log.Printf("Total schemas: %#d\n\n", len(g.wsdl.Types.Schemas))

	return nil, nil
}

func (g *GoWsdl) genOperations() ([]byte, error) {
	for _, pt := range g.wsdl.PortTypes {
		// for _, o := range pt.Operations {
		// 	g.logger.Printf("Operation: %s", o.Name)
		// }
		log.Printf("Total ops: %d\n", len(pt.Operations))
	}

	return nil, nil
}

func (g *GoWsdl) genMessages() ([]byte, error) {
	return nil, nil
}

func (g *GoWsdl) genSoapProxy() ([]byte, error) {
	return nil, nil
}
