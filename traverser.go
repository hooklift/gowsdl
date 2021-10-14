package gowsdl

import (
	"encoding/xml"
	"strings"
)

type traverseMode int32

const (
	refResolution traverseMode = iota
	findNameByType
)

type traverser struct {
	c   *XSDSchema
	all []*XSDSchema
	tm  traverseMode
	// fields used by findNameByType mode
	typeName             string
	foundElmName         string
	conflictingTypeUsage bool
}

func newTraverser(c *XSDSchema, all []*XSDSchema) *traverser {
	return &traverser{
		c:   c,
		all: all,
		tm:  refResolution, // default traverse mode is refResolution
	}
}

func (t *traverser) traverse() {
	t.tm = refResolution

	for _, ct := range t.c.ComplexTypes {
		t.traverseComplexType(ct)
	}
	for _, st := range t.c.SimpleType {
		t.traverseSimpleType(st)
	}
	for _, elm := range t.c.Elements {
		t.traverseElement(elm)
	}
}

// Given a type, check if there is an Element with that type, and return its name.
// If multiple elements with identical names of the given type are found,
// the name is returned.
// If multiple elements with different names of the given type are found,
// the original type name is returned instead.
// If no elements are found, the original type name is returned instead.
func (t *traverser) findNameByType(name string) string {
	t.initFindNameByType(name)

	// Search for elements of given type
	for _, schema := range t.all {
		for _, elm := range schema.Elements {
			t.traverseElement(elm)
		}
		for _, ct := range schema.ComplexTypes {
			t.traverseComplexType(ct)
		}
		for _, st := range schema.SimpleType {
			t.traverseSimpleType(st)
		}
	}

	// Return found element name if given type is used only once
	if len(t.foundElmName) > 0 && !t.conflictingTypeUsage {
		return t.foundElmName
	}

	// Return original type name
	// No element found or conflicting element names found
	return t.typeName
}

func (t *traverser) initFindNameByType(name string) {
	// Initialize fields for processing
	t.tm = findNameByType
	t.typeName = stripns(name)
	t.foundElmName = ""
	t.conflictingTypeUsage = false
}

func (t *traverser) traverseElements(ct []*XSDElement) {
	for _, elm := range ct {
		t.traverseElement(elm)
	}
}

func (t *traverser) traverseElement(elm *XSDElement) {
	t.findElmName(elm)

	if elm.ComplexType != nil {
		t.traverseComplexType(elm.ComplexType)
	}
	if elm.SimpleType != nil {
		t.traverseSimpleType(elm.SimpleType)
	}
}

func (t *traverser) findElmName(elm *XSDElement) {
	// Check if we are called by findNameByType
	if t.tm != findNameByType {
		return
	}

	// Conflicting type usage already detected -> no need to search any further
	if t.conflictingTypeUsage {
		return
	}

	if stripns(elm.Type) == t.typeName {
		if len(t.foundElmName) == 0 {
			// First time usage t.typeName
			t.foundElmName = elm.Name
		} else if t.foundElmName != elm.Name {
			// Duplicate use of t.typeName with different element names
			t.conflictingTypeUsage = true
		}
	}
}

func (t *traverser) traverseSimpleType(st *XSDSimpleType) {
}

func (t *traverser) traverseComplexType(ct *XSDComplexType) {
	t.traverseElements(ct.Sequence)
	t.traverseElements(ct.Choice)
	t.traverseElements(ct.SequenceChoice)
	t.traverseElements(ct.All)
	t.traverseAttributes(ct.Attributes)
	t.traverseAttributes(ct.ComplexContent.Extension.Attributes)
	t.traverseElements(ct.ComplexContent.Extension.Sequence)
	t.traverseElements(ct.ComplexContent.Extension.Choice)
	t.traverseElements(ct.ComplexContent.Extension.SequenceChoice)
	t.traverseAttributes(ct.SimpleContent.Extension.Attributes)
}

func (t *traverser) traverseAttributes(attrs []*XSDAttribute) {
	for _, attr := range attrs {
		t.traverseAttribute(attr)
	}
}

func (t *traverser) traverseAttribute(attr *XSDAttribute) {
	// Check if we are in ref resolution mode
	if t.tm != refResolution {
		return
	}

	if attr.Ref != "" {
		refAttr := t.getGlobalAttribute(attr.Ref)
		if refAttr != nil && refAttr.Ref == "" {
			t.traverseAttribute(refAttr)
			attr.Name = refAttr.Name
			attr.Type = refAttr.Type
			if attr.Fixed == "" {
				attr.Fixed = refAttr.Fixed
			}
		}
	} else if attr.Type == "" {
		if attr.SimpleType != nil {
			t.traverseSimpleType(attr.SimpleType)
			attr.Type = attr.SimpleType.Restriction.Base
		}
	}
}

func (t *traverser) getGlobalAttribute(name string) *XSDAttribute {
	ref := t.qname(name)

	for _, schema := range t.all {
		if schema.TargetNamespace == ref.Space {
			for _, attr := range schema.Attributes {
				if attr.Name == ref.Local {
					return attr
				}
			}
		}
	}

	return nil
}

// qname resolves QName into xml.Name.
func (t *traverser) qname(name string) (qname xml.Name) {
	x := strings.SplitN(name, ":", 2)
	if len(x) == 1 {
		qname.Local = x[0]
	} else {
		qname.Local = x[1]
		qname.Space = x[0]
		if ns, ok := t.c.Xmlns[qname.Space]; ok {
			qname.Space = ns
		}
	}

	return qname
}
