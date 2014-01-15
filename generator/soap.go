package generator

import (
	"encoding/xml"
)

type SoapEnvelope struct {
	XMLName       xml.Name   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	EncodingStyle string     `xml:"http://schemas.xmlsoap.org/soap/encoding/ encodingStyle,attr"`
	Header        SoapHeader `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body          SoapBody   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type SoapHeader struct {
	Header interface{}
}

type SoapBody struct {
	Body  string
	Fault SoapFault
}

type SoapFault struct {
	faultcode   string `xml:"http://schemas.xmlsoap.org/soap/envelope/ faultcode"`
	faultstring string `xml:"faultstring"`
	faulactor   string `xml:"faultactor"`
	detail      string `xml:"detail"`
}
