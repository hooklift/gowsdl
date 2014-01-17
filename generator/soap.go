package generator

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
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

type SoapClient struct {
	url string
	tls bool
}

func NewSoapClient(url string, tls bool) *SoapClient {
	return &SoapClient{
		url: url,
		tls: tls,
	}
}

func (s *SoapClient) Call(operation, soapAction string, request, response interface{}) error {
	envelope := SoapEnvelope{
		Header:        SoapHeader{},
		EncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
	}

	reqXml, err := xml.Marshal(request)
	if err != nil {
		return err
	}

	envelope.Body = SoapBody{
		Body: string(reqXml),
	}

	buffer := &bytes.Buffer{}

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	err = encoder.Encode(envelope)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url, buffer)
	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)
	req.Header.Set("User-Agent", "gowsdl/0.1")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: s.tls,
		},
		Dial: dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	respEnvelope := &SoapEnvelope{}

	err = xml.Unmarshal(body, respEnvelope)
	if err != nil {
		return err
	}

	if respEnvelope.Body.Body == "" {
		log.Printf("%#v\n", respEnvelope.Body)
		return nil
	}

	err = xml.Unmarshal([]byte(respEnvelope.Body.Body), response)
	if err != nil {
		return err
	}

	return nil
}
