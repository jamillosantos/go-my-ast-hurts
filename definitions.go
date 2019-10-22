package myasthurts

import (
	"go/build"
	"regexp"
	"strings"
)

type constant struct {
	Name string
	Type Type
}

type Doc struct {
	Comments []string
}

// EnvConfig is a Struct to set config in Environment
type EnvConfig struct {
	DevMode bool
	ASTI    bool
}

// Field is utilized in Struct type in the present moment.
type Field struct {
	Name    string
	RefType RefType
	Tag     Tag
	Doc     Doc
}

// File is utilized to represent each file read in Package.
type File struct {
	Package    *Package
	FileName   string
	Doc        Doc
	Variables  []*Variable
	Constants  []*constant
	Structs    []*Struct
	Interfaces []*Interface
	Files      []*File
}

type Interface struct {
	pkg     *Package
	name    string
	Methods []MethodDescriptor
	Doc     Doc
}

// MethodArgument represent type of fields and arguments.
type MethodArgument struct {
	Name string
	Type RefType
}

type MethodDescriptor struct {
	pkg       *Package
	name      string
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

type Package struct {
	Name        string
	ImportPath  string
	RealPath    string
	explored    bool
	Doc         Doc
	Variables   []*Variable
	Constants   []*constant
	Methods     []*MethodDescriptor
	Structs     []*Struct
	Interfaces  []*Interface
	RefType     []RefType
	Types       []Type
	Files       []*File
	Parent      *Package
	Subpackages []*Package
}

func NewPackage(buildPackage *build.Package) *Package {
	return &Package{
		Name:       buildPackage.Name,
		ImportPath: buildPackage.ImportPath,
		RealPath:   buildPackage.Dir,
	}
}

type RefType interface {
	Name() string
	Pkg() *Package
	Type() Type
	AppendType(t Type)
}

type BaseRefType struct {
	name string
	pkg  *Package
	t    Type
}

func NewRefType(name string, pkg *Package, t Type) RefType {
	return &BaseRefType{
		name: name,
		pkg:  pkg,
		t:    t,
	}
}

func (refType *BaseRefType) Name() string {
	return refType.name
}

func (refType *BaseRefType) Pkg() *Package {
	return refType.pkg
}

func (refType *BaseRefType) Type() Type {
	return refType.t
}

// AppendType add Type in RefType
func (rt *BaseRefType) AppendType(tp Type) {
	rt.t = tp
}

type Type interface {
	Package() *Package
	Name() string
	Methods() []*TypeMethod
	AddMethod(*TypeMethod)
}

type baseType struct {
	pkg     *Package
	name    string
	methods []*TypeMethod
}

// NewBaseType creates a new initialized baseType.
func NewBaseType(pkg *Package, name string) *baseType {
	return &baseType{
		pkg:     pkg,
		name:    name,
		methods: make([]*TypeMethod, 0),
	}
}

func (t *baseType) Package() *Package {
	return t.pkg
}

func (t *baseType) Name() string {
	return t.name
}

func (t *baseType) Methods() []*TypeMethod {
	return t.methods
}

func (t *baseType) AddMethod(method *TypeMethod) {
	t.methods = append(t.methods, method)
	// TODO(jota): Add map to improve future use.
}

type Struct struct {
	baseType
	Doc        Doc
	Fields     []*Field
	Interfaces []*Interface
}

type TypeMethod struct {
	Name       string
	Descriptor *MethodDescriptor
}

type Tag struct {
	Raw    string
	Params []TagParam
}

type TagParam struct {
	Name    string
	Value   string
	Options []string
}

type Variable struct {
	Name    string
	RefType RefType
}

// FormatComment is simple method to remove // or /* */ of comment
func (doc *Doc) FormatComment() string {
	str := ""
	reg := regexp.MustCompile(`(\/{2}|\/\*|\*\/)`)
	for _, c := range doc.Comments {
		str += strings.TrimSpace(reg.ReplaceAllString(c, ""))
		strArr := reg.FindStringSubmatch(c)
		if len(strArr) > 0 && strArr[0] == "//" {
			str += "\n"
		}
	}
	return str
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

// NewMethodDescriptor return the pointer of new MethodDescriptor
func NewMethodDescriptor(pkg *Package, name string) *MethodDescriptor {
	return &MethodDescriptor{
		pkg:  pkg,
		name: name,
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

// AppendStruct add new Struct in Package
func (p *Package) AppendStruct(s *Struct) {
	p.Structs = append(p.Structs, s)
}

// StructByName find Struct by name.
func (p *Package) StructByName(name string) *Struct {
	for _, e := range p.Structs {
		if e.Name() == name {
			return e
		}
	}
	return nil
}

// RefTypeByName find RefType by name.
func (p *Package) RefTypeByName(name string) (RefType, bool) {
	for _, pp := range p.RefType {
		if name == pp.Name() {
			return pp, true
		}
	}
	return nil, false
}

func (p *Package) AddRefType(ref RefType) {
	p.RefType = append(p.RefType, ref)
}

// AppendRefType add new RefType in Package.
func (p *Package) AppendRefType(name string) (ref RefType) {
	ref = &BaseRefType{
		pkg:  p,
		name: name,
	}
	p.RefType = append(p.RefType, ref)
	return ref
}

func (p *Package) VariableByName(name string) (vrle *Variable) {
	for _, v := range p.Variables {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (p *Package) AppendVariable(vrle *Variable) (v *Variable) {
	p.Variables = append(p.Variables, vrle)
	return vrle
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

// FormatComments show struct with all comments
/*func (s *Struct) FormatComments() string {
	str := fmt.Sprintf("%s\ntype %s struct {\n", s.Doc.FormatComment(), s.Name())
	for _, f := range s.Fields {
		c := f.Doc.FormatComment()
		brk := "\n"
		if c == "" {
			brk = ""
		}
		str += fmt.Sprintf("%s%s%s %s %s\n", c, brk, f.Name, f.RefType.Name, f.Tag.Raw)
	}
	return fmt.Sprintf("%s}", str)
}*/

// AppendTagParam add new TagParam in Tag
func (t *Tag) AppendTagParam(tNew *TagParam) bool {
	tp := t.TagParamByName(tNew.Name)
	if tp != nil {
		return false
	}
	t.Params = append(t.Params, *tNew)
	return false
}

// TagParamByName find TagParam by name.
func (t *Tag) TagParamByName(name string) *TagParam {
	for _, tp := range t.Params {
		if tp.Name == name {
			return &tp
		}
	}
	return nil
}
