package myasthurts

type Interface struct {
	BaseType
	Doc Doc
}

// NewInterface Create new Interface.
func NewInterface(pkg *Package, name string) *Interface {
	return &Interface{
		BaseType: *NewBaseType(pkg, name),
	}
}
