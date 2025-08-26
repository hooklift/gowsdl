package gowsdl

import (
	"errors"
	"fmt"
	"strings"
)

// resolveAttrRefs resolves all attribute references across the provided schemas.
// It modifies the schemas in-place, copying properties from the referenced attributes
// onto the referencing attributes.
func resolveAttrRefs(schemas []*XSDSchema) error {
	// First, build an index from (namespace, attrName) -> attrDef
	// for all global attrs across all schemas.
	attrIndex := make(map[namespacedKey]*XSDAttribute)
	for _, s := range schemas {
		for _, attr := range s.Attributes {
			if attr.Name != "" {
				attrIndex[newNamespacedKey(s.TargetNamespace, attr.Name)] = attr
			}
		}
	}

	keyFromAttrRef := func(s *XSDSchema, attrRef string) (namespacedKey, error) {
		before, after, hadColon := strings.Cut(attrRef, ":")
		if !hadColon {
			return newNamespacedKey(s.TargetNamespace, before), nil
		}
		if ns, ok := s.Xmlns[before]; ok {
			return newNamespacedKey(ns, after), nil
		}
		return "", fmt.Errorf("unable to resolve namespace prefix %q in attribute ref %q", before, attrRef)
	}

	// Next, traverse all attrs with refs and copy over the properties from the referenced attrs.
	var currentSchema *XSDSchema
	var errs []error
	visitor{schemas}.visit(&visitorConfig{
		onEnterSchema: func(s *XSDSchema) {
			currentSchema = s
		},
		onEnterAttribute: func(attr *XSDAttribute) {
			if attr.Ref != "" {
				nsk, err := keyFromAttrRef(currentSchema, attr.Ref)
				if err != nil {
					errs = append(errs, err)
					return
				}
				refAttr, ok := attrIndex[nsk]
				if !ok || refAttr.Ref != "" {
					errs = append(errs, fmt.Errorf("unable to resolve attribute ref %q in schema with namespace %q", attr.Ref, currentSchema.TargetNamespace))
					return
				}
				attr.Name = refAttr.Name
				attr.Type = refAttr.Type
				if attr.Fixed == "" {
					attr.Fixed = refAttr.Fixed
				}
				attr.TargetNamespace = currentSchema.XMLNameForAttribute(refAttr).Space
			} else if attr.Type == "" && attr.SimpleType != nil {
				attr.Type = attr.SimpleType.Restriction.Base
			}
		},
	})

	return errors.Join(errs...)
}
