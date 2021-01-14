package myasthurts

type Struct struct {
	BaseType
	Doc    Doc
	Fields []*Field
}

// NewStruct return new pointer Struct
func NewStruct(pkg *Package, name string) *Struct {
	srct := &Struct{
		BaseType: *NewBaseType(pkg, name),
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

// Implements checks if this struct implements a given interface.
//
// This method uses the `MethodDescriptor.Compatible` to check if all interface
// methods have are implemented on the struct.
//
// TODO(jota): Take into consideration Interface composing...
func (s *Struct) Implements(i *Interface) bool {
	for _, m := range i.Methods() {
		method, ok := s.methodsMap[m.Descriptor.Name()]
		if !ok {
			return false
		}
		if !method.Descriptor.Compatible(m.Descriptor) {
			return false
		}
	}
	return true
}
