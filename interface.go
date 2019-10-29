package myasthurts

type Interface struct {
	baseType
	Doc Doc
}

// NewInterface Create new Interface.
func NewInterface(pkg *Package, name string) *Interface {
	return &Interface{
		baseType: *NewBaseType(pkg, name),
	}
}
