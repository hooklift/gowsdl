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
			{{$responseAttachments := responseAttachments .Output.Message }}
			{{$attachments := len $responseAttachments}}

			{{/*if ne $soapAction ""*/}}
			{{if gt $faults 0}}
			// Error can be either of the following types:
			// {{range .Faults}}
			//   - {{.Name}} {{.Doc}}{{end}}{{end}}
			{{if ne .Doc ""}}/* {{.Doc}} */{{end}}
			{{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}{{if gt $attachments 0}}[]soap.MIMEMultipartAttachment, {{end}}error)
			{{/*end*/}}
			{{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}{{if gt $attachments 0}}[]soap.MIMEMultipartAttachment, {{end}}error)
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
		{{$responseAttachments := responseAttachments .Output.Message }}
		{{$attachments := len $responseAttachments}}

		{{if gt $attachments 0}}
			func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}[]soap.MIMEMultipartAttachment, error) {
				{{if ne $responseType ""}}response := new({{$responseType}}){{end}}
				attachments, err := service.client.CallContextWithAttachmentsAndFaultDetail(ctx, "{{if ne $soapAction ""}}{{$soapAction}}{{else}}''{{end}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, {{if ne $responseType ""}}response{{else}}struct{}{}{{end}}, nil)
				if err != nil {
					return {{if ne $responseType ""}}nil, {{end}}nil, err
				}

				return {{if ne $responseType ""}}response, {{end}}attachments, nil
			}

			func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}[]soap.MIMEMultipartAttachment, error) {
				return service.{{makePublic .Name | replaceReservedWords}}Context(
					context.Background(),
					{{if ne $requestType ""}}request,{{end}}
				)
			}
		{{else}}
			func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}}Context (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {
				{{if ne $responseType ""}}response := new({{$responseType}}){{end}}
				err := service.client.CallContext(ctx, "{{if ne $soapAction ""}}{{$soapAction}}{{else}}''{{end}}", {{if ne $requestType ""}}request{{else}}nil{{end}}, {{if ne $responseType ""}}response{{else}}struct{}{}{{end}})
				if err != nil {
					return {{if ne $responseType ""}}nil, {{end}}err
				}

				return {{if ne $responseType ""}}response, {{end}}nil
			}

			func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}} ({{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {
				return service.{{makePublic .Name | replaceReservedWords}}Context(
					context.Background(),
					{{if ne $requestType ""}}request,{{end}}
				)
			}
		{{end}}

	{{end}}
{{end}}
`
