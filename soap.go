package main

import (
	"encoding/xml"
)

type SoapEnvelope struct {
	XMLName       xml.Name   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	EncodingStyle string     `xml:"http://schemas.xmlsoap.org/soap/envelope/ encodingStyle,attr"`
	Header        SoapHeader `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body          SoapBody   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type SoapHeader struct {
	Header interface{}
}

type SoapBody struct {
	Body  interface{}
	Fault SoapFault
}

type SoapFault struct {
	Fault interface{}
}
