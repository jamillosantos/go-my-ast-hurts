package myasthurts

// MethodArgument represent type of fields and arguments.
type MethodArgument struct {
	Name string
	Type RefType
	Doc  Doc
}

type MethodDescriptor struct {
	baseType
	Doc       Doc
	Recv      []MethodArgument
	Arguments []MethodArgument
	Result    []MethodResult
	Tag       Tag
}

type MethodResult struct {
	Name string
	Type Type
}

// NewMethodDescriptor return the pointer of new MethodDescriptor
func NewMethodDescriptor(pkg *Package, name string) *MethodDescriptor {
	return &MethodDescriptor{
		baseType: *NewBaseType(pkg, name),
	}
}

// Package return pointer of Package
func (method *MethodDescriptor) Package() *Package {
	return method.pkg
}

// Name return name of Method
func (method *MethodDescriptor) Name() string {
	return method.name
}
