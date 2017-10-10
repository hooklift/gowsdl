// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

import "encoding/xml"

const wsdlNamespace = "http://schemas.xmlsoap.org/wsdl/"

// WSDL represents the global structure of a WSDL file.
type WSDL struct {
	Xmlns           map[string]string `xml:"-"`
	Name            string          `xml:"name,attr"`
	TargetNamespace string          `xml:"targetNamespace,attr"`
	Imports         []*WSDLImport   `xml:"import"`
	Doc             string          `xml:"documentation"`
	Types           WSDLType        `xml:"http://schemas.xmlsoap.org/wsdl/ types"`
	Messages        []*WSDLMessage  `xml:"http://schemas.xmlsoap.org/wsdl/ message"`
	PortTypes       []*WSDLPortType `xml:"http://schemas.xmlsoap.org/wsdl/ portType"`
	Binding         []*WSDLBinding  `xml:"http://schemas.xmlsoap.org/wsdl/ binding"`
	Service         []*WSDLService  `xml:"http://schemas.xmlsoap.org/wsdl/ service"`
}

// UnmarshalXML implements interface xml.Unmarshaler for XSDSchema.
func (w *WSDL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	w.Xmlns = make(map[string]string)
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			w.Xmlns[attr.Name.Local] = attr.Value
			continue
		}

		switch attr.Name.Local {
		case "name":
			w.Name = attr.Value
		case "targetNamespace":
			w.TargetNamespace = attr.Value
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
			switch {
			case t.Name.Local == "import":
				x := new(WSDLImport)
				if err := d.DecodeElement(x, &t); err != nil {
					return err
				}
				w.Imports = append(w.Imports, x)
			case t.Name.Local == "documentation":
				if err := d.DecodeElement(&w.Doc, &t); err != nil {
					return err
				}
			case t.Name.Space == wsdlNamespace:
				switch t.Name.Local {
				case "types":
					if err := d.DecodeElement(&w.Types, &t); err != nil {
						return err
					}
					for prefix, namespace := range w.Xmlns {
						for _, s := range w.Types.Schemas {
							if _, ok := s.Xmlns[prefix]; !ok {
								s.Xmlns[prefix] = namespace
							}
						}
					}
				case "message":
					x := new(WSDLMessage)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}
					w.Messages = append(w.Messages, x)
				case "portType":
					x := new(WSDLPortType)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}
					w.PortTypes = append(w.PortTypes, x)
				case "binding":
					x := new(WSDLBinding)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}
					w.Binding = append(w.Binding, x)
				case "service":
					x := new(WSDLService)
					if err := d.DecodeElement(x, &t); err != nil {
						return err
					}
					w.Service = append(w.Service, x)
				default:
					d.Skip()
					continue Loop
				}
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

// WSDLImport is the struct used for deserializing WSDL imports.
type WSDLImport struct {
	Namespace string `xml:"namespace,attr"`
	Location  string `xml:"location,attr"`
}

// WSDLType represents the entry point for deserializing XSD schemas used by the WSDL file.
type WSDLType struct {
	Doc     string       `xml:"documentation"`
	Schemas []*XSDSchema `xml:"schema"`
}

// WSDLPart defines the struct for a function parameter within a WSDL.
type WSDLPart struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}

// WSDLMessage represents a function, which in turn has one or more parameters.
type WSDLMessage struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"documentation"`
	Parts []*WSDLPart `xml:"http://schemas.xmlsoap.org/wsdl/ part"`
}

// WSDLFault represents a WSDL fault message.
type WSDLFault struct {
	Name      string        `xml:"name,attr"`
	Message   string        `xml:"message,attr"`
	Doc       string        `xml:"documentation"`
	SOAPFault WSDLSOAPFault `xml:"http://schemas.xmlsoap.org/wsdl/soap/ fault"`
}

// WSDLInput represents a WSDL input message.
type WSDLInput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"documentation"`
	SOAPBody   WSDLSOAPBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SOAPHeader []*WSDLSOAPHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
}

