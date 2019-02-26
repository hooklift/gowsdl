package soap

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// supported SOAP versions
const (
	SOAP11 = 11
	SOAP12 = 12
)

// SOAPVersion determines which version of SOAP protocol to use for requests
type SOAPVersion int

// SOAPEnvelope is general representation of 3 types supported XML SOAP Envelopes
// for SOAP 1.1, 1.2 and custom user version
type SOAPEnvelope interface {
	SetHeaders([]interface{})
	SetContent(content interface{})
	Fault() *SOAPFault
}

// newSOAPEnvelope produces new SOAPEnvelope struct according passed params
func newSOAPEnvelope(c *Client, request bool) (SOAPEnvelope, error) {
	// if user provided callbacks for custom constructors, use them
	if request && c.opts.customReq != nil {
		return c.opts.customReq(), nil
	}
	if !request && c.opts.customResp != nil {
		return c.opts.customResp(), nil
	}

	// overwise use default for standard headers
	switch c.opts.version {
	case SOAP11:
		return new(SOAPEnvelope11), nil
	case SOAP12:
		return new(SOAPEnvelope12), nil
	}

	return nil, errors.New("version is not supported")
}

type SOAPEnvelope11 struct {
	XMLName xml.Name      `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Headers []interface{} `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body    SOAPBody11
}

func (s *SOAPEnvelope11) SetHeaders(headers []interface{}) {
	s.Headers = headers
}

func (s *SOAPEnvelope11) SetContent(content interface{}) {
	s.Body = SOAPBody11{Content: content}
}

func (s SOAPEnvelope11) Fault() *SOAPFault {
	return s.Body.Fault
}

type SOAPEnvelope12 struct {
	XMLName xml.Name      `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Headers []interface{} `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	Body    SOAPBody12
}

func (s *SOAPEnvelope12) SetHeaders(headers []interface{}) {
	s.Headers = headers
}

func (s *SOAPEnvelope12) SetContent(content interface{}) {
	s.Body = SOAPBody12{Content: content}
}

func (s SOAPEnvelope12) Fault() *SOAPFault {
	return s.Body.Fault
}

type SOAPBody11 struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

func (s SOAPBody11) content() interface{} {
	return s.Content
}

func (s *SOAPBody11) setContent(c interface{}) {
	s.Content = c
}

func (s SOAPBody11) fault() *SOAPFault {
	return s.Fault
}

func (s *SOAPBody11) setFault(f *SOAPFault) {
	s.Fault = f
}

type SOAPBody interface {
	content() interface{}
	setContent(interface{})
	fault() *SOAPFault
	setFault(*SOAPFault)
}

type SOAPBody12 struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

func (s SOAPBody12) content() interface{} {
	return s.Content
}

func (s *SOAPBody12) setContent(c interface{}) {
	s.Content = c
}

func (s SOAPBody12) fault() *SOAPFault {
	return s.Fault
}

func (s *SOAPBody12) setFault(f *SOAPFault) {
	s.Fault = f
}

// UnmarshalXML unmarshals SOAPBody11 xml
func (b *SOAPBody11) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	return unmarshalSOAPBody(d, b)
}

// UnmarshalXML unmarshals SOAPBody12 xml
func (b *SOAPBody12) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	return unmarshalSOAPBody(d, b)
}

