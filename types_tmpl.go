// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
  
	{{define "SimpleType"}}
		{{$type := replaceReservedWords .Name | makePublic}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{if ne .List.ItemType ""}}
			type {{$type}} []{{toGoType .List.ItemType }}
		{{else if ne .Union.MemberTypes ""}}
			type {{$type}} string
		{{else if .Union.SimpleType}}
			type {{$type}} string
		{{else}}
			type {{$type}} {{toGoTypeNoPointer .Restriction.Base}}
		{{end}}
        {{if .Restriction.SimpleType}} 
        {{template "SimpleType" .Restriction.SimpleType}}
        {{end}}

		{{if .Restriction.Enumeration}}
			const (
					{{with .Restriction}}
						{{range .Enumeration}}
						{{if .Doc}} {{.Doc | comment}} {{end}}
						{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}" {{end}}
					{{end}}
					)
		{{end}}
	{{end}}

	{{define "ComplexContent"}}
		{{$baseType := toGoType .Extension.Base}}
		{{ if $baseType }}
			{{$baseType}}
		{{end}}

		{{template "Elements" .Extension.SequenceChoice}}
        {{template "Elements" .Extension.SequenceChoiceSequence}} 
		{{template "Elements" .Extension.Sequence}}
		{{template "Attributes" .Extension.Attributes}}
		{{template "Elements" .Extension.Choice}}
		{{template "AttributeGroups" .Extension.AttributeGroup}}
	{{end}}

	{{define "Attributes"}}
		{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
			{{ if ne .Type "" }}
				{{ .Name | makeFieldPublic}} {{toGoType .Type}} ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
			{{ else }}
				{{ .Name | makeFieldPublic}} string ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
			{{ end }}
		{{end}}
	{{end}}
	{{define "AttributeGroups"}}
		{{range .}}
			{{if ne .Ref ""}}
				{{template "Attributes" getAttributesFromGroup .Ref}}
			{{ end }}
		{{end}}
	{{end}}

	{{define "SimpleContent"}}
		Value {{toGoTypeNoPointer .Extension.Base}} ` + "`" + `xml:",chardata"` + "`" + `
		{{template "Attributes" .Extension.Attributes}}
		{{template "AttributeGroups" .Extension.AttributeGroup}}
	{{end}}

	{{define "ComplexTypeInline"}}
		{{replaceReservedWords .Key | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
		{{with .ComplexType}}
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{else}}
				{{template "Elements" .Sequence}}
				{{template "Elements" .Choice}}
				{{template "Elements" .SequenceChoice}}
                {{template "Elements" .SequenceChoiceSequence}} 
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
				{{template "Elements" .ChoiceSequence}}
				{{template "AttributeGroups" .AttributeGroup}}
			{{end}}
		{{end}}
		} ` + "`" + `xml:"{{.Key}},omitempty"` + "`" + `
	{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Ref | toGoType}} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
                {{ $complexInlineName := setElementInComplexInlineCache .}}
                 {{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}*  {{$complexInlineName}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Type | toGoType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := .TargetNamespace }}

	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

   {{range .Elements}}
  		{{$name := .Name}}
		{{if not .Type}}
		{{/* ComplexTypeLocal */}}
         {{with .ComplexType}}
              type {{$name | replaceReservedWords | makePublic}} struct {  XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
                        {{template "Elements" .SequenceSequence}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
                        {{template "Elements" .SequenceChoiceSequence}} 
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
                        {{template "Elements" .ChoiceSequence}}
                        {{template "AttributeGroups" .AttributeGroup}}
					{{end}}
				}
			{{end}}
		{{else}}
			type {{$name | replaceReservedWords | makePublic}} {{toGoType .Type | removePointerFromType}}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		{{$name := replaceReservedWords .Name | makePublic}}
		type {{$name}} struct {
			{{$typ := findNameByType .Name}}
			{{if ne $name $typ}}
				
			{{end}}
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{else}}
				{{template "Elements" .Sequence}}
                {{template "Elements" .SequenceSequence}}
                {{template "Elements" .Choice}}
				{{template "Elements" .SequenceChoice}}
                {{template "Elements" .SequenceChoiceSequence}} 
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
                {{template "Elements" .ChoiceSequence}}
                {{template "AttributeGroups" .AttributeGroup}}
			{{end}}
		}
	{{end}}
{{end}}
`
var typesTmplComplexInline = `
{{define "SimpleType"}}
	{{$type := replaceReservedWords .Name | makePublic}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$type}} []{{toGoType .List.ItemType }}
	{{else if ne .Union.MemberTypes ""}}
		type {{$type}} string
	{{else if .Union.SimpleType}}
		type {{$type}} string
	{{else}}
		type {{$type}} {{toGoType .Restriction.Base}}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$type}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}
    
    {{template "Elements" .Extension.SequenceChoice}}
    {{template "Elements" .Extension.SequenceChoiceSequence}} 
	{{template "Elements" .Extension.Sequence}}
	{{template "Attributes" .Extension.Attributes}}
    {{template "Elements" .Extension.Choice}}
    {{template "AttributeGroups" .Extension.AttributeGroup}}
{{end}}

{{define "Attributes"}}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ .Name | makeFieldPublic}} {{toGoType .Type}} ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
		{{ else }}
			{{ .Name | makeFieldPublic}} string ` + "`" + `xml:"{{.Name}},attr,omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}
{{define "AttributeGroups"}}
	{{range .}}
		{{if ne .Ref ""}}
          {{template "Attributes" getAttributesFromGroup .Ref}}
         {{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base}} ` + "`" + `xml:",chardata"` + "`" + `
	{{template "Attributes" .Extension.Attributes}}
    {{template "AttributeGroups" .Extension.AttributeGroup}}
{{end}}

{{define "ComplexTypeInline"}}
	{{replaceReservedWords .Key | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
            {{template "Elements" .SequenceSequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
            {{template "Elements" .SequenceChoiceSequence}} 
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
            {{template "Elements" .ChoiceSequence}}
            {{template "AttributeGroups" .AttributeGroup}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Key}},omitempty"` + "`" + `
{{end}}

{{define "Elements"}}
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Ref | toGoType}} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
				 {{ $complexInlineName := setElementInComplexInlineCache .}}
                 {{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}*  {{$complexInlineName}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` 
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Type | toGoType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{ range $Key, $Value  := getComplexInlineCache }}
 type {{replaceReservedWords $Key | makePublic}} struct {
	{{with $Value.ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
            {{template "Elements" .SequenceChoiceSequence}} 
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
            {{template "Elements" .ChoiceSequence}}
            {{template "AttributeGroups" .AttributeGroup}}
		{{end}}
	{{end}}
	} 
{{end}}
`
