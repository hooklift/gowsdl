// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var soapTmpl = `
var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"` + "`" + `
	Body SOAPBody ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"` + "`" + `
}

type SOAPHeader struct {
	Header interface{}
}

type SOAPBody struct {
	Fault   *SOAPFault ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"` + "`" + `
	Content string     ` + "`" + `xml:",innerxml"` + "`" + `
}

type SOAPFault struct {
	Code   string ` + "`" + `xml:"faultcode,omitempty"` + "`" + `
	String string ` + "`" + `xml:"faultstring,omitempty"` + "`" + `
	Actor  string ` + "`" + `xml:"faultactor,omitempty"` + "`" + `
	Detail string ` + "`" + `xml:"detail,omitempty"` + "`" + `
}

type BasicAuth struct {
	Login string
	Password string
}

type SOAPClient struct {
	url string
	tls bool
	auth *BasicAuth
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, tls bool, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url: url,
		tls: tls,
		auth: auth,
	}
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{
	//Header:        SoapHeader{},
	}

	if request != nil {
		reqXml, err := xml.Marshal(request)
		if err != nil {
			return err
		}

		envelope.Body.Content = string(reqXml)
	}
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	err := encoder.Encode(envelope)
	if err == nil {
		err = encoder.Flush()
	}

	log.Println(buffer.String())
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url, buffer)
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

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
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	body := respEnvelope.Body.Content
	fault := respEnvelope.Body.Fault
	if body == "" {
		log.Println("empty response body", "envelope", respEnvelope, "body", body)
		return nil
	}

	log.Println("response", "envelope", respEnvelope, "body", body)
	if fault != nil {
		return fault
	}

	err = xml.Unmarshal([]byte(body), response)
	if err != nil {
		return err
	}

	return nil
}
`
