package generator

import (
	"encoding/xml"
)

type XsdSchema struct {
	XMLName            xml.Name          `xml:"schema"`
	Tns                string            `xml:"xmlns tns",attr`
	Xs                 string            `xml:"xmlns xs,attr"`
	Version            string            `xml:"version,attr"`
	TargetNamespace    string            `xml:"targetNamespace,attr"`
	ElementFormDefault string            `xml:"elementFormDefault,attr"`
	Includes           []*XsdInclude     `xml:"include"`
	Imports            []*XsdImport      `xml:"import"`
	Elements           []*XsdElement     `xml:"element"`
	ComplexTypes       []*XsdComplexType `xml:"complexType"` //global
	SimpleType         []*XsdSimpleType  `xml:"simpleType"`
}

type XsdInclude struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type XsdImport struct {
	XMLName        xml.Name `xml:"import"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Namespace      string   `xml:"namespace,attr"`
}

type XsdElement struct {
	XMLName     xml.Name        `xml:"element"`
	Name        string          `xml:"name,attr"`
	Nillable    bool            `xml:"nillable,attr"`
	Type        string          `xml:"type,attr"`
	Ref         string          `xml:"ref,attr"`
	MinOccurs   string          `xml:"minOccurs,attr"`
	MaxOccurs   string          `xml:"maxOccurs,attr"`
	ComplexType *XsdComplexType `xml:"complexType"` //local
	SimpleType  *XsdSimpleType  `xml:"simpleType"`
	Groups      []*XsdGroup     `xml:"group"`
}

type XsdComplexType struct {
	XMLName        xml.Name          `xml:"complexType"`
	Abstract       bool              `xml:"abstract,attr"`
	Name           string            `xml:"name,attr"`
	Mixed          bool              `xml:"mixed,attr"`
	Sequence       XsdSequence       `xml:"sequence"`
	Choice         []XsdElement      `xml:"choice>element"`
	All            []XsdElement      `xml:"all>element"`
	ComplexContent XsdComplexContent `xml:"complexContent"`
	SimpleContent  XsdSimpleContent  `xml:"simpleContent"`
	Attributes     []*XsdAttribute   `xml:"attribute"`
}

type XsdGroup struct {
	Name     string       `xml:"name,attr"`
	Ref      string       `xml:"ref,attr"`
	Sequence XsdSequence  `xml:"sequence"`
	Choice   []XsdElement `xml:"choice"`
	All      []XsdElement `xml:"all"`
}

type XsdComplexContent struct {
	XMLName   xml.Name     `xml:"complexContent"`
	Extension XsdExtension `xml:"extension"`
}

type XsdSimpleContent struct {
	XMLName   xml.Name     `xml:"simpleContent"`
	Extension XsdExtension `xml:"extension"`
}

type XsdExtension struct {
	XMLName    xml.Name        `xml:"extension"`
	Base       string          `xml:"base,attr"`
	Attributes []*XsdAttribute `xml:"attribute"`
	Sequence   XsdSequence     `xml:"sequence"`
}

type XsdAttribute struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type XsdSimpleType struct {
	Name        string         `xml:"name,attr"`
	Restriction XsdRestriction `xml:"restriction"`
}

type XsdSequence struct {
	Elements []XsdElement `xml:"element"`
}

type XsdRestriction struct {
	Base         string                `xml:"base,attr"`
	Enumeration  []XsdRestrictionValue `xml:"enumeration"`
	Pattern      XsdRestrictionValue   `xml:"pattern"`
	MinInclusive XsdRestrictionValue   `xml:"minInclusive"`
	MaxInclusive XsdRestrictionValue   `xml:"maxInclusive"`
	WhiteSpace   XsdRestrictionValue   `xml:"whitespace"`
	Length       XsdRestrictionValue   `xml:"length"`
	MinLength    XsdRestrictionValue   `xml:"minLength"`
	MaxLength    XsdRestrictionValue   `xml:"maxLength"`
}

type XsdRestrictionValue struct {
	Value string `xml:"value,attr"`
}
