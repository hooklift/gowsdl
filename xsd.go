// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"encoding/xml"
)

// XSDSchema represents an entire Schema structure.
type XSDSchema struct {
	XMLName            xml.Name          `xml:"schema"`
	Tns                string            `xml:"xmlns tns,attr"`
	Xs                 string            `xml:"xmlns xs,attr"`
	Version            string            `xml:"version,attr"`
	TargetNamespace    string            `xml:"targetNamespace,attr"`
	ElementFormDefault string            `xml:"elementFormDefault,attr"`
	Includes           []*XSDInclude     `xml:"include"`
	Imports            []*XSDImport      `xml:"import"`
	Elements           []*XSDElement     `xml:"element"`
	ComplexTypes       []*XSDComplexType `xml:"complexType"` //global
	SimpleType         []*XSDSimpleType  `xml:"simpleType"`
}

// XSDInclude represents schema includes.
type XSDInclude struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
}

// XSDImport represents XSD imports within the main schema.
type XSDImport struct {
	XMLName        xml.Name `xml:"import"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Namespace      string   `xml:"namespace,attr"`
}

// XSDElement represents a Schema element.
type XSDElement struct {
	XMLName     xml.Name        `xml:"element"`
	Name        string          `xml:"name,attr"`
	Doc         string          `xml:"annotation>documentation"`
	Nillable    bool            `xml:"nillable,attr"`
	Type        string          `xml:"type,attr"`
	Ref         string          `xml:"ref,attr"`
	MinOccurs   string          `xml:"minOccurs,attr"`
	MaxOccurs   string          `xml:"maxOccurs,attr"`
	ComplexType *XSDComplexType `xml:"complexType"` //local
	SimpleType  *XSDSimpleType  `xml:"simpleType"`
	Groups      []*XSDGroup     `xml:"group"`
}

// XSDComplexType represents a Schema complex type.
type XSDComplexType struct {
	XMLName        xml.Name          `xml:"complexType"`
	Abstract       bool              `xml:"abstract,attr"`
	Name           string            `xml:"name,attr"`
	Mixed          bool              `xml:"mixed,attr"`
	Sequence       []XSDElement      `xml:"sequence>element"`
	Choice         []XSDElement      `xml:"choice>element"`
	SequenceChoice []XSDElement      `xml:"sequence>choice>element"`
	All            []XSDElement      `xml:"all>element"`
	ComplexContent XSDComplexContent `xml:"complexContent"`
	SimpleContent  XSDSimpleContent  `xml:"simpleContent"`
	Attributes     []*XSDAttribute   `xml:"attribute"`
}

// XSDGroup element is used to define a group of elements to be used in complex type definitions.
type XSDGroup struct {
	Name     string       `xml:"name,attr"`
	Ref      string       `xml:"ref,attr"`
	Sequence []XSDElement `xml:"sequence>element"`
	Choice   []XSDElement `xml:"choice>element"`
	All      []XSDElement `xml:"all>element"`
}

// XSDComplexContent element defines extensions or restrictions on a complex
// type that contains mixed content or elements only.
type XSDComplexContent struct {
	XMLName   xml.Name     `xml:"complexContent"`
	Extension XSDExtension `xml:"extension"`
}

// XSDSimpleContent element contains extensions or restrictions on a text-only
// complex type or on a simple type as content and contains no elements.
type XSDSimpleContent struct {
	XMLName   xml.Name     `xml:"simpleContent"`
	Extension XSDExtension `xml:"extension"`
}

// XSDExtension element extends an existing simpleType or complexType element.
type XSDExtension struct {
	XMLName    xml.Name        `xml:"extension"`
	Base       string          `xml:"base,attr"`
	Attributes []*XSDAttribute `xml:"attribute"`
	Sequence   []XSDElement    `xml:"sequence>element"`
}

// XSDAttribute represent an element attribute. Simple elements cannot have
// attributes. If an element has attributes, it is considered to be of a
// complex type. But the attribute itself is always declared as a simple type.
type XSDAttribute struct {
	Name       string         `xml:"name,attr"`
	Doc        string         `xml:"annotation>documentation"`
	Type       string         `xml:"type,attr"`
	SimpleType *XSDSimpleType `xml:"simpleType"`
}

// XSDSimpleType element defines a simple type and specifies the constraints
// and information about the values of attributes or text-only elements.
type XSDSimpleType struct {
	Name        string         `xml:"name,attr"`
	Doc         string         `xml:"annotation>documentation"`
	Restriction XSDRestriction `xml:"restriction"`
}

// XSDRestriction defines restrictions on a simpleType, simpleContent, or complexContent definition.
type XSDRestriction struct {
	Base         string                `xml:"base,attr"`
	Enumeration  []XSDRestrictionValue `xml:"enumeration"`
	Pattern      XSDRestrictionValue   `xml:"pattern"`
	MinInclusive XSDRestrictionValue   `xml:"minInclusive"`
	MaxInclusive XSDRestrictionValue   `xml:"maxInclusive"`
	WhiteSpace   XSDRestrictionValue   `xml:"whitespace"`
	Length       XSDRestrictionValue   `xml:"length"`
	MinLength    XSDRestrictionValue   `xml:"minLength"`
	MaxLength    XSDRestrictionValue   `xml:"maxLength"`
}

// XSDRestrictionValue represents a restriction value.
type XSDRestrictionValue struct {
	Doc   string `xml:"annotation>documentation"`
	Value string `xml:"value,attr"`
}
