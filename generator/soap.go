package generator

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"gopkg.in/inconshreveable/log15.v2"
)

var Log = log15.New()

func init() {
	Log.SetHandler(log15.DiscardHandler())
}

type SoapEnvelope struct {
	XMLName       xml.Name   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	EncodingStyle string     `xml:"http://schemas.xmlsoap.org/soap/encoding/ encodingStyle,attr"`
	Header        SoapHeader `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header,omitempty"`
	Body          SoapBody   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type SoapHeader struct {
	Header interface{}
}

type SoapBody struct {
	Body  string     `xml:",innerxml"`
	Fault *SoapFault `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
}

type SoapFault struct {
	faultcode   string `xml:"http://schemas.xmlsoap.org/soap/envelope/ faultcode"`
	faultstring string `xml:"faultstring"`
	faultactor  string `xml:"faultactor"`
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

func (s *SoapClient) Call(soapAction string, request, response interface{}) error {
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
	if err == nil {
		err = encoder.Flush()
	}
	Log.Debug("encoded", "envelope", log15.Lazy{func() string { return buffer.String() }})
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
	if len(body) == 0 {
		Log.Warn("empty response")
		return nil
	}

	respEnvelope := &SoapEnvelope{}

	err = xml.Unmarshal(body, respEnvelope)
	if err != nil {
		return err
	}

	if respEnvelope.Body.Body == "" {
		Log.Warn("empty response body", "envelope", respEnvelope, "body", string(body))
		return nil
	}

	err = xml.Unmarshal([]byte(respEnvelope.Body.Body), response)
	if err != nil {
		return err
	}

	return nil
}
