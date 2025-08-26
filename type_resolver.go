package gowsdl

import (
	"encoding/xml"
	"strings"
)

type typeResolver struct {
	elementsByTypeName map[namespacedKey][]elementAndSchema
}

type namespacedKey string

func newNamespacedKey(namespace, local string) namespacedKey {
	return namespacedKey(namespace + "|" + local)
}

type elementAndSchema struct {
	element *XSDElement
	schema  *XSDSchema
}

func newTypeResolver(schemas []*XSDSchema) *typeResolver {
	var currentSchema *XSDSchema
	elementsByTypeName := make(map[namespacedKey][]elementAndSchema)

	visitor{schemas}.visit(&visitorConfig{
		onEnterSchema: func(s *XSDSchema) {
			currentSchema = s
		},
		onEnterElement: func(e *XSDElement) {
			if e.Type == "" {
				return
			}
			before, after, hadColon := strings.Cut(e.Type, ":")
			var key namespacedKey
			if hadColon {
				key = newNamespacedKey(currentSchema.Xmlns[before], after)
			} else {
				key = newNamespacedKey(
					currentSchema.XMLNameForElement(e).Space,
					e.Type,
				)
			}
			elementsByTypeName[key] = append(elementsByTypeName[key], elementAndSchema{
				element: e,
				schema:  currentSchema,
			})
		},
	})
	return &typeResolver{
		elementsByTypeName: elementsByTypeName,
	}
}

// Given a type qname and a schema that it appears within, determine the XMLName
// tag that should be used on the Go type representing that XML type.
//
// If no elements with the given type exist return the original typeName plus the schema's
// target namespace.
//
// If all elements with the given type have the same name and namespace, return that name
// and namespace.
//
// Otherwise, if there are mismatched names or namespaces among the elements with the given
// type, return the original typeName plus the schema's target namespace.
//
// Note that this last case will likely result in broken generated code.
// The more correct solution here is probably to not generate XMLName fields at all,
// except for types used directly as request/response wrappers, and instead emit namespaces
// on field tags, but that would be a somewhat involved change.
func (tr *typeResolver) xmlNameForType(typeName string, schema *XSDSchema) xml.Name {
	key := newNamespacedKey(schema.TargetNamespace, typeName)
	elements, ok := tr.elementsByTypeName[key]
	if !ok {
		return xml.Name{
			Space: schema.TargetNamespace,
			Local: typeName,
		}
	}

	elementNames := make(map[string]struct{})
	namespaces := make(map[string]struct{})
	for _, e := range elements {
		elementNames[e.element.Name] = struct{}{}
		namespaces[e.schema.XMLNameForElement(e.element).Space] = struct{}{}
	}
	if len(elementNames) == 1 && len(namespaces) == 1 {
		e := elements[0]
		return e.schema.XMLNameForElement(e.element)
	}

	return xml.Name{
		Space: schema.TargetNamespace,
		Local: typeName,
	}
}
