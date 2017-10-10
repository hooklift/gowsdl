package gowsdl

import (
	"encoding/xml"
	"strings"
)

type traverser struct {
	c   *XSDSchema
	all []*XSDSchema
}

func newTraverser(c *XSDSchema, all []*XSDSchema) *traverser {
	return &traverser{
		c:   c,
		all: all,
	}
}

func (t *traverser) traverse() {
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

func (t *traverser) traverseElements(ct []*XSDElement) {
	for _, elm := range ct {
		t.traverseElement(elm)
	}
}

func (t *traverser) traverseElement(elm *XSDElement) {
	if elm.ComplexType != nil {
		t.traverseComplexType(elm.ComplexType)
	}
	if elm.SimpleType != nil {
		t.traverseSimpleType(elm.SimpleType)
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
	t.traverseAttributes(ct.SimpleContent.Extension.Attributes)
}

func (t *traverser) traverseAttributes(attrs []*XSDAttribute) {
	for _, attr := range attrs {
		t.traverseAttribute(attr)
	}
}

func (t *traverser) traverseAttribute(attr *XSDAttribute) {
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