// unmarshalSOAPBody is a common code for both SOAP11 and SOAP12 unmarshal funcs
func unmarshalSOAPBody(d *xml.Decoder, b SOAPBody) error {
	if b.content() == nil {
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
				b.setFault(&SOAPFault{})
				b.setContent(nil)

				err = d.DecodeElement(b.fault(), &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.content(), &se); err != nil {
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

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

func (f *SOAPFault) Error() string {
	return f.String
}

const (
	// Predefined WSS namespaces to be used in
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName   xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	MustUnderstand string `xml:"mustUnderstand,attr,omitempty"`

	Token *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name `xml:"wsse:UsernameToken"`
	XmlNSWsu  string   `xml:"xmlns:wsu,attr"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Id string `xml:"wsu:Id,attr,omitempty"`

	Username *WSSUsername `xml:",omitempty"`
	Password *WSSPassword `xml:",omitempty"`
}

type WSSUsername struct {
	XMLName   xml.Name `xml:"wsse:Username"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Data string `xml:",chardata"`
}

type WSSPassword struct {
	XMLName   xml.Name `xml:"wsse:Password"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	XmlNSType string   `xml:"Type,attr"`

	Data string `xml:",chardata"`
}

// NewWSSSecurityHeader creates WSSSecurityHeader instance
func NewWSSSecurityHeader(user, pass, tokenID, mustUnderstand string) *WSSSecurityHeader {
	hdr := &WSSSecurityHeader{XmlNSWsse: WssNsWSSE, MustUnderstand: mustUnderstand}
	hdr.Token = &WSSUsernameToken{XmlNSWsu: WssNsWSU, XmlNSWsse: WssNsWSSE, Id: tokenID}
	hdr.Token.Username = &WSSUsername{XmlNSWsse: WssNsWSSE, Data: user}
	hdr.Token.Password = &WSSPassword{XmlNSWsse: WssNsWSSE, XmlNSType: WssNsType, Data: pass}
	return hdr
}

type basicAuth struct {
	Login    string
	Password string
}

type options struct {
	tlsCfg           *tls.Config
	auth             *basicAuth
	timeout          time.Duration
	contimeout       time.Duration
	tlshshaketimeout time.Duration
	client           HTTPClient
	httpHeaders      map[string]string
	version          SOAPVersion
	customReq        func() SOAPEnvelope
	customResp       func() SOAPEnvelope
}

var defaultOptions = options{
	timeout:          time.Duration(30 * time.Second),
	contimeout:       time.Duration(90 * time.Second),
	tlshshaketimeout: time.Duration(15 * time.Second),
	version:          SOAP11,
}

// A Option sets options such as credentials, tls, etc.
type Option func(*options)

// WithHTTPClient is an Option to set the HTTP client to use
// This cannot be used with WithTLSHandshakeTimeout, WithTLS,
// WithTimeout options
func WithHTTPClient(c HTTPClient) Option {
	return func(o *options) {
		o.client = c
	}
}

// WithTLSHandshakeTimeout is an Option to set default tls handshake timeout
// This option cannot be used with WithHTTPClient
func WithTLSHandshakeTimeout(t time.Duration) Option {
	return func(o *options) {
		o.tlshshaketimeout = t
	}
}

// WithRequestTimeout is an Option to set default end-end connection timeout
// This option cannot be used with WithHTTPClient
func WithRequestTimeout(t time.Duration) Option {
	return func(o *options) {
		o.contimeout = t
	}
}

// WithBasicAuth is an Option to set BasicAuth
func WithBasicAuth(login, password string) Option {
	return func(o *options) {
		o.auth = &basicAuth{Login: login, Password: password}
	}
}

// WithTLS is an Option to set tls config
// This option cannot be used with WithHTTPClient
func WithTLS(tls *tls.Config) Option {
	return func(o *options) {
		o.tlsCfg = tls
	}
}

// WithTimeout is an Option to set default HTTP dial timeout
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.timeout = t
	}
}

// WithHTTPHeaders is an Option to set global HTTP headers for all requests
func WithHTTPHeaders(headers map[string]string) Option {
	return func(o *options) {
		o.httpHeaders = headers
	}
}

// WithSOAPVersion is an Option to set SOAP protocol version to use
func WithSOAPVersion(version SOAPVersion) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithCustomRequester is an Option to set specific SOAP Envelope contructor for
// broken services (used only for send requests)
func WithCustomRequester(req func() SOAPEnvelope) Option {
	return func(o *options) {
		o.customReq = req
	}
}

// Client is soap client
type Client struct {
	url     string
	opts    *options
	headers []interface{}
}

// HTTPClient is a client which can make HTTP requests
// An example implementation is net/http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient creates new SOAP client instance
func NewClient(url string, opt ...Option) *Client {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return &Client{
		url:  url,
		opts: &opts,
	}
}

// AddHeader adds envelope header
func (s *Client) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

// Call performs HTTP POST request
func (s *Client) Call(soapAction string, request, response interface{}) error {
	envelope, err := newSOAPEnvelope(s, true /* request */)
	if err != nil {
		return nil
	}

	if s.headers != nil && len(s.headers) > 0 {
		envelope.SetHeaders(s.headers)
	}

	envelope.SetContent(request)
	buffer := new(bytes.Buffer)
	encoder := xml.NewEncoder(buffer)

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.opts.auth != nil {
		req.SetBasicAuth(s.opts.auth.Login, s.opts.auth.Password)
	}

	switch s.opts.version {
	case SOAP11:
		req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	case SOAP12:
		req.Header.Add("Content-Type", "application/soap+xml")
	}
	req.Header.Add("SOAPAction", soapAction)
	req.Header.Set("User-Agent", "gowsdl/0.1")
	if s.opts.httpHeaders != nil {
		for k, v := range s.opts.httpHeaders {
			req.Header.Set(k, v)
		}
	}
	req.Close = true

	client := s.opts.client
	if client == nil {
		tr := &http.Transport{
			TLSClientConfig: s.opts.tlsCfg,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, s.opts.timeout)
			},
			TLSHandshakeTimeout: s.opts.tlshshaketimeout,
		}
		client = &http.Client{Timeout: s.opts.contimeout, Transport: tr}
	}

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
		return nil
	}

	respEnvelope, err := newSOAPEnvelope(s, false /* request */)
	if err != nil {
		return nil
	}
	respEnvelope.SetContent(response)
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Fault()
	if fault != nil {
		return fault
	}

	return nil
}
