package gowsdl

type visitor struct {
	all []*XSDSchema
}

type visitorConfig struct {
	onEnterSchema      func(*XSDSchema)
	onExitSchema       func(*XSDSchema)
	onEnterElement     func(*XSDElement)
	onExitElement      func(*XSDElement)
	onEnterComplexType func(*XSDComplexType)
	onExitComplexType  func(*XSDComplexType)
	onEnterSimpleType  func(*XSDSimpleType)
	onExitSimpleType   func(*XSDSimpleType)
	onEnterAttribute   func(*XSDAttribute)
	onExitAttribute    func(*XSDAttribute)
}

func (v visitor) visit(cfg *visitorConfig) {
	for _, schema := range v.all {
		if cfg.onEnterSchema != nil {
			cfg.onEnterSchema(schema)
		}
		v.visitSchema(schema, cfg)
		if cfg.onExitSchema != nil {
			cfg.onExitSchema(schema)
		}
	}
}

func (v visitor) visitSchema(s *XSDSchema, cfg *visitorConfig) {
	for _, elm := range s.Elements {
		v.visitElement(elm, cfg)
	}
	for _, ct := range s.ComplexTypes {
		v.visitComplexType(ct, cfg)
	}
	for _, st := range s.SimpleType {
		v.visitSimpleType(st, cfg)
	}
}

func (v visitor) visitElements(es []*XSDElement, cfg *visitorConfig) {
	for _, e := range es {
		v.visitElement(e, cfg)
	}
}

func (v visitor) visitAttribute(attr *XSDAttribute, cfg *visitorConfig) {
	if cfg.onEnterAttribute != nil {
		cfg.onEnterAttribute(attr)
	}
	if cfg.onExitAttribute != nil {
		cfg.onExitAttribute(attr)
	}
}

func (v visitor) visitAttributes(attrs []*XSDAttribute, cfg *visitorConfig) {
	for _, attr := range attrs {
		v.visitAttribute(attr, cfg)
	}
}

func (v visitor) visitElement(e *XSDElement, cfg *visitorConfig) {
	if cfg.onEnterElement != nil {
		cfg.onEnterElement(e)
	}

	if e.ComplexType != nil {
		v.visitComplexType(e.ComplexType, cfg)
	}
	if e.SimpleType != nil {
		v.visitSimpleType(e.SimpleType, cfg)
	}

	if cfg.onExitElement != nil {
		cfg.onExitElement(e)
	}
}

func (v visitor) visitComplexType(ct *XSDComplexType, cfg *visitorConfig) {
	if cfg.onEnterComplexType != nil {
		cfg.onEnterComplexType(ct)
	}

	v.visitElements(ct.Sequence, cfg)
	v.visitElements(ct.Choice, cfg)
	v.visitElements(ct.SequenceChoice, cfg)
	v.visitElements(ct.All, cfg)
	v.visitAttributes(ct.Attributes, cfg)
	v.visitAttributes(ct.ComplexContent.Extension.Attributes, cfg)
	v.visitElements(ct.ComplexContent.Extension.Sequence, cfg)
	v.visitElements(ct.ComplexContent.Extension.Choice, cfg)
	v.visitElements(ct.ComplexContent.Extension.SequenceChoice, cfg)
	v.visitAttributes(ct.SimpleContent.Extension.Attributes, cfg)

	if cfg.onExitComplexType != nil {
		cfg.onExitComplexType(ct)
	}
}

func (v visitor) visitSimpleType(st *XSDSimpleType, cfg *visitorConfig) {
	if cfg.onEnterSimpleType != nil {
		cfg.onEnterSimpleType(st)
	}
	if cfg.onExitSimpleType != nil {
		cfg.onExitSimpleType(st)
	}
}
