package myasthurts

type Interface struct {
	pkg     *Package
	name    string
	Methods []MethodDescriptor
	Doc     Doc
}

// NewInterface Create new Interface.
func NewInterface(pkg *Package, name string) *Interface {
	return &Interface{
		pkg:  pkg,
		name: name,
	}
}

// Package get name of Pacakge.
func (i *Interface) Package() *Package {
	return i.pkg
}

// Name return name of Struct than implement this Interface.
func (i *Interface) Name() string {
	return i.name
}
