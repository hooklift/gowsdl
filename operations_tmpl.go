// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var opsTmpl = `
{{range .}}
	{{$privateType := .Name | makePrivate}}
	{{$exportType := .Name | makePublic}}

	type {{$exportType}} interface {
		{{range .Operations}}
			{{$faults := len .Faults}}
			{{$soapAction := findSOAPAction .Name $privateType}}
			{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
			{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}

			{{/*if ne $soapAction ""*/}}
			{{if gt $faults 0}}
			// Error can be either of the following types:
			// {{range .Faults}}
			//   - {{.Name}} {{.Doc}}{{end}}{{end}}
			{{if ne .Doc ""}}/* {{.Doc}} */{{end}}
			{{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) (*{{$responseType}}, error)
			{{/*end*/}}
		{{end}}
	}

	type {{$privateType}} struct {
		client *soap.Client
	}

	func New{{$exportType}}(client *soap.Client) {{$exportType}} {
		return &{{$privateType}}{
			client: client,
		}
	}

	{{range .Operations}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$soapAction := findSOAPAction .Name $privateType}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}
		func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) (*{{$responseType}}, error) {
			response := new({{$responseType}})
			err := service.client.Call("{{$soapAction}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, response)
			if err != nil {
				return nil, err
			}

			return response, nil
		}
	{{end}}
{{end}}
`
