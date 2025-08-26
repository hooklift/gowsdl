// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestElementGenerationDoesntCommentOutStructProperty(t *testing.T) {
	g, err := NewGoWSDL("fixtures/test.wsdl", "myservice", false, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Error(err)
	}

	if strings.Contains(string(resp["types"]), "// this is a comment  GetInfoResult string `xml:\"GetInfoResult,omitempty\"`") {
		t.Error("Type comment should not comment out struct type property")
		t.Error(string(resp["types"]))
	}
}

func TestComplexTypeWithInlineSimpleType(t *testing.T) {
	g, err := NewGoWSDL("fixtures/test.wsdl", "myservice", false, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Fatal(err)
	}
	actual, err := getTypeDeclaration(resp, "GetInfo")
	if err != nil {
		t.Fatal(err)
	}

	expected := `type GetInfo struct {
	XMLName	xml.Name	` + "`" + `xml:"http://www.mnb.hu/webservices/ GetInfo"` + "`" + `

	Id	string	` + "`" + `xml:"Id,omitempty" json:"Id,omitempty"` + "`" + `
}`
	if actual != expected {
		t.Error("got " + actual + " want " + expected)
	}
}

func TestAttributeRef(t *testing.T) {
	g, err := NewGoWSDL("fixtures/test.wsdl", "myservice", false, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Fatal(err)
	}
	actual, err := getTypeDeclaration(resp, "ResponseStatus")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected := `type ResponseStatus struct {
	Status	[]struct {
		Value	string  ` + "`" + `xml:",chardata" json:"-,"` + "`" + `

		Code	string	` + "`" + `xml:"code,attr,omitempty" json:"code,omitempty"` + "`" + `
	}	` + "`" + `xml:"status,omitempty" json:"status,omitempty"` + "`" + `

	ResponseCode	string	` + "`" + `xml:"http://www.mnb.hu/webservices/ responseCode,attr,omitempty" json:"responseCode,omitempty"` + "`" + `
}`
	actual = string(bytes.ReplaceAll([]byte(actual), []byte("\t"), []byte("  ")))
	expected = string(bytes.ReplaceAll([]byte(expected), []byte("\t"), []byte("  ")))
	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}
}

func TestElementWithLocalSimpleType(t *testing.T) {
	g, err := NewGoWSDL("fixtures/test.wsdl", "myservice", false, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Type declaration
	actual, err := getTypeDeclaration(resp, "ElementWithLocalSimpleType")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected := `type ElementWithLocalSimpleType string`

	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}

	// Const declaration of first enum value
	actual, err = getTypeDeclaration(resp, "ElementWithLocalSimpleTypeEnum1")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected = `const ElementWithLocalSimpleTypeEnum1 ElementWithLocalSimpleType = "enum1"`

	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}

	// Const declaration of second enum value
	actual, err = getTypeDeclaration(resp, "ElementWithLocalSimpleTypeEnum2")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected = `const ElementWithLocalSimpleTypeEnum2 ElementWithLocalSimpleType = "enum2"`

	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}
}

func TestDateTimeType(t *testing.T) {
	g, err := NewGoWSDL("fixtures/test.wsdl", "myservice", false, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Type declaration
	actual, err := getTypeDeclaration(resp, "StartDate")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected := `type StartDate soap.XSDDateTime`

	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}

	// Method declaration MarshalXML
	actual, err = getFuncDeclaration(resp, "MarshalXML", "StartDate")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected = `func (xdt StartDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return soap.XSDDateTime(xdt).MarshalXML(e, start)
}`

	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}

	// Method declaration UnmarshalXML
	actual, err = getFuncDeclaration(resp, "UnmarshalXML", "StartDate")
	if err != nil {
		fmt.Println(string(resp["types"]))
		t.Fatal(err)
	}

	expected = `func (xdt *StartDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return (*soap.XSDDateTime)(xdt).UnmarshalXML(d, start)
}`

	if actual != expected {
		t.Error("got \n" + actual + " want \n" + expected)
	}
}

func TestVboxGeneratesWithoutSyntaxErrors(t *testing.T) {
	files, err := filepath.Glob("fixtures/*.wsdl")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		if file == "fixtures/vim.wsdl" {
			// This fixture references missing .xsd files from the VMWare SDK, which aren't actually
			// vendored. We use it explicitly from test_wsdl.go to test unmarshalling, but not here.
			continue
		}
		t.Run(file, func(t *testing.T) {
			g, err := NewGoWSDL(file, "myservice", false, true)
			if err != nil {
				t.Error(err)
			}

			resp, err := g.Start()
			if err != nil {
				t.Error(err)
			}

			data := new(bytes.Buffer)
			data.Write(resp["header"])
			data.Write(resp["types"])
			data.Write(resp["operations"])
			data.Write(resp["soap"])

			_, err = format.Source(data.Bytes())
			if err != nil {
				fmt.Println(data.String())
				t.Error(err)
			}
		})
	}
}

