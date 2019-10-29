package myasthurts

type Struct struct {
	baseType
	Doc        Doc
	Fields     []*Field
	Interfaces []*Interface
}

// NewStruct return new pointer Struct
func NewStruct(pkg *Package, name string) *Struct {
	srct := &Struct{
		baseType: baseType{
			pkg:  pkg,
			name: name,
		},
	}
	return srct
}

// Package return pointer package of Struct
func (s *Struct) Package() *Package {
	return s.pkg
}

// Name return name of Struct
func (s *Struct) Name() string {
	return s.name
}
