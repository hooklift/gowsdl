// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"bytes"
	"go/format"
	"path/filepath"
	"strings"
	"testing"
)

func TestElementGenerationDoesntCommentOutStructProperty(t *testing.T) {
	g := GoWSDL{
		file:         "fixtures/test.wsdl",
		pkg:          "myservice",
		makePublicFn: makePublic,
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

func TestVboxGeneratesWithoutSyntaxErrors(t *testing.T) {
	files, err := filepath.Glob("fixtures/*.wsdl")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		g := GoWSDL{
			file:         file,
			pkg:          "myservice",
			makePublicFn: makePublic,
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
			t.Error(err)
		}
	}
}
