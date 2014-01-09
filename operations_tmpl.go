package main

var opsTmpl = `
{{range .}}
	{{$portType := .Name}}
	type {{$portType}} struct {
		client *SoapClient
	}

	{{range .Operations}}
		func (service *{{$portType}}) {{.Name}} (request *{{findType .Input.Message}}) (*{{findType .Output.Message}}, error) {
			return nil, nil
		}
	{{end}}
{{end}}
`
