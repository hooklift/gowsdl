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

{{define "ComplexContent"}}
	{{$baseType := toGoType .Extension.Base}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence.Elements}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
	{{range .}}
		{{ .Name | makePublic}} {{toGoType .Type}}{{end}}
{{end}}

{{define "SimpleContent"}}
	Value {{toGoType .Extension.Base}}{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeGlobal"}}
	{{$name := replaceReservedWords .Name}}
	type {{$name}} struct {
		XMLName xml.Name ` + "`xml:\"{{getTargetNamespace}} {{$name}}\"`" + `
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{ else }}
			{{template "Elements" .Sequence.Elements}}
			{{template "Elements" .Choice}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
		{{end}}
	}
{{end}}

{{define "ComplexTypeLocal"}}
	{{$name := replaceReservedWords .Name}}
	{{with .ComplexType}}
		type {{$name}} struct {
			XMLName xml.Name ` + "`xml:\"{{getTargetNamespace}} {{$name}}\"`" + `
			{{if ne .ComplexContent.Extension.Base ""}}
				{{template "ComplexContent" .ComplexContent}}
			{{else if ne .SimpleContent.Extension.Base ""}}
				{{template "SimpleContent" .SimpleContent}}
			{{ else }}
				{{template "Elements" .Sequence.Elements}}
				{{template "Elements" .Choice}}
				{{template "Elements" .All}}
				{{template "Attributes" .Attributes}}
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
