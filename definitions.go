package myasthurts

type Type interface {
	Package() *Package
	Name() string
}

type MethodDescriptor struct {
	pkg       *Package
	name      string
	Comment   []string
	Recv      []MethodArgument
	Arguments []MethodArgument
	Result    []MethodResult
	Tag       Tag
}

func NewMethodDescriptor(pkg *Package, name string) *MethodDescriptor {
	return &MethodDescriptor{
		pkg:  pkg,
		name: name,
	}
}

func (method *MethodDescriptor) Package() *Package {
	return method.pkg
}

func (method *MethodDescriptor) Name() string {
	return method.name
}

type Interface struct {
	pkg     *Package
	name    string
	Methods []MethodDescriptor
	Comment []string
}

func NewInterface(pkg *Package, name string) *Interface {
	return &Interface{
		pkg:  pkg,
		name: name,
	}
}

func (i *Interface) Package() *Package {
	return i.pkg
}

func (i *Interface) Name() string {
	return i.name
}

type Struct struct {
	pkg        *Package
	name       string
	Comment    []string
	Fields     []*Field
	Methods    []*StructMethod
	Interfaces []*Interface
}

func (s *Struct) Package() *Package {
	return s.pkg
}

func (s *Struct) Name() string {
	return s.name
}

// NewStruct sdsdsds
func NewStruct(pkg *Package, name string) *Struct {
	srct := &Struct{
		pkg:  pkg,
		name: name,
	}
	pkg.AppendStruct(srct)
	return srct
}

type Variable struct {
	Name string
	Type *RefType
}

type Constant struct {
	Name string
	Type Type
}

type MethodArgument struct {
	Name string
	Type *RefType
}

type RefType struct {
	Name string
	Pkg  *Package
	Type Type
}

func NewRefType(pkg *Package) *RefType {
	return &RefType{
		Pkg: pkg,
	}
}

func (rt *RefType) AppendType(tp Type) {
	if rt.Type == nil {
		switch t := tp.(type) {
		case *Struct:
			var mDescriptor *StructMethod
			for _, s := range t.Package().Methods {
				mDescriptor = &StructMethod{
					Descriptor: s,
				}
				t.Methods = append(t.Methods, mDescriptor)
			}
			rt.Type = tp
		}
	} else {
		rt.Type = tp
	}
}

type Tag struct {
	Raw    string
	Params []TagParam
}

func (t *Tag) AppendTagParam(tNew *TagParam) bool {
	_, ok := t.TagParamByName(tNew.Name)
	if ok {
		return !ok
	}
	t.Params = append(t.Params, *tNew)
	return ok
}

func (t *Tag) TagParamByName(name string) (*TagParam, bool) {
	for _, tp := range t.Params {
		if tp.Name == name {
			return &tp, true
		}
	}
	return nil, false
}

type TagParam struct {
	Name    string
	Value   string
	Options []string
}

type MethodResult struct {
	Name string
	Type Type
}

type Field struct {
	Name    string
	Type    *RefType
	Tag     Tag
	Comment []string
}

type StructMethod struct {
	Descriptor *MethodDescriptor
	// TODO
}

type File struct {
	Package    *Package
	FileName   string
	Comment    []string
	Variables  []*Variable
	Constants  []*Constant
	Structs    []*Struct
	Interfaces []*Interface
	Files      []*File
}

type Package struct {
	Name        string
	Comment     []string
	Directory   string
	Variables   []*Variable
	Constants   []*Constant
	Methods     []*MethodDescriptor
	Structs     []*Struct
	Interfaces  []*Interface
	RefType     []*RefType
	Types       []Type
	Files       []*File
	Parent      *Package
	Subpackages []*Package
}

func (p *Package) AppendStruct(s *Struct) {
	p.Structs = append(p.Structs, s)
}

func (p *Package) RefTypeByName(name string) (*RefType, bool) {
	for _, pr := range p.RefType {
		if name == pr.Name {
			return pr, true
		}
	}
	return nil, false
}

type Environment struct {
	packages    []*Package
	packagesMap map[string]*Package
}

func NewEnvironment() *Environment {
	return &Environment{
		packages:    []*Package{},
		packagesMap: map[string]*Package{},
	}
}

func (e *Environment) PackageByName(name string) (*Package, bool) {
	pkg, ok := e.packagesMap[name]
	return pkg, ok
}

func (e *Environment) AppendPackage(pkg *Package) {
	e.packages = append(e.packages, pkg)
	e.packagesMap[pkg.Name] = pkg
}
