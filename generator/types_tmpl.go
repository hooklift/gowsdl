package generator

var typesTmpl = `
{{define "SimpleType"}}
	{{$type := replaceReservedWords .Name}}
	type {{$type}} {{toGoType .Restriction.Base}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{$type}}_{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$type}} = "{{$value}}" {{end}}
		{{end}}
	)
{{end}}

{{define "ComplexTypeGlobal"}}
	{{$name := replaceReservedWords .Name}}
	type {{$name}} struct {
		{{if ne .ComplexContent.Extension.Base ""}}
			{{$baseType := .ComplexContent.Extension.Base}}
			{{ if $baseType }}
				*{{stripns $baseType}}
			{{end}}

			{{template "Elements" .ComplexContent.Extension.Sequence.Elements}}
		{{ else }}
			{{template "Elements" .Sequence.Elements}}
		{{end}}
	}
{{end}}

{{define "ComplexTypeLocal"}}
	{{$name := replaceReservedWords .Name}}

	{{with .ComplexType}}
		type {{$name}} struct {
			{{if ne .ComplexContent.Extension.Base ""}}
				{{$baseType := .ComplexContent.Extension.Base}}
				{{ if $baseType }}
					*{{stripns $baseType}}
				{{end}}

				{{template "Elements" .ComplexContent.Extension.Sequence.Elements}}
			{{ else }}
				{{template "Elements" .Sequence.Elements}}
			{{end}}
		}
	{{end}}
{{end}}

{{define "Elements"}}
	{{range .}}
		{{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Type | toGoType}} ` + "`" + `xml:"{{.Name}},omitempty"` + "`" + `{{end}}
{{end}}

{{range .Schemas}}
	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}
	{{range .Elements}}
		{{if not .Type}}
			{{template "ComplexTypeLocal" .}}
		{{end}}
	{{end}}
	{{range .ComplexTypes}}
		{{template "ComplexTypeGlobal" .}}
	{{end}}
{{end}}
`
