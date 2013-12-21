package main

type Wsdl struct {
	Name            string          `xml:"name,attr"`
	TargetNamespace string          `xml:"targetNamespace,attr"`
	Imports         []*WsdlImport   `xml:"http://schemas.xmlsoap.org/wsdl/ import"`
	Doc             string          `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
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
	Doc    string       `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	Schema []*XsdSchema `xml:"http://www.w3.org/2001/XMLSchema schema"`
}

type WsdlPart struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}

type WsdlMessage struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	Parts []*WsdlPart `xml:"http://schemas.xmlsoap.org/wsdl/ part"`
}

type WsdlFault struct {
	Name      string        `xml:"name,attr"`
	Message   string        `xml:"message,attr"`
	Doc       string        `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	SoapFault WsdlSoapFault `xml:"http://schemas.xmlsoap.org/wsdl/soap/ fault"`
}

type WsdlInput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	SoapBody   WsdlSoapBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SoapHeader []*WsdlSoapHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
	//Mime       MimeBinding
}

type WsdlOutput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	SoapBody   WsdlSoapBody      `xml:"http://schemas.xmlsoap.org/wsdl/soap/ body"`
	SoapHeader []*WsdlSoapHeader `xml:"http://schemas.xmlsoap.org/wsdl/soap/ header"`
	//Mime       MimeBinding
}

type WsdlOperation struct {
	Name          string            `xml:"name,attr"`
	Doc           string            `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	Input         WsdlInput         `xml:"http://schemas.xmlsoap.org/wsdl/ input"`
	Output        WsdlOutput        `xml:"http://schemas.xmlsoap.org/wsdl/ output"`
	Faults        []*WsdlFault      `xml:"http://schemas.xmlsoap.org/wsdl/ fault"`
	SoapOperation WsdlSoapOperation `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
	HttpOperation WsdlHttpOperation `xml:"http://schemas.xmlsoap.org/wsdl/http/ operation"`
}

type WsdlPortType struct {
	Name       string           `xml:"name,attr"`
	Doc        string           `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	Operations []*WsdlOperation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}

type WsdlHttpBinding struct {
	Verb string `xml:"verb,attr"`
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
	HeadersFault  []*WsdlSoapHeaderFault `xml:"http://schemas.xmlsoap.org/wsdl/soap/ headerfault"`
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

type WsdlHttpOperation struct {
	Location string `xml:"location,attr"`
}

type WsdlBinding struct {
	Name        string           `xml:"name,attr"`
	Type        string           `xml:"type,attr"`
	Doc         string           `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	HttpBinding WsdlHttpBinding  `xml:"http://schemas.xmlsoap.org/wsdl/http/ binding"`
	SoapBinding WsdlSoapBinding  `xml:"http://schemas.xmlsoap.org/wsdl/soap/ binding"`
	Operations  []*WsdlOperation `xml:"http://schemas.xmlsoap.org/wsdl/ operation"`
}

type WsdlPort struct {
	Name        string          `xml:"name,attr"`
	Binding     string          `xml:"binding,attr"`
	Doc         string          `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	SoapAddress WsdlSoapAddress `xml:"http://schemas.xmlsoap.org/wsdl/soap/ address"`
}

type WsdlService struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"http://schemas.xmlsoap.org/wsdl/ documentation"`
	Ports []*WsdlPort `xml:"http://schemas.xmlsoap.org/wsdl/ port"`
}
