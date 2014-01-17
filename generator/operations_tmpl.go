package generator

var opsTmpl = `
{{range .}}
	{{$portType := .Name}}
	type {{$portType}} struct {
		client *gowsdl.SoapClient
	}

	func New{{$portType}}(url string, tls bool) *{{$portType}} {
		client := gowsdl.NewSoapClient(url, tls)

		return &{{$portType}}{
			client: client,
		}
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
			response := &{{$output}}{}
			err := service.client.Call("{{.Name}}", "{{$soapAction}}", request, response)
			if err != nil {
				return nil, err
			}

			return response, nil			
		}
	{{end}}
{{end}}
`
