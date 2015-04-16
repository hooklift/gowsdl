// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
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
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	//Header SoapHeader `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header,omitempty"`
	Body SoapBody `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type SoapHeader struct {
	Header interface{}
}

type SoapBody struct {
	Fault   *SoapFault `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
	Content string     `xml:",innerxml"`
}

type SoapFault struct {
	Faultcode   string `xml:"faultcode,omitempty"`
	Faultstring string `xml:"faultstring,omitempty"`
	Faultactor  string `xml:"faultactor,omitempty"`
	Detail      string `xml:"detail,omitempty"`
}

type SoapClient struct {
	url string
	tls bool
}

func (f *SoapFault) Error() string {
	return f.Faultstring
}

func NewSoapClient(url string, tls bool) *SoapClient {
	return &SoapClient{
		url: url,
		tls: tls,
	}
}

func (s *SoapClient) Call(soapAction string, request, response interface{}) error {
	envelope := SoapEnvelope{
	//Header:        SoapHeader{},
	}

	if request != nil {
		reqXml, err := xml.Marshal(request)
		if err != nil {
			return err
		}

		envelope.Body.Content = string(reqXml)
	}
	buffer := &bytes.Buffer{}

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	err := encoder.Encode(envelope)
	if err == nil {
		err = encoder.Flush()
	}
	Log.Debug("request", "envelope", log15.Lazy{func() string { return buffer.String() }})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url, buffer)
	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}
	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

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

	rawbody, err := ioutil.ReadAll(res.Body)
	if len(rawbody) == 0 {
		Log.Warn("empty response")
		return nil
	}

	respEnvelope := &SoapEnvelope{}

	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	body := respEnvelope.Body.Content
	fault := respEnvelope.Body.Fault
	if body == "" {
		Log.Warn("empty response body", "envelope", respEnvelope, "body", body)
		return nil
	}

	Log.Debug("response", "envelope", respEnvelope, "body", body)
	if fault != nil {
		return fault
	}

	err = xml.Unmarshal([]byte(body), response)
	if err != nil {
		return err
	}

	return nil
}
