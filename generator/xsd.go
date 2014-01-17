package generator

import (
	"encoding/xml"
)

type XsdSchema struct {
	XMLName            xml.Name          `xml:"http://www.w3.org/2001/XMLSchema schema"`
	Tns                string            `xml:"xmlns tns",attr`
	Xs                 string            `xml:"xmlns xs,attr"`
	Version            string            `xml:"version,attr"`
	TargetNamespace    string            `xml:"targetNamespace,attr"`
	ElementFormDefault string            `xml:"elementFormDefault,attr"`
	Includes           []*XsdInclude     `xml:"http://www.w3.org/2001/XMLSchema include"`
	Imports            []*XsdImport      `xml:"http://www.w3.org/2001/XMLSchema import"`
	Elements           []*XsdElement     `xml:"http://www.w3.org/2001/XMLSchema element"`
	ComplexTypes       []*XsdComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"` //global
	SimpleType         []*XsdSimpleType  `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
}

type XsdInclude struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type XsdImport struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2001/XMLSchema import"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Namespace      string   `xml:"namespace,attr"`
}

type XsdElement struct {
	XMLName     xml.Name        `xml:"http://www.w3.org/2001/XMLSchema element"`
	Name        string          `xml:"name,attr"`
	Nillable    bool            `xml:"nillable,attr"`
	Type        string          `xml:"type,attr"`
	Ref         string          `xml:"ref,attr"`
	MinOccurs   string          `xml:"minOccurs,attr"`
	MaxOccurs   string          `xml:"maxOccurs,attr"`
	ComplexType *XsdComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"` //local
	SimpleType  *XsdSimpleType  `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Groups      []*XsdGroup     `xml:"http://www.w3.org/2001/XMLSchema group"`
}

type XsdComplexType struct {
	XMLName        xml.Name          `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	Abstract       bool              `xml:"abstract,attr"`
	Name           string            `xml:"name,attr"`
	Mixed          bool              `xml:"mixed,attr"`
	Sequence       XsdSequence       `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	Choice         []XsdElement      `xml:"choice>element"`
	All            []XsdElement      `xml:"all>element"`
	ComplexContent XsdComplexContent `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	SimpleContent  XsdSimpleContent  `xml:"http://www.w3.org/2001/XMLSchema simpleContent"`
	Attributes     []*XsdAttribute   `xml:"http://www.w3.org/2001/XMLSchema attribute"`
}

type XsdGroup struct {
	Name     string       `xml:"name,attr"`
	Ref      string       `xml:"ref,attr"`
	Sequence XsdSequence  `xml:"http://www.w3.org/2001/XMLSchema sequence"`
	Choice   []XsdElement `xml:"http://www.w3.org/2001/XMLSchema choice"`
	All      []XsdElement `xml:"http://www.w3.org/2001/XMLSchema all"`
}

type XsdComplexContent struct {
	XMLName   xml.Name     `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension XsdExtension `xml:"http://www.w3.org/2001/XMLSchema extension"`
}

type XsdSimpleContent struct {
	XMLName   xml.Name     `xml:"http://www.w3.org/2001/XMLSchema simpleContent"`
	Extension XsdExtension `xml:"http://www.w3.org/2001/XMLSchema extension"`
}

type XsdExtension struct {
	XMLName    xml.Name        `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base       string          `xml:"base,attr"`
	Attributes []*XsdAttribute `xml:"http://www.w3.org/2001/XMLSchema attribute"`
	Sequence   XsdSequence     `xml:"http://www.w3.org/2001/XMLSchema sequence"`
}

type XsdAttribute struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type XsdSimpleType struct {
	Name        string         `xml:"name,attr"`
	Restriction XsdRestriction `xml:"http://www.w3.org/2001/XMLSchema restriction"`
}

type XsdSequence struct {
	Elements []XsdElement `xml:"http://www.w3.org/2001/XMLSchema element"`
}

type XsdRestriction struct {
	Base         string                `xml:"base,attr"`
	Enumeration  []XsdRestrictionValue `xml:"http://www.w3.org/2001/XMLSchema enumeration"`
	Pattern      XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema pattern"`
	MinInclusive XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema minInclusive"`
	MaxInclusive XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema maxInclusive"`
	WhiteSpace   XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema whitespace"`
	Length       XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema length"`
	MinLength    XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema minLength"`
	MaxLength    XsdRestrictionValue   `xml:"http://www.w3.org/2001/XMLSchema maxLength"`
}

type XsdRestrictionValue struct {
	Value string `xml:"value,attr"`
}
