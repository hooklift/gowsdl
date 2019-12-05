// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"strings"
)

const xmlschema11 = "http://www.w3.org/2001/XMLSchema"

type NameSpaceCache struct {
	latest string
	cache  map[string]string
}

var nameSpaceCache NameSpaceCache = NameSpaceCache{
	latest: "ns",
	cache:  map[string]string{},
}

// XSDSchema represents an entire Schema structure.
type XSDSchema struct {
	XMLName            xml.Name             `xml:"schema"`
	Xmlns              map[string]string    `xml:"-"`
	Tns                string               `xml:"xmlns tns,attr"`
	Xs                 string               `xml:"xmlns xs,attr"`
	Version            string               `xml:"version,attr"`
	TargetNamespace    string               `xml:"targetNamespace,attr"`
	ElementFormDefault string               `xml:"elementFormDefault,attr"`
	Includes           []*XSDInclude        `xml:"include"`
	Imports            []*XSDImport         `xml:"import"`
	Elements           []*XSDElement        `xml:"element"`
	Attributes         []*XSDAttribute      `xml:"attribute"`
	ComplexTypes       []*XSDComplexType    `xml:"complexType"` //global
	SimpleType         []*XSDSimpleType     `xml:"simpleType"`
	AttributeGroups    []*XSDAttributeGroup `xml:"attributeGroup"`
}

// UnmarshalXML implements interface xml.Unmarshaler for XSDSchema.
func (s *XSDSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	nameSpaceCache.cache["http://www.iata.org/IATA/EDIST/2017.2"] = "ns"

	s.Xmlns = make(map[string]string)
	s.XMLName = start.Name
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			s.Xmlns[attr.Name.Local] = attr.Value
			continue
		}

		switch attr.Name.Local {
		case "version":
			s.Version = attr.Value
		case "targetNamespace":
			s.TargetNamespace = attr.Value
		case "elementFormDefault":
			s.ElementFormDefault = attr.Value
		}
	}

Loop:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Space != xmlschema11 {
				d.Skip()
				continue Loop
			}

			switch t.Name.Local {
			case "include":
				x := new(XSDInclude)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				s.Includes = append(s.Includes, x)
			case "import":
				x := new(XSDImport)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				s.Imports = append(s.Imports, x)
			case "element":
				x := new(XSDElement)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				//checkNameSpace(s.TargetNamespace, x)
				s.Elements = append(s.Elements, x)
			case "attribute":
				x := new(XSDAttribute)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				s.Attributes = append(s.Attributes, x)
			case "complexType":
				x := new(XSDComplexType)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				//addDisplayNameForComplexTypes(s.TargetNamespace, x)
				s.ComplexTypes = append(s.ComplexTypes, x)
			case "simpleType":
				x := new(XSDSimpleType)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				//processXSDSimpleType(s.TargetNamespace, x)
				s.SimpleType = append(s.SimpleType, x)
			case "attributeGroup":
				x := new(XSDAttributeGroup)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				s.AttributeGroups = append(s.AttributeGroups, x)
			default:
				d.Skip()
				continue Loop
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
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
	XMLName        xml.Name             `xml:"element"`
	Name           string               `xml:"name,attr"`
	Doc            string               `xml:"annotation>documentation"`
	Nillable       bool                 `xml:"nillable,attr"`
	Type           string               `xml:"type,attr"`
	Ref            string               `xml:"ref,attr"`
	MinOccurs      string               `xml:"minOccurs,attr"`
	MaxOccurs      string               `xml:"maxOccurs,attr"`
	ComplexType    *XSDComplexType      `xml:"complexType"` //local
	SimpleType     *XSDSimpleType       `xml:"simpleType"`
	Groups         []*XSDGroup          `xml:"group"`
	AttributeGroup []*XSDAttributeGroup `xml:"attributeGroup"`
	DisplayName    string
	RefDisplayName string
}

// XSDComplexType represents a Schema complex type.
type XSDComplexType struct {
	XMLName                xml.Name             `xml:"complexType"`
	Abstract               bool                 `xml:"abstract,attr"`
	Name                   string               `xml:"name,attr"`
	Mixed                  bool                 `xml:"mixed,attr"`
	Sequence               []*XSDElement        `xml:"sequence>element"`
	Choice                 []*XSDElement        `xml:"choice>element"`
	SequenceChoice         []*XSDElement        `xml:"sequence>choice>element"`
	SequenceChoiceSequence []*XSDElement        `xml:"sequence>choice>sequence>element"`
	All                    []*XSDElement        `xml:"all>element"`
	ComplexContent         XSDComplexContent    `xml:"complexContent"`
	SimpleContent          XSDSimpleContent     `xml:"simpleContent"`
	Attributes             []*XSDAttribute      `xml:"attribute"`
	ChoiceSequence         []*XSDElement        `xml:"choice>sequence>element"`
	AttributeGroup         []*XSDAttributeGroup `xml:"attributeGroup"`
	SequenceSequence       []*XSDElement        `xml:"sequence>sequence>element"`
}

