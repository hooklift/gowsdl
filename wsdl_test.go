package main

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("vim.wsdl")
	if err != nil {
		t.Errorf("incorrect result\ngot:  %#v\nwant: %#v", err, nil)
	}

	v := Wsdl{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		t.Errorf("incorrect result\ngot:  %#v\nwant: %#v", err, nil)
	}

	// for _, pt := range v.PortTypes {
	// 	t.Logf("PortType name: %s\n", pt.Name)
	// 	for _, o := range pt.Operations {
	// 		t.Logf("Operation: %s", o.Name)
	// 	}
	// 	t.Logf("Total ops: %d\n", len(pt.Operations))
	// }

	// t.Logf("%#v\n", v.Types.Schema[0].Includes)
}
