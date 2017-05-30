// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
{{define "MoreAnnotations"}}json:"{{.}},omitempty" yaml:"{{.}},omitempty"{{end}}

{{define "SimpleType"}}
	{{$type := replaceReservedWords .Name | makePublic}}
	type {{$type}} {{toGoType .Restriction.Base}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}} {{if not .Type}}
			{{ .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},attr,omitempty" {{template "MoreAnnotations" .Name}}` + "`" + `
		{{else}}
			{{ .Name | makeFieldPublic}} {{toGoType .Type}} ` + "`" + `xml:"{{.Name}},attr,omitempty" {{template "MoreAnnotations" .Name}}` + "`" + `
		{{end}}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base}}{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Name | makePublic}} struct {
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
	} ` + "`" + `xml:"{{.Name}},omitempty" {{template "MoreAnnotations" .Name}}` + "`" + `
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Ref | toGoType}} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty" {{template "MoreAnnotations" .Name}}` + "`" + `
		{{else}}
		{{if not .Type}}
			{{template "ComplexTypeInline" .}}
		{{else}}
			{{if .Doc}}
				{{.Doc | comment}} {{"\n"}}
			{{end}}
			{{replaceReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Type | toGoType}} ` + "`" + `xml:"{{.Name}},omitempty" {{template "MoreAnnotations" .Name}}` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := .TargetNamespace }}

	{{range .SimpleTypes}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{if not .Type}}
			{{/* ComplexTypeLocal */}}
			{{$name := .Name}}
			{{with .ComplexType}}
				type {{$name | replaceReservedWords | makePublic}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
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
				}
			{{end}}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$name := replaceReservedWords .Name | makePublic}}
		type {{$name}} struct {
			{{$typ := findNameByType .Name}}
			{{if ne $name $typ}}
				XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$typ}}\"`" + `
			{{end}}
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
		}
	{{end}}
{{end}}
`