func TestEnumerationsGeneratedCorrectly(t *testing.T) {
	enumStringTest := func(t *testing.T, fixtureWsdl string, varName string, typeName string, enumString string) {
		g, err := NewGoWSDL("fixtures/"+fixtureWsdl, "myservice", false, true)
		if err != nil {
			t.Error(err)
		}

		resp, err := g.Start()
		if err != nil {
			t.Error(err)
		}
		re := regexp.MustCompile(varName + " " + typeName + " = \"([^\"]*)\"")
		matches := re.FindStringSubmatch(string(resp["types"]))

		if len(matches) != 2 {
			t.Errorf("No match or too many matches found for %s", varName)
		} else if matches[1] != enumString {
			t.Errorf("%s got '%s' but expected '%s'", varName, matches[1], enumString)
		}
	}
	enumStringTest(t, "chromedata.wsdl", "DriveTrainFrontWheelDrive", "DriveTrain", "Front Wheel Drive")
	enumStringTest(t, "chromedata.wsdl", "DriveTrainEmptyString", "DriveTrain", "")
	enumStringTest(t, "vboxweb.wsdl", "SettingsVersionV1_14", "SettingsVersion", "v1_14")

}

func TestComplexTypeGeneratedCorrectly(t *testing.T) {
	g, err := NewGoWSDL("fixtures/workday-time-min.wsdl", "myservice", false, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Error(err)
	}

	decl, err := getTypeDeclaration(resp, "WorkerObjectIDType")

	expected := "type WorkerObjectIDType struct"
	re := regexp.MustCompile(expected)
	matches := re.FindStringSubmatch(decl)

	if len(matches) != 1 {
		t.Errorf("No match or too many matches found for WorkerObjectIDType")
	} else if matches[0] != expected {
		t.Errorf("WorkerObjectIDType got '%s' but expected '%s'", matches[1], expected)
	}
}

func TestEPCISWSDL(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	g, err := NewGoWSDL("./fixtures/epcis/EPCglobal-epcis-query-1_2.wsdl", "myservice", true, true)
	if err != nil {
		t.Error(err)
	}

	resp, err := g.Start()
	if err != nil {
		t.Fatal(err)
	}
	data := new(bytes.Buffer)
	data.Write(resp["header"])
	data.Write(resp["types"])
	data.Write(resp["operations"])
	data.Write(resp["soap"])

	// go fmt the generated code
	source, err := format.Source(data.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	expectedBytes, err := os.ReadFile("./fixtures/epcis/epcisquery.src")
	if err != nil {
		t.Fatal(err)
	}

	actual := string(source)
	expected := string(expectedBytes)
	if actual != expected {
		_ = os.WriteFile("./fixtures/epcis/epcisquery_gen.src", source, 0664)
		t.Error("got source ./fixtures/epcis/epcisquery_gen.src but expected ./fixtures/epcis/epcisquery.src")
	}
}

func getTypeDeclaration(resp map[string][]byte, name string) (string, error) {
	source, err := format.Source([]byte(string(resp["header"]) + string(resp["types"])))
	if err != nil {
		return "", err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "myservice.go", string(source), parser.DeclarationErrors)
	if err != nil {
		return "", err
	}
	o := f.Scope.Lookup(name)
	if o == nil {
		return "", errors.New("type " + name + " is missing")
	}
	var buf bytes.Buffer
	buf.WriteString(o.Kind.String())
	buf.WriteString(" ")
	err = printer.Fprint(&buf, fset, o.Decl)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func findFuncDecl(f *ast.File, name string, recv string) *ast.Decl {
	// Loop over all declarations
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Found FuncDecl declaration type
			if funcDecl.Name.Name == name {
				// Found match with function name
				if recvType, ok := funcDecl.Recv.List[0].Type.(*ast.Ident); ok {
					// Value receiver type
					if recvType.Name == recv {
						// Found receiver type match
						return &decl
					}
				} else if recvTypePtr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					// Pointer receiver type
					if recvType, ok := recvTypePtr.X.(*ast.Ident); ok {
						if recvType.Name == recv {
							// Found receiver type match
							return &decl
						}
					}
				}
			}
		}
	}

	return nil
}

func getFuncDeclaration(resp map[string][]byte, name string, recv string) (string, error) {
	source, err := format.Source([]byte(string(resp["header"]) + string(resp["types"])))
	if err != nil {
		return "", err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "myservice.go", string(source), parser.DeclarationErrors)
	if err != nil {
		return "", err
	}

	decl := findFuncDecl(f, name, recv)
	if decl == nil {
		return "", errors.New("Function declaration " + name + " not found")
	}
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, *decl)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