// XSDGroup element is used to define a group of elements to be used in complex type definitions.
type XSDGroup struct {
	Name     string        `xml:"name,attr"`
	Ref      string        `xml:"ref,attr"`
	Sequence []*XSDElement `xml:"sequence>element"`
	Choice   []*XSDElement `xml:"choice>element"`
	All      []*XSDElement `xml:"all>element"`
}

// XSDComplexContent element defines extensions or restrictions on a complex
// type that contains mixed content or elements only.
type XSDComplexContent struct {
	XMLName     xml.Name       `xml:"complexContent"`
	Extension   XSDExtension   `xml:"extension"`
	Restriction XSDRestriction `xml:"restriction"`
}

// XSDSimpleContent element contains extensions or restrictions on a text-only
// complex type or on a simple type as content and contains no elements.
type XSDSimpleContent struct {
	XMLName   xml.Name     `xml:"simpleContent"`
	Extension XSDExtension `xml:"extension"`
}
type XSDAttributeGroup struct {
	XMLName        xml.Name             `xml:"attributeGroup"`
	Name           string               `xml:"name,attr"`
	Ref            string               `xml:"ref,attr"`
	Attributes     []*XSDAttribute      `xml:"attribute"`
	AttributeGroup []*XSDAttributeGroup `xml:"attributeGroup"`
}

// XSDExtension element extends an existing simpleType or complexType element.
type XSDExtension struct {
	XMLName                xml.Name             `xml:"extension"`
	Base                   string               `xml:"base,attr"`
	Attributes             []*XSDAttribute      `xml:"attribute"`
	Sequence               []*XSDElement        `xml:"sequence>element"`
	SequenceChoice         []*XSDElement        `xml:"sequence>choice>element"`
	SequenceChoiceSequence []*XSDElement        `xml:"sequence>choice>sequence>element"`
	Choice                 []*XSDElement        `xml:"choice>element"`
	AttributeGroup         []*XSDAttributeGroup `xml:"attributeGroup"`
}

// XSDAttribute represent an element attribute. Simple elements cannot have
// attributes. If an element has attributes, it is considered to be of a
// complex type. But the attribute itself is always declared as a simple type.
type XSDAttribute struct {
	Doc        string         `xml:"annotation>documentation"`
	Name       string         `xml:"name,attr"`
	Ref        string         `xml:"ref,attr"`
	Type       string         `xml:"type,attr"`
	Use        string         `xml:"use,attr"`
	Fixed      string         `xml:"fixed,attr"`
	SimpleType *XSDSimpleType `xml:"simpleType"`
}

// XSDSimpleType element defines a simple type and specifies the constraints
// and information about the values of attributes or text-only elements.
type XSDSimpleType struct {
	Name        string         `xml:"name,attr"`
	Doc         string         `xml:"annotation>documentation"`
	Restriction XSDRestriction `xml:"restriction"`
	List        XSDList        `xml:"list"`
	Union       XSDUnion       `xml:"union"`
	Final       string         `xml:"final"`
}

// XSDList represents a element list
type XSDList struct {
	Doc        string         `xml:"annotation>documentation"`
	ItemType   string         `xml:"itemType,attr"`
	SimpleType *XSDSimpleType `xml:"simpleType"`
}

// XSDUnion represents a union element
type XSDUnion struct {
	SimpleType  []*XSDSimpleType `xml:"simpleType,omitempty"`
	MemberTypes string           `xml:"memberTypes,attr"`
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
	Sequence     []*XSDElement         `xml:"sequence>element"`
	Attributes   []*XSDAttribute       `xml:"attribute"`
	SimpleType   *XSDSimpleType        `xml:"simpleType"`
}

// XSDRestrictionValue represents a restriction value.
type XSDRestrictionValue struct {
	Doc   string `xml:"annotation>documentation"`
	Value string `xml:"value,attr"`
}

