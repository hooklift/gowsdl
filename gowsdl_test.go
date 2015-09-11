package gowsdl

import (
    "testing"
    "strings"
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
