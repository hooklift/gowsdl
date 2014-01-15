package generator

var opsTmpl = `
var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

{{range .}}
	{{$portType := .Name}}
	type {{$portType}} struct {
		url string
		tls bool
	}

	func New{{$portType}}(url string, tls bool) *{{$portType}} {
		service := &{{$portType}}{
			url: url,
			tls: tls,
		}

		return service
	}

	func (service *{{$portType}}) call(operation, soapAction string, request interface{}) ([]byte, error) {
		envelope := gowsdl.SoapEnvelope{
			Header:        gowsdl.SoapHeader{},
			EncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		}

		reqXml, err := xml.Marshal(request)
		if err != nil {
			return nil, err
		}

		envelope.Body = gowsdl.SoapBody{
			Body: string(reqXml),
		}

		buffer := &bytes.Buffer{}

		encoder := xml.NewEncoder(buffer)
		//encoder.Indent("  ", "    ")

		err = encoder.Encode(envelope)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", service.url, buffer)
		req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
		req.Header.Add("SOAPAction", soapAction)
		req.Header.Set("User-Agent", "gowsdl/0.1")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: service.tls,
			},
			Dial: dialTimeout,
		}

		client := &http.Client{Transport: tr}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		return body, nil		
	}

	{{range .Operations}}
		{{$requestType := findType .Input.Message}}
		{{$soapAction := findSoapAction .Name $portType}}
		{{$output := findType .Output.Message}}
		/**
		* Error can be either of the following types:
		* {{range .Faults}}
		* - {{.Name}} {{.Doc}}{{end}}
		*/
		func (service *{{$portType}}) {{makePublic .Name}} (request *{{$requestType}}) (*{{$output}}, error) {
			data, err := service.call("{{.Name}}", "{{$soapAction}}", request)
			if err != nil {
				return nil, err
			}

			envelope := &gowsdl.SoapEnvelope{}

			err = xml.Unmarshal(data, envelope)
			if err != nil {
				return nil, err
			}

			if envelope.Body.Body == "" {
				log.Printf("%#v\n", envelope.Body)
				return nil, nil
			}

			res := &{{$output}}{}
			err = xml.Unmarshal([]byte(envelope.Body.Body), res)
			if err != nil {
				return nil, err
			}

			return res, nil
		}
	{{end}}
{{end}}
`
