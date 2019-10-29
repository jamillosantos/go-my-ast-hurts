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
	Type RefType
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

// Compatible checks if the signature of both method descriptor are compatible.
//
// It checks if the all arguments refer to the same RefType. The same happens
// with the result.
//
// Receivers are not taken into consideration, neither names.
func (method *MethodDescriptor) Compatible(m *MethodDescriptor) bool {
	if len(method.Arguments) != len(m.Arguments) {
		return false
	}
	if len(method.Result) != len(m.Result) {
		return false
	}
	for i, arg := range method.Arguments {
		if m.Arguments[i].Type != arg.Type {
			return false
		}
	}
	for i, r := range method.Result {
		if m.Result[i].Type != r.Type {
			return false
		}
	}
	return true
}
