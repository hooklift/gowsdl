package generator

type Wsdl struct {
	Name            string          `xml:"name,attr"`
	TargetNamespace string          `xml:"targetNamespace,attr"`
	Imports         []*WsdlImport   `xml:"import"`
	Doc             string          `xml:"documentation"`
	Types           WsdlType        `xml:"types"`
	Messages        []*WsdlMessage  `xml:"message"`
	PortTypes       []*WsdlPortType `xml:"portType"`
	Binding         []*WsdlBinding  `xml:"binding"`
	Service         []*WsdlService  `xml:"service"`
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
	Parts []*WsdlPart `xml:"part"`
}

type WsdlFault struct {
	Name      string        `xml:"name,attr"`
	Message   string        `xml:"message,attr"`
	Doc       string        `xml:"documentation"`
	SoapFault WsdlSoapFault `xml:"fault"`
}

type WsdlInput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"documentation"`
	SoapBody   WsdlSoapBody      `xml:"body"`
	SoapHeader []*WsdlSoapHeader `xml:"header"`
	//Mime       MimeBinding
}

type WsdlOutput struct {
	Name       string            `xml:"name,attr"`
	Message    string            `xml:"message,attr"`
	Doc        string            `xml:"documentation"`
	SoapBody   WsdlSoapBody      `xml:"body"`
	SoapHeader []*WsdlSoapHeader `xml:"header"`
	//Mime       MimeBinding
}

type WsdlOperation struct {
	Name          string            `xml:"name,attr"`
	Doc           string            `xml:"documentation"`
	Input         WsdlInput         `xml:"input"`
	Output        WsdlOutput        `xml:"output"`
	Faults        []*WsdlFault      `xml:"fault"`
	SoapOperation WsdlSoapOperation `xml:"http://schemas.xmlsoap.org/wsdl/soap/ operation"`
	HttpOperation WsdlHttpOperation `xml:"http://schemas.xmlsoap.org/wsdl/http/ operation"`
}

type WsdlPortType struct {
	Name       string           `xml:"name,attr"`
	Doc        string           `xml:"documentation"`
	Operations []*WsdlOperation `xml:"operation"`
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

type WsdlHttpOperation struct {
	Location string `xml:"location,attr"`
}

type WsdlBinding struct {
	Name        string           `xml:"name,attr"`
	Type        string           `xml:"type,attr"`
	Doc         string           `xml:"documentation"`
	HttpBinding WsdlHttpBinding  `xml:"http://schemas.xmlsoap.org/wsdl/http/ binding"`
	SoapBinding WsdlSoapBinding  `xml:"http://schemas.xmlsoap.org/wsdl/soap/ binding"`
	Operations  []*WsdlOperation `xml:"operation"`
}

type WsdlPort struct {
	Name        string          `xml:"name,attr"`
	Binding     string          `xml:"binding,attr"`
	Doc         string          `xml:"documentation"`
	SoapAddress WsdlSoapAddress `xml:"address"`
}

type WsdlService struct {
	Name  string      `xml:"name,attr"`
	Doc   string      `xml:"documentation"`
	Ports []*WsdlPort `xml:"port"`
}
