// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

const XmlNameSpaceAttrAndEmpty = "`xml:\"xmlns:xsi,attr,omitempty\"`"
const XmlPlatformCoreAttrAndEmpty = "`xml:\"xmlns:platformCore,attr,omitempty\"`"

var typesTmpl = `
{{define "SimpleType"}}
	{{$typeName := replaceReservedWords .Name | makePublic}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
	{{else if ne .Union.MemberTypes ""}}
		type {{$typeName}} string
	{{else if .Union.SimpleType}}
		type {{$typeName}} string
	{{else if .Restriction.Base}}
		type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
    {{else}}
		type {{$typeName}} interface{}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base false}}
	{{ if $baseType }}
			{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Elements" .Extension.Choice}}
	{{template "Elements" .Extension.SequenceChoice}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
    {{ $targetNamespace := getNS }}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ normalize .Name | makeFieldPublic}} {{toGoType .Type false}} ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ else }}
			{{ normalize .Name | makeFieldPublic}} string ` + "`" + `xml:"{{with $targetNamespace}}{{.}} {{end}}{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base false}} ` + "`xml:\",chardata\" json:\"-,\"`" + `
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Ref .Nillable }} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty" json:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ normalize .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ normalize .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
				{{template "ComplexTypeInline" .}}
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceAttrReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Type .Nillable }} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{define "Any"}}
	{{range .}}
		Items     []string ` + "`" + `xml:",any" json:"items,omitempty"` + "`" + `
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := setNS .TargetNamespace }}

	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{$name := .Name}}
		{{$typeName := replaceReservedWords $name | makePublic}}
		{{if not .Type}}
			{{/* ComplexTypeLocal */}}
			{{with .ComplexType}}
				type {{$typeName}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
						{{template "Any" .Any}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
					{{end}}
				}
			{{end}}
			{{/* SimpleTypeLocal */}}
			{{with .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .List.ItemType ""}}
					type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
				{{else if ne .Union.MemberTypes ""}}
					type {{$typeName}} string
				{{else if .Union.SimpleType}}
					type {{$typeName}} string
				{{else if .Restriction.Base}}
					type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
				{{else}}
					type {{$typeName}} interface{}
				{{end}}
			
				{{if .Restriction.Enumeration}}
				const (
					{{with .Restriction}}
						{{range .Enumeration}}
							{{if .Doc}} {{.Doc | comment}} {{end}}
							{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
					{{end}}
				)
				{{end}}
			{{end}}
		{{else}}
			{{$type := toGoType .Type .Nillable | removePointerFromType}}
			{{if ne ($typeName) ($type)}}
				type {{$typeName}} {{$type}}
				{{if eq ($type) ("soap.XSDDateTime")}}
					func (xdt {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDDateTime(xdt).MarshalXML(e, start)
					}

					func (xdt *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDDateTime)(xdt).UnmarshalXML(d, start)
					}
				{{else if eq ($type) ("soap.XSDDate")}}
					func (xd {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDDate(xd).MarshalXML(e, start)
					}

					func (xd *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDDate)(xd).UnmarshalXML(d, start)
					}
				{{else if eq ($type) ("soap.XSDTime")}}
					func (xt {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDTime(xt).MarshalXML(e, start)
					}

					func (xt *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDTime)(xt).UnmarshalXML(d, start)
					}
				{{end}}
			{{end}}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$typeName := replaceReservedWords .Name | makePublic}}
		type {{$typeName}}_AddRequest struct {
			/* Here */
			XMLName  xml.Name ` + "`xml:\"add\"`" + `
			XmlNSXSI string ` + XmlNameSpaceAttrAndEmpty + `
			XmlNSPC  string ` + XmlPlatformCoreAttrAndEmpty + `
			XmlNS1   string ` + "`xml:\"xmlns:ns1,attr,omitempty\"`" + `
			Record   *{{$typeName}} ` + "`xml:\"record,omitempty\" json:\"record,omitempty\"`" + `
		}

		type {{$typeName}}_AddResponse struct {
			XMLName xml.Name ` + "`xml:\"addResponse\"`" + `
			WriteResponse *{{$typeName}}_WriteResponse ` + "`xml:\"writeResponse,omitempty\" json:\"writeResponse,omitempty\"`" + `
		}
		
		type {{$typeName}}_WriteResponse struct {
			XMLName xml.Name ` + "`xml:\"writeResponse\"`" + `
			Status  *Status  ` + "`xml:\"status,omitempty\" json:\"status,omitempty\"`" + `
			BaseRef *BaseRef ` + "`xml:\"baseRef,omitempty\" json:\"baseRef,omitempty\"`" + `
		}
		
		//----------------------------------------------------
		
		type {{$typeName}}_GetRequest struct {
			XMLName  xml.Name ` + "`xml:\"get\"`" + `
			XmlNSXSI string   ` + XmlNameSpaceAttrAndEmpty + `
			XmlNSPC  string   ` + XmlPlatformCoreAttrAndEmpty + `
			BaseRef  *BaseRef ` + "`xml:\"baseRef,omitempty\" json:\"baseRef,omitempty\"`" + `
		}
		
		type {{$typeName}}_GetResponse struct {
			//XMLName xml.Name ` + "`xml:\"urn:messages_2022_1.platform.webservices.netsuite.com getResponse\"`" + `
			XMLName xml.Name ` + "`xml:\"getResponse\"`" + `
			ReadResponse *{{$typeName}}_ReadResponse ` + "`xml:\"readResponse,omitempty\" json:\"readResponse,omitempty\"`" + `
		}
		
		//----------------------------------------------------
		
		type {{$typeName}}_ReadResponse struct {
			XMLName xml.Name  ` + "`xml:\"readResponse\"`" + `
			Status  *Status   ` + "`xml:\"status,omitempty\" json:\"status,omitempty\"`" + `
			Record  *{{$typeName}} ` + "`xml:\"record,omitempty\" json:\"record,omitempty\"`" + `
		}		

	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$typeName := replaceReservedWords .Name | makePublic}}
		{{if and (eq (len .SimpleContent.Extension.Attributes) 0) (eq (toGoType .SimpleContent.Extension.Base false) "string") }}
			type {{$typeName}} string
		{{else}}
			type {{$typeName}} struct {
				/* Here Jay */
				XsiType string ` + "`xml:\"xsi:type,attr,omitempty\"  json:\"-\"`" + `
				{{if eq $typeName "BaseRef"}}
					InternalId string ` + "`xml:\"internalId,attr,omitempty\" json:\"internalId,omitempty\"`" + `
					ExternalId string ` + "`xml:\"externalId,attr,omitempty\" json:\"externalId,omitempty\"`" + `
					Type *RecordType ` + "`xml:\"type,attr,omitempty\" json:\"type,omitempty\"`" + `
				{{end}}
				{{if eq $typeName "DeleteRequest"}}
					XmlNSXSI string ` + XmlNameSpaceAttrAndEmpty + `
					XmlNSPC string ` + XmlPlatformCoreAttrAndEmpty + `
				{{end}}
				{{if eq $typeName "DeleteResponse"}}
					DeleteReturn *DeleteReturn ` + "`xml:\"deleteReturn,omitempty\" json:\"writeResponse,omitempty\"`" + `
				{{end}}
				{{$type := findNameByType .Name}}
				{{if ne .Name $type}}
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$type}}\"`" + `
				{{end}}

				{{if ne .ComplexContent.Extension.Base ""}}
					/* Jay 2 */
					{{template "ComplexContent" .ComplexContent}}
				{{else if ne .SimpleContent.Extension.Base ""}}
					/* Jay 3 */
				{{template "SimpleContent" .SimpleContent}}
				{{else}}
					/* Jay 4 */
					{{if ne $typeName "DeleteResponse"}}
						{{template "Elements" .Sequence}}
					{{end}}
					{{template "Any" .Any}}
					{{template "Elements" .Choice}}
					{{template "Elements" .SequenceChoice}}
					{{template "Elements" .All}}
					{{template "Attributes" .Attributes}}
				{{end}}
			}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* netSuitePortFuncs */}}
		{{$typeName := replaceReservedWords .Name | makePublic}}

		func (service *netSuitePortType) {{$typeName}}_AddContext(ctx context.Context, request *{{$typeName}}_AddRequest) (*{{$typeName}}_AddResponse, error) {
			response := new({{$typeName}}_AddResponse)
			err := service.client.CallContext(ctx, "add", request, response)
			if err != nil {
				return nil, err
			}
		
			return response, nil
		}


		func (service *netSuitePortType) {{$typeName}}_Add(request *{{$typeName}}_AddRequest) (*{{$typeName}}_AddResponse, error) {
			return service.{{$typeName}}_AddContext(
				context.Background(),
				request,
			)
		}

		func (service *netSuitePortType) {{$typeName}}_GetContext(ctx context.Context, request *{{$typeName}}_GetRequest) (*{{$typeName}}_GetResponse, error) {
			response := new({{$typeName}}_GetResponse)
			err := service.client.CallContext(ctx, "get", request, response)
			if err != nil {
				return nil, err
			}
		
			return response, nil
		}
		
		func (service *netSuitePortType) {{$typeName}}_Get(request *{{$typeName}}_GetRequest) (*{{$typeName}}_GetResponse, error) {
			return service.{{$typeName}}_GetContext(
				context.Background(),
				request,
			)
		}
		
	{{end}}
{{end}}
`
