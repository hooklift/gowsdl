package gowsdl

var serverTmpl = `

var WSDLUndefinedError = errors.New("Server was unable to process request. --> Object reference not set to an instance of an object.")

type SOAPEnvelopeRequest struct {
	XMLName xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"` + "`" + `
	Body SOAPBodyRequest
}

type SOAPBodyRequest struct {
	XMLName xml.Name ` + "`" + `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"` + "`" + `
	{{range .}}
		{{range .Operations}}
				{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}} ` + `
  				{{$requestType}} *{{$requestType}} ` + "`" + `xml:,omitempty` + "`" + `
		{{end}}
	{{end}}
}

type SOAPEnvelopeResponse struct { ` + `
	XMLName    xml.Name` + "`" + `xml:"soap:Envelope"` + "`" + `
	PrefixSoap string  ` + "`" + `xml:"xmlns:soap,attr"` + "`" + `
	PrefixXsi  string  ` + "`" + `xml:"xmlns:xsi,attr"` + "`" + `
	PrefixXsd  string  ` + "`" + `xml:"xmlns:xsd,attr"` + "`" + `

	Body SOAPBodyResponse
}

func NewSOAPEnvelopResponse() *SOAPEnvelopeResponse {
	return &SOAPEnvelopeResponse{
		PrefixSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		PrefixXsd:  "http://www.w3.org/2001/XMLSchema",
		PrefixXsi:  "http://www.w3.org/2001/XMLSchema-instance",
	}
}

type Fault struct { ` + `
	XMLName xml.Name ` + "`" + `xml:"SOAP-ENV:Fault"` + "`" + `
	Space   string   ` + "`" + `xml:"xmlns:SOAP-ENV,omitempty,attr"` + "`" + `

	Code   string    ` + "`" + `xml:"faultcode,omitempty"` + "`" + `
	String string    ` + "`" + `xml:"faultstring,omitempty"` + "`" + `
	Actor  string 	 ` + "`" + `xml:"faultactor,omitempty"` + "`" + `
	Detail string    ` + "`" + `xml:"detail,omitempty"` + "`" + `
}


type SOAPBodyResponse struct { ` + `
	XMLName xml.Name   ` + "`" + `xml:"soap:Body"` + "`" + `
	Fault   *Fault ` + "`" + `xml:",omitempty"` + "`" + `
{{range .}}
	{{range .Operations}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}} ` + `
			{{$requestType}} *{{$responseType}} ` + "`" + `xml:",omitempty"` + "`" + `
	{{end}}
{{end}}

}

{{range .}}
	{{range .Operations}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$requestTypeSource := findType .Input.Message | replaceReservedWords }}
func (service *SOAPBodyRequest) {{$requestType}}Func(request *{{$requestType}}) (*{{$responseType}}, error) {
	return nil, WSDLUndefinedError
}
	{{end}}
{{end}}


func (service *SOAPEnvelopeRequest) call(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/xml; charset=utf-8")
	val := reflect.ValueOf(&service.Body).Elem()
	n := val.NumField()
	var field reflect.Value
	var name string
	find := false

	if r.Method == http.MethodGet {
		w.Write([]byte(wsdl))
		return
	}

	resp := NewSOAPEnvelopResponse()
	defer func() {
		if r := recover(); r != nil {
			resp.Body.Fault = &Fault{}
			resp.Body.Fault.Space = "http://schemas.xmlsoap.org/soap/envelope/"
			resp.Body.Fault.Code = "soap:Server"
			resp.Body.Fault.Detail = fmt.Sprintf("%v", r)
			resp.Body.Fault.String = fmt.Sprintf("%v", r)
		}
		xml.NewEncoder(w).Encode(resp)
	}()

	header := r.Header.Get("Content-Type")
	if strings.Index(header, "application/soap+xml") >= 0 {
		panic("Could not find an appropriate Transport Binding to invoke.")
	}

	err := xml.NewDecoder(r.Body).Decode(service)
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		field = val.Field(i)
		name = val.Type().Field(i).Name
		if field.Kind() != reflect.Ptr {
			continue
		}
		if field.IsNil() {
			continue
		}
		if field.IsValid() {
			find = true
			break
		}
	}

	if !find {
		panic(WSDLUndefinedError)
	} else {
		m := val.Addr().MethodByName(name + "Func")
		if !m.IsValid() {
			panic(WSDLUndefinedError)
		}

		vals := m.Call([]reflect.Value{field})
		if vals[1].IsNil() {
			reflect.ValueOf(&resp.Body).Elem().FieldByName(name).Set(vals[0])
		} else {
			panic(vals[1].Interface())
		}
	}

}

func Endpoint(w http.ResponseWriter, r *http.Request) {
	request := SOAPEnvelopeRequest{}
	request.call(w, r)
}

`
