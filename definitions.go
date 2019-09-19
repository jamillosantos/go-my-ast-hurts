package myasthurts

type Type struct {
	Package *Package
	Name    string
	// TODO
}

type Variable struct {
	Name string
	Type *Type
}

type Constant struct {
	Name string
	Type *Type
}

type MethodArgument struct {
	Name string
	Type string
}

type TagParam struct {
	Name  string
	Value string
}

type Tag struct {
	Raw    string
	Params []TagParam
}

type MethodResult struct {
	Name string
	Type *Type
}

type MethodDescriptor struct {
	Name      string
	Comment   string
	Recv      []MethodArgument
	Arguments []MethodArgument
	Result    []MethodResult
	Tag       Tag
}

type Interface struct {
	Name    string
	Methods []MethodDescriptor
	Comment string
}

type Field struct {
	Name    string
	Type    *Type
	Tag     Tag
	Comment string
}

type StructMethod struct {
	Descriptor *MethodDescriptor
	// TODO
}

type Struct struct {
	Name       string
	Comment    string
	Fields     []*Field
	Methods    []*StructMethod
	Interfaces []*Interface
}

type File struct {
	Package    *Package
	FileName   string
	Comment    string
	Variables  []*Variable
	Constants  []*Constant
	Structs    []*Struct
	Interfaces []*Interface
	Files      []*File
}

type Package struct {
	Name        string
	Comment     string
	Directory   string
	Variables   []*Variable
	Constants   []*Constant
	Methods     []*MethodDescriptor
	Structs     []*Struct
	Interfaces  []*Interface
	Files       []*File
	Parent      *Package
	Subpackages []*Package
}

type Environment struct {
	Packages []*Package
}
