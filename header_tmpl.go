// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var headerTmpl = `
// Code generated by gowsdl DO NOT EDIT.

package {{.}}

import (
	"encoding/xml"
	"fmt"
	"time"
	"github.com/factrylabs/gowsdl/soap"

	{{/*range .Imports*/}}
		{{/*.*/}}
	{{/*end*/}}
)

// against "unused imports"
var _ time.Time
var _ xml.Name


// YB : Added custom timestamp to parse all kinds of different datetime formats
type CustomTimestamp struct {
	time.Time
}

// UnmarshalXML will parse the time.Time in different ways - as Navision uses them
func (c *CustomTimestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)

	parse, err := time.Parse("2006-01-02T15:04:05Z07:00", v)
	if err == nil {
		*c = CustomTimestamp{parse}
		return nil
	}

	parse, err = time.Parse("2006-01-02T15:04:05", v)
	if err == nil {
		*c = CustomTimestamp{parse}
		return nil
	}

	parse, err = time.Parse("2006-01-02", v)
	if err == nil {
		*c = CustomTimestamp{parse}
		return nil
	}

	parse, err = time.Parse("15:04:05", v)
	if err == nil {
		*c = CustomTimestamp{parse}
		return nil
	}

	// if we reach here: an error occurred..
	return fmt.Errorf("Could not parse datetime for %v", v)

}
`
