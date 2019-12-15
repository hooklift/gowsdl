// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
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

	ResponseCode	string	` + "`" + `xml:"responseCode,attr,omitempty" json:"responseCode,omitempty"` + "`" + `
}`
	actual = string(bytes.ReplaceAll([]byte(actual), []byte("\t"), []byte("  ")))
	expected = string(bytes.ReplaceAll([]byte(expected), []byte("\t"), []byte("  ")))
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
		g, err := NewGoWSDL(file, "myservice", false, true)
		if err != nil {
			t.Error(err)
		}

		resp, err := g.Start()
		if err != nil {
			continue
			//t.Error(err)
		}

		data := new(bytes.Buffer)
		data.Write(resp["header"])
		data.Write(resp["types"])
		data.Write(resp["operations"])
		data.Write(resp["soap"])

		_, err = format.Source(data.Bytes())
		if err != nil {
			fmt.Println(string(data.Bytes()))
			t.Error(err)
		}
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

		re := regexp.MustCompile(varName + " " + typeName + " = \"([^\"]+)\"")
		matches := re.FindStringSubmatch(string(resp["types"]))

		if len(matches) != 2 {
			t.Errorf("No match or too many matches found for %s", varName)
		} else if matches[1] != enumString {
			t.Errorf("%s got '%s' but expected '%s'", varName, matches[1], enumString)
		}
	}
	enumStringTest(t, "chromedata.wsdl", "DriveTrainFrontWheelDrive", "DriveTrain", "Front Wheel Drive")
	enumStringTest(t, "vboxweb.wsdl", "SettingsVersionV1_14", "SettingsVersion", "v1_14")

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
	expectedBytes, err := ioutil.ReadFile("./fixtures/epcis/epcisquery.src")
	if err != nil {
		t.Fatal(err)
	}

	actual := string(source)
	expected := string(expectedBytes)
	if actual != expected {
		_ = ioutil.WriteFile("./fixtures/epcis/epcisquery_gen.src", source, 0664)
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
