// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package gowsdl

type Wsdl struct {
	Name            string          `xml:"name,attr"`
	TargetNamespace string          `xml:"targetNamespace,attr"`
	Imports         []*WsdlImport   `xml:"import"`
	Doc             string          `xml:"documentation"`
	Types           WsdlType        `xml:"http://schemas.xmlsoap.org/wsdl/ types"`
	Messages        []*WsdlMessage  `xml:"http://schemas.xmlsoap.org/wsdl/ message"`
	PortTypes       []*WsdlPortType `xml:"http://schemas.xmlsoap.org/wsdl/ portType"`
	Binding         []*WsdlBinding  `xml:"http://schemas.xmlsoap.org/wsdl/ binding"`
	Service         []*WsdlService  `xml:"http://schemas.xmlsoap.org/wsdl/ service"`
}

type WsdlImport struct {
	Namespace string `xml:"namespace,attr"`
	Location  string `xml:"location,attr"`
}

type WsdlType struct {
	Doc     string       `xml:"documentation"`
	Schemas []*XsdSchema `xml:"schema"`
}

type WsdlPart struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}

type WsdlMessage struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"documentation"`
	Parts []*WsdlPart `xml:"http://schemas.xmlsoap.org/wsdl/ part"`
}

type WsdlFault struct {
	Name      string        `xml:"name,attr"`
	Message   string        `xml:"message,attr"`
	Doc       string        `xml:"documentation"`
	SoapFault WsdlSoapFault `xml:"http://schemas.xmlsoap.org/wsdl/soap/ fault"`
}

type WsdlInput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"documentation"`
	SoapBody   WsdlSoapBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SoapHeader []*WsdlSoapHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
}

type WsdlOutput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"documentation"`
	SoapBody   WsdlSoapBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SoapHeader []*WsdlSoapHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
}

type WsdlOperation struct {
	Name          string            `xml:"name,attr"`
	Doc           string            `xml:"documentation"`
	Input         WsdlInput         `xml:"input"`
	Output        WsdlOutput        `xml:"output"`
	Faults        []*WsdlFault      `xml:"fault"`
	SoapOperation WsdlSoapOperation `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
}

type WsdlPortType struct {
	Name       string           `xml:"name,attr"`
	Doc        string           `xml:"documentation"`
	Operations []*WsdlOperation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}

type WsdlSoapBinding struct {
	Style     string `xml:"style,attr"`
	Transport string `xml:"transport,attr"`
}

type WsdlSoapOperation struct {
	SoapAction string `xml:"soapAction,attr"`
	Style      string `xml:"style,attr"`
}

type WsdlSoapHeader struct {
	Message       string                 `xml:"message,attr"`
	Part          string                 `xml:"part,attr"`
	Use           string                 `xml:"use,attr"`
	EncodingStyle string                 `xml:"encodingStyle,attr"`
	Namespace     string                 `xml:"namespace,attr"`
	HeadersFault  []*WsdlSoapHeaderFault `xml:"headerfault"`
}

type WsdlSoapHeaderFault struct {
	Message       string `xml:"message,attr"`
	Part          string `xml:"part,attr"`
	Use           string `xml:"use,attr"`
	EncodingStyle string `xml:"encodingStyle,attr"`
	Namespace     string `xml:"namespace,attr"`
}

type WsdlSoapBody struct {
	Parts         string `xml:"parts,attr"`
	Use           string `xml:"use,attr"`
	EncodingStyle string `xml:"encodingStyle,attr"`
	Namespace     string `xml:"namespace,attr"`
}

type WsdlSoapFault struct {
	Parts         string `xml:"parts,attr"`
	Use           string `xml:"use,attr"`
	EncodingStyle string `xml:"encodingStyle,attr"`
	Namespace     string `xml:"namespace,attr"`
}

type WsdlSoapAddress struct {
	Location string `xml:"location,attr"`
}

type WsdlBinding struct {
	Name        string           `xml:"name,attr"`
	Type        string           `xml:"type,attr"`
	Doc         string           `xml:"documentation"`
	SoapBinding WsdlSoapBinding  `xml:"http://schemas.xmlsoap.org/wsdl/soap/ binding"`
	Operations  []*WsdlOperation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}

type WsdlPort struct {
	Name        string          `xml:"name,attr"`
	Binding     string          `xml:"binding,attr"`
	Doc         string          `xml:"documentation"`
	SoapAddress WsdlSoapAddress `xml:"http://schemas.xmlsoap.org/wsdl/soap/ address"`
}

type WsdlService struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"documentation"`
	Ports []*WsdlPort `xml:"http://schemas.xmlsoap.org/wsdl/ port"`
}
