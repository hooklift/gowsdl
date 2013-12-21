package main

// import (
// 	"encoding/xml"
// )

type XsdSchema struct {
	TargetNamespace    string            `xml:"targetNamespace,attr"`
	ElementFormDefault string            `xml:"elementFormDefault,attr"`
	Includes           []*XsdInclude     `xml:"http://www.w3.org/2001/XMLSchema include"`
	Imports            []*XsdImport      `xml:"http://www.w3.org/2001/XMLSchema import"`
	Elements           []*XsdElement     `xml:"http://www.w3.org/2001/XMLSchema element"`
	ComplexTypes       []*XsdComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"`
}

type XsdInclude struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type XsdImport struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
	Namespace      string `xml:"namespace,attr"`
}

type XsdElement struct {
	Name        string          `xml:"name,attr"`
	Nillable    bool            `xml:"nillable,attr"`
	Type        string          `xml:"type,attr"`
	MinOccurs   string          `xml:"minOccurs,attr"`
	MaxOccurs   string          `xml:"maxOccurs,attr"`
	ComplexType *XsdComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	SimpleType  *XsdSimpleType  `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
}

type XsdComplexType struct {
	Name     string       `xml:"name,attr"`
	Sequence *XsdSequence `xml:"http://www.w3.org/2001/XMLSchema sequence"`
}

type XsdSimpleType struct {
	Name     string          `xml:"name,attr"`
	Sequence *XsdRestriction `xml:"http://www.w3.org/2001/XMLSchema restriction"`
}

type XsdSequence struct {
	Elements []*XsdElement `xml:"http://www.w3.org/2001/XMLSchema element"`
}

type XsdRestriction struct {
	Base         string           `xml:"base,attr"`
	Pattern      *XsdPattern      `xml:"http://www.w3.org/2001/XMLSchema pattern"`
	MinInclusive *XsdMinInclusive `xml:"http://www.w3.org/2001/XMLSchema minInclusive"`
	MaxInclusive *XsdMaxInclusive `xml:"http://www.w3.org/2001/XMLSchema maxInclusive"`
}

type XsdPattern struct {
	Value string `xml:"value,attr"`
}

type XsdMinInclusive struct {
	Value string `xml:"value,attr"`
}

type XsdMaxInclusive struct {
	Value string `xml:"value,attr"`
}