func checkNameSpace(targetNameSpaceName string, element *XSDElement) {
	var b bytes.Buffer
	latest := nameSpaceCache.latest
	if len(nameSpaceCache.cache) > 0 && nameSpaceCache.cache[targetNameSpaceName] != "" {
		latest = nameSpaceCache.cache[targetNameSpaceName]
	} else {
		var tokens []string = strings.Split(latest, "s")
		var num int
		if len(tokens) < 2 {
			num = 1
		} else {
			num, _ = strconv.Atoi(tokens[1])
			num = num + 1
		}
		b.WriteString("ns")
		b.WriteString(strconv.Itoa(num))
		latest = b.String()
		nameSpaceCache.latest = latest
		nameSpaceCache.cache[targetNameSpaceName] = latest
	}
	element.DisplayName = latest + ":" + element.Name
	if element.Ref != "" {
		element.RefDisplayName = removeNameSpace(latest, element.Ref)
	}

	if element.ComplexType != nil {
		addDisplayNameForComplexTypes(targetNameSpaceName, element.ComplexType)
	}
	if element.Groups != nil {
		for _, element := range element.Groups {
			processXSDGroup(targetNameSpaceName, element)
		}
	}

}

func processXSDGroup(targetNameSpaceName string, group *XSDGroup) {
	if group.Sequence != nil {
		for _, element := range group.Sequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if group.Choice != nil {
		for _, element := range group.Choice {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if group.All != nil {
		for _, element := range group.All {
			checkNameSpace(targetNameSpaceName, element)
		}
	}

}

func addDisplayNameForComplexTypes(targetNameSpaceName string, complexType *XSDComplexType) {
	if complexType.ChoiceSequence != nil {
		for _, element := range complexType.ChoiceSequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if complexType.Choice != nil {
		for _, element := range complexType.Choice {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if complexType.All != nil {
		for _, element := range complexType.All {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if complexType.Sequence != nil {
		for _, element := range complexType.Sequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if complexType.SequenceChoice != nil {
		for _, element := range complexType.SequenceChoice {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if complexType.SequenceChoiceSequence != nil {
		for _, element := range complexType.SequenceChoiceSequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if complexType.SequenceSequence != nil {
		for _, element := range complexType.SequenceSequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	complexType.SimpleContent = processXSDSimpleContent(targetNameSpaceName, complexType.SimpleContent)
	complexType.ComplexContent = processXSDComplexContent(targetNameSpaceName, complexType.ComplexContent)

}

func processXSDSimpleContent(targetNameSpaceName string, simpleContent XSDSimpleContent) XSDSimpleContent {
	simpleContent.Extension = processExtension(targetNameSpaceName, simpleContent.Extension)
	return simpleContent
}
func processXSDComplexContent(targetNameSpaceName string, complexContent XSDComplexContent) XSDComplexContent {
	complexContent.Extension = processExtension(targetNameSpaceName, complexContent.Extension)
	return complexContent
}

func processExtension(targetNameSpaceName string, extension XSDExtension) XSDExtension {
	if extension.SequenceChoiceSequence != nil {
		for _, element := range extension.SequenceChoiceSequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if extension.SequenceChoice != nil {
		for _, element := range extension.SequenceChoice {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if extension.Sequence != nil {
		for _, element := range extension.Sequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}

	if extension.Choice != nil {
		for _, element := range extension.Choice {
			checkNameSpace(targetNameSpaceName, element)
		}
	}

	return extension
}

func processRestriction(targetNameSpaceName string, restriction XSDRestriction) XSDRestriction {
	if restriction.Sequence != nil {
		for _, element := range restriction.Sequence {
			checkNameSpace(targetNameSpaceName, element)
		}
	}
	if restriction.SimpleType != nil {
		processXSDSimpleType(targetNameSpaceName, restriction.SimpleType)
	}

	return restriction
}

func processXSDSimpleType(targetNameSpaceName string, simpleType *XSDSimpleType) {
	simpleType.Restriction = processRestriction(targetNameSpaceName, simpleType.Restriction)

	if simpleType.List.SimpleType != nil {
		processXSDSimpleType(targetNameSpaceName, simpleType.List.SimpleType)
	}
	if simpleType.Union.SimpleType != nil {
		for _, element := range simpleType.Union.SimpleType {
			processXSDSimpleType(targetNameSpaceName, element)
		}
	}
}

func removeNameSpace(latestSpace string, xsdType string) string {
	// Handles name space, ie. xsd:string, xs:string
	r := strings.Split(xsdType, ":")

	if len(r) == 2 {
		return r[1]
	}

	return latestSpace + ":" + r[0]
}
