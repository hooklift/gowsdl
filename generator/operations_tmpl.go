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
		{{$faults := len .Faults}}
		{{$requestType := findType .Input.Message}}
		{{$soapAction := findSoapAction .Name $portType}}
		{{$output := findType .Output.Message}}

		{{if gt $faults 0}}
		/**
		* Error can be either of the following types:
		* {{range .Faults}}
		* - {{.Name}} {{.Doc}}{{end}}
		*/
		{{end}}
		func (service *{{$portType}}) {{makePublic .Name}} (request *{{$requestType}}) (*{{$output}}, error) {
			response := &{{$output}}{}
			err := service.client.Call("{{$soapAction}}", request, response)
			if err != nil {
				return nil, err
			}

			return response, nil
		}
	{{end}}
{{end}}
`
