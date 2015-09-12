package gowsdl

import (
    "testing"
    "strings"
    "go/format"
    "bytes"
    "path/filepath"
)

func TestElementGenerationDoesntCommentOutStructProperty(t *testing.T) {
    g := GoWsdl{
        file: "fixtures/test.wsdl",
        pkg: "myservice",
    }

    resp,err := g.Start()
    if err != nil {
        t.Error(err)
    }

    if (strings.Contains(string(resp["types"]), "// this is a comment  GetInfoResult string `xml:\"GetInfoResult,omitempty\"`")) {
        t.Error("Type comment should not comment out struct type property")
        t.Error(string(resp["types"]))
    }
}

func TestVboxGeneratesWithoutSyntaxErrors(t *testing.T) {
    files, err := filepath.Glob("fixtures/*.wsdl")
    if err != nil {
        t.Error(err)
    }

    for _,file := range files {
        g := GoWsdl{
            file: file,
            pkg: "myservice",
        }

        resp,err := g.Start()
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
    		t.Error(err)
    	}
    }
}