// WSDLOutput represents a WSDL output message.
type WSDLOutput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"documentation"`
	SOAPBody   WSDLSOAPBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SOAPHeader []*WSDLSOAPHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
}

// WSDLOperation represents the contract of an entire operation or function.
type WSDLOperation struct {
	Name          string            `xml:"name,attr"`
	Doc           string            `xml:"documentation"`
	Input         WSDLInput         `xml:"input"`
	Output        WSDLOutput        `xml:"output"`
	Faults        []*WSDLFault      `xml:"fault"`
	SOAPOperation WSDLSOAPOperation `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
}

// WSDLPortType defines the service, operations that can be performed and the messages involved.
// A port type can be compared to a function library, module or class.
type WSDLPortType struct {
	Name       string           `xml:"name,attr"`
	Doc        string           `xml:"documentation"`
	Operations []*WSDLOperation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}

// WSDLSOAPBinding represents a SOAP binding to the web service.
type WSDLSOAPBinding struct {
	Style     string `xml:"style,attr"`
	Transport string `xml:"transport,attr"`
}

// WSDLSOAPOperation represents a service operation in SOAP terms.
type WSDLSOAPOperation struct {
	SOAPAction string `xml:"soapAction,attr"`
	Style      string `xml:"style,attr"`
}

// WSDLSOAPHeader defines the header for a SOAP service.
type WSDLSOAPHeader struct {
	Message       string                 `xml:"message,attr"`
	Part          string                 `xml:"part,attr"`
	Use           string                 `xml:"use,attr"`
	EncodingStyle string                 `xml:"encodingStyle,attr"`
	Namespace     string                 `xml:"namespace,attr"`
	HeadersFault  []*WSDLSOAPHeaderFault `xml:"headerfault"`
}

// WSDLSOAPHeaderFault defines a SOAP fault header.
type WSDLSOAPHeaderFault struct {
	Message       string `xml:"message,attr"`
	Part          string `xml:"part,attr"`
	Use           string `xml:"use,attr"`
	EncodingStyle string `xml:"encodingStyle,attr"`
	Namespace     string `xml:"namespace,attr"`
}

// WSDLSOAPBody defines SOAP body characteristics.
type WSDLSOAPBody struct {
	Parts         string `xml:"parts,attr"`
	Use           string `xml:"use,attr"`
	EncodingStyle string `xml:"encodingStyle,attr"`
	Namespace     string `xml:"namespace,attr"`
}

// WSDLSOAPFault defines a SOAP fault message characteristics.
type WSDLSOAPFault struct {
	Parts         string `xml:"parts,attr"`
	Use           string `xml:"use,attr"`
	EncodingStyle string `xml:"encodingStyle,attr"`
	Namespace     string `xml:"namespace,attr"`
}

// WSDLSOAPAddress defines the location for the SOAP service.
type WSDLSOAPAddress struct {
	Location string `xml:"location,attr"`
}

// WSDLBinding defines only a SOAP binding and its operations
type WSDLBinding struct {
	Name        string           `xml:"name,attr"`
	Type        string           `xml:"type,attr"`
	Doc         string           `xml:"documentation"`
	SOAPBinding WSDLSOAPBinding  `xml:"http://schemas.xmlsoap.org/wsdl/soap/ binding"`
	Operations  []*WSDLOperation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}

// WSDLPort defines the properties for a SOAP port only.
type WSDLPort struct {
	Name        string          `xml:"name,attr"`
	Binding     string          `xml:"binding,attr"`
	Doc         string          `xml:"documentation"`
	SOAPAddress WSDLSOAPAddress `xml:"http://schemas.xmlsoap.org/wsdl/soap/ address"`
}

// WSDLService defines the list of SOAP services associated with the WSDL.
type WSDLService struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"documentation"`
	Ports []*WSDLPort `xml:"http://schemas.xmlsoap.org/wsdl/ port"`
}
