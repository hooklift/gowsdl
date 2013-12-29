package main

var typesTmpl = `
package main
import (
	"encoding/xml"
	"github.com/c4milo/gowsdl"
)
{{range .Schemas}}
	{{range .SimpleType}}
		{{$name := .Name}}
		type {{.Name}} string
		const (
			{{with .Restriction}}
				{{range .Enumeration}}
				{{.Value}} {{$name}} = "{{.Value}}" {{end}}
			{{end}}
		)
	{{end}}
	{{range .Elements}}
	{{end}}
	{{range .ComplexTypes}}
	{{end}}
{{end}}`
