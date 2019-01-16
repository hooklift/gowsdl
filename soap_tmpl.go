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
	Header *SOAPHeader
	Body SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"` + "`" + `

	Items []interface{} ` + "`" + `xml:",omitempty"` + "`" + `
}

type SOAPBody struct {
	XMLName xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"` + "`" + `

	Fault   *SOAPFault ` + "`" + `xml:",omitempty"` + "`" + `
	Content interface{} ` + "`" + `xml:",omitempty"` + "`" + `
}

type SOAPFault struct {
	XMLName xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"` + "`" + `

	Code   string ` + "`" + `xml:"faultcode,omitempty"` + "`" + `
	String string ` + "`" + `xml:"faultstring,omitempty"` + "`" + `
	Actor  string ` + "`" + `xml:"faultactor,omitempty"` + "`" + `
	Detail string ` + "`" + `xml:"detail,omitempty"` + "`" + `
}

const (
	// Predefined WSS namespaces to be used in
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName   xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"` + "`" + `
	XmlNSWsse string   ` + "`" + `xml:"xmlns:wsse,attr"` + "`" + `

	MustUnderstand string ` + "`" + `xml:"mustUnderstand,attr,omitempty"` + "`" + `

	Token *WSSUsernameToken ` + "`" + `xml:",omitempty"` + "`" + `
}

type WSSUsernameToken struct {
	XMLName   xml.Name ` + "`" + `xml:"wsse:UsernameToken"` + "`" + `
	XmlNSWsu  string   ` + "`" + `xml:"xmlns:wsu,attr"` + "`" + `
	XmlNSWsse string   ` + "`" + `xml:"xmlns:wsse,attr"` + "`" + `

	Id string ` + "`" + `xml:"wsu:Id,attr,omitempty"` + "`" + `

	Username *WSSUsername ` + "`" + `xml:",omitempty"` + "`" + `
	Password *WSSPassword ` + "`" + `xml:",omitempty"` + "`" + `
}

type WSSUsername struct {
	XMLName   xml.Name ` + "`" + `xml:"wsse:Username"` + "`" + `
	XmlNSWsse string   ` + "`" + `xml:"xmlns:wsse,attr"` + "`" + `

	Data string ` + "`" + `xml:",chardata"` + "`" + `
}

type WSSPassword struct {
	XMLName   xml.Name ` + "`" + `xml:"wsse:Password"` + "`" + `
	XmlNSWsse string   ` + "`" + `xml:"xmlns:wsse,attr"` + "`" + `
	XmlNSType string   ` + "`" + `xml:"Type,attr"` + "`" + `

	Data string ` + "`" + `xml:",chardata"` + "`" + `
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url     string
	tlsCfg  *tls.Config
	auth    *BasicAuth
	headers []interface{}
}

// **********
// Accepted solution from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// Author: Icza - http://stackoverflow.com/users/1705598/icza

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// **********

func NewWSSSecurityHeader(user, pass, mustUnderstand string) *WSSSecurityHeader {
	hdr := &WSSSecurityHeader{XmlNSWsse: WssNsWSSE, MustUnderstand: mustUnderstand}
	hdr.Token = &WSSUsernameToken{XmlNSWsu: WssNsWSU, XmlNSWsse: WssNsWSSE, Id: "UsernameToken-" + randStringBytesMaskImprSrc(9)}
	hdr.Token.Username = &WSSUsername{XmlNSWsse: WssNsWSSE, Data: user}
	hdr.Token.Password = &WSSPassword{XmlNSWsse: WssNsWSSE, XmlNSType: WssNsType, Data: pass}
	return hdr
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, insecureSkipVerify bool, auth *BasicAuth) *SOAPClient {
	tlsCfg := &tls.Config{
	       InsecureSkipVerify: insecureSkipVerify,
	}
	return NewSOAPClientWithTLSConfig(url, tlsCfg, auth)
}

func NewSOAPClientWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url: url,
		tlsCfg: tlsCfg,
		auth: auth,
	}
}

func (s *SOAPClient) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{}

	if s.headers != nil && len(s.headers) > 0 {
		soapHeader := &SOAPHeader{Items: make([]interface{}, len(s.headers))}
		copy(soapHeader.Items, s.headers)
		envelope.Header = soapHeader
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println(buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: s.tlsCfg,
		Dial: dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
`
