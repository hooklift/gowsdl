// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocation_ParseLocation_URL(t *testing.T) {
	r, err := ParseLocation("http://example.org/my.wsdl")
	if err != nil {
		t.Fatal(err)
	}

	if !r.isURL() || r.isFile() {
		t.Error("Location should be a URL type")
	}
	if r.String() != "http://example.org/my.wsdl" {
		t.Error("got " + r.String() + " wanted " + "http://example.org/my.wsdl")
	}
}

func TestLocation_Parse_URL(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{"http://example.org/my.wsdl", "some.xsd", "http://example.org/some.xsd"},
		{"http://example.org/folder/my.wsdl", "some.xsd", "http://example.org/folder/some.xsd"},
		{"http://example.org/folder/my.wsdl", "../some.xsd", "http://example.org/some.xsd"},
	}
	for _, test := range tests {
		r, err := ParseLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}
		r, err = r.Parse(test.ref)
		if err != nil {
			t.Error(err)
			continue
		}

		if !r.isURL() || r.isFile() {
			t.Error("Location should be a URL type")
		}
		if r.String() != test.expected {
			t.Error("got " + r.String() + " wanted " + test.name)
		}
	}
}

func TestLocation_ParseLocation_File(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"fixtures/test.wsdl"},
		{"cmd/../fixtures/test.wsdl"},
	}
	for _, test := range tests {
		r, err := ParseLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}

		if r.isURL() || !r.isFile() {
			t.Error("Location should be a FILE type")
			continue
		}
		if !filepath.IsAbs(r.String()) {
			t.Error("Path should be absolute")
		}
		if _, err := os.Stat(r.String()); err != nil {
			t.Errorf("Location should point to existing loc: %s", err.Error())
		}
	}
}

func TestLocation_Parse_File(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{"fixtures/test.wsdl", "some.xsd", "fixtures/some.xsd"},
		{"fixtures/test.wsdl", "../xsd/some.xsd", "xsd/some.xsd"},
		{"fixtures/test.wsdl", "xsd/some.xsd", "fixtures/xsd/some.xsd"},
	}
	for _, test := range tests {
		r, err := ParseLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}
		r, err = r.Parse(test.ref)
		if err != nil {
			t.Error(err)
			continue
		}

		if r.isURL() || !r.isFile() {
			t.Error("Location should be a File type")
			continue
		}
		x, _ := filepath.Abs("")
		rel, _ := filepath.Rel(x, r.String())
		if rel != test.expected {
			t.Error("got " + rel + " wanted " + test.expected)
		}
	}
}

func TestLocation_Parse_FileToURL(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{"fixtures/test.wsdl", "http://example.org/some.xsd", "http://example.org/some.xsd"},
	}
	for _, test := range tests {
		r, err := ParseLocation(test.name)
		if err != nil {
			t.Error(err)
			continue
		}
		r, err = r.Parse(test.ref)
		if err != nil {
			t.Error(err)
			continue
		}

		if !r.isURL() || r.isFile() {
			t.Error("Location should be a URL type")
			continue
		}
		if r.String() != test.expected {
			t.Error("got " + r.String() + " wanted " + test.expected)
		}
	}
}
