// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"encoding/xml"
	"os"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data, err := os.ReadFile("fixtures/vim.wsdl")
	if err != nil {
		t.Errorf("incorrect result\ngot:  %#v\nwant: %#v", err, nil)
	}

	v := WSDL{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		t.Errorf("incorrect result\ngot:  %#v\nwant: %#v", err, nil)
	}
}
