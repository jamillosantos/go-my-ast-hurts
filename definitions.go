package myasthurts

import (
	"fmt"
	"go/build"
	"regexp"
	"strings"
)

type Constant struct {
	Name string
	Type Type
}

type Doc struct {
	Comments []string
}

// EnvConfig is a Struct to set config in Environment
type EnvConfig struct {
	DevMode    bool
	ASTI       bool
	CurrentDir string
}

func (ec EnvConfig) CWD() string {
	return ec.CurrentDir
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
	Constants  []*Constant
	Structs    []*Struct
	Interfaces []*Interface
	Files      []*File
}

type Package struct {
	Name        string
	ImportPath  string
	RealPath    string
	BuildInfo   *build.Package
	Doc         Doc
	explored    bool
	Variables   []*Variable
	Constants   []*Constant
	Methods     []*MethodDescriptor
	methodsMap  map[string]*MethodDescriptor
	Structs     []*Struct
	Interfaces  []*Interface
	RefType     []RefType
	refTypeMap  map[string]RefType
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
		BuildInfo:  buildPackage,
		Constants:  make([]*Constant, 0),
		Methods:    make([]*MethodDescriptor, 0),
		Variables:  make([]*Variable, 0),
		methodsMap: make(map[string]*MethodDescriptor),
		Structs:    make([]*Struct, 0),
		Interfaces: make([]*Interface, 0),
		RefType: []RefType{
			NullRefType,
			InterfaceRefType,
		},
		refTypeMap:  make(map[string]RefType),
		Types:       make([]Type, 0),
		Files:       make([]*File, 0),
		Subpackages: make([]*Package, 0),
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

var (
	NullRefType = &BaseRefType{
		name: "nil",
	}
	InterfaceRefType = &BaseRefType{
		name: "interface",
	}
)

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

type StarRefType struct {
	RefType
}

func NewStarRefType(refType RefType) *StarRefType {
	return &StarRefType{
		RefType: refType,
	}
}

func (refType *StarRefType) Name() string {
	return refType.RefType.Name()
}

func (refType *StarRefType) Pkg() *Package {
	return refType.RefType.Pkg()
}

func (refType *StarRefType) Type() Type {
	return refType.RefType.Type()
}

func (refType *StarRefType) AppendType(tp Type) {
	refType.RefType.AppendType(tp)
}

type ArrayRefType struct {
	RefType
}

func NewArrayRefType(refType RefType) *ArrayRefType {
	return &ArrayRefType{
		RefType: refType,
	}
}

func (refType *ArrayRefType) Name() string {
	return refType.RefType.Name()
}

func (refType *ArrayRefType) Pkg() *Package {
	return refType.RefType.Pkg()
}

func (refType *ArrayRefType) Type() Type {
	return refType.RefType.Type()
}

func (refType *ArrayRefType) AppendType(tp Type) {
	refType.RefType.AppendType(tp)
}

type ChanRefType struct {
	RefType
}

func NewChanRefType(refType RefType) *ChanRefType {
	return &ChanRefType{
		RefType: refType,
	}
}

func (refType *ChanRefType) Name() string {
	return refType.RefType.Name()
}

func (refType *ChanRefType) Pkg() *Package {
	return refType.RefType.Pkg()
}

func (refType *ChanRefType) Type() Type {
	return refType.RefType.Type()
}

func (refType *ChanRefType) AppendType(tp Type) {
	refType.RefType.AppendType(tp)
}

type EllipsisRefType struct {
	RefType
}

func NewEllipsisRefType(refType RefType) *EllipsisRefType {
	return &EllipsisRefType{
		RefType: refType,
	}
}

func (refType *EllipsisRefType) Name() string {
	return refType.RefType.Name()
}

func (refType *EllipsisRefType) Pkg() *Package {
	return refType.RefType.Pkg()
}

func (refType *EllipsisRefType) Type() Type {
	return refType.RefType.Type()
}

func (refType *EllipsisRefType) AppendType(tp Type) {
	refType.RefType.AppendType(tp)
}

type MapType struct {
	pkg   *Package
	Key   RefType
	Value RefType
}

func NewMap(pkg *Package, key RefType, value RefType) *MapType {
	return &MapType{
		Key:   key,
		Value: value,
	}
}

func (mt *MapType) Package() *Package {
	return mt.pkg
}

func (mt *MapType) Name() string {
	return fmt.Sprintf("map[%s]%s", mt.Key.Name(), mt.Value.Name())
}

func (mt *MapType) Methods() []*TypeMethod {
	return nil
}

func (mt *MapType) AddMethod(method *TypeMethod) {}

type Ellipsis struct {
	pkg *Package
	Elt RefType
}

func NewEllipsis(pkg *Package, elt RefType) *Ellipsis {
	return &Ellipsis{
		pkg: pkg,
		Elt: elt,
	}
}

func (el *Ellipsis) Package() *Package {
	return el.pkg
}

func (el *Ellipsis) Name() string {
	return ""
}

func (el *Ellipsis) Methods() []*TypeMethod {
	return nil
}

func (el *Ellipsis) AddMethod(method *TypeMethod) {}

type Type interface {
	Package() *Package
	Name() string
	Methods() []*TypeMethod
	AddMethod(*TypeMethod)
}

type baseType struct {
	pkg        *Package
	name       string
	methods    []*TypeMethod
	methodsMap map[string]*TypeMethod
}

// NewBaseType creates a new initialized baseType.
func NewBaseType(pkg *Package, name string) *baseType {
	return &baseType{
		pkg:        pkg,
		name:       name,
		methods:    make([]*TypeMethod, 0),
		methodsMap: make(map[string]*TypeMethod, 0),
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
	t.methodsMap[method.Descriptor.Name()] = method
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
	Doc     Doc
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

// AppendStruct add new Struct in Package
func (p *Package) AppendStruct(s *Struct) {
	p.Structs = append(p.Structs, s)
	p.Types = append(p.Types, s)
}

// StructByName find Struct by name.
func (p *Package) StructByName(name string) (*Struct, bool) {
	for _, e := range p.Structs {
		if e.Name() == name {
			return e, true
		}
	}
	return nil, false
}

// AppendInterface registers a new Interface to the package.
func (p *Package) AppendInterface(s *Interface) {
	p.Interfaces = append(p.Interfaces, s)
	p.Types = append(p.Types, s)
}

// StructByName find Struct by name.
func (p *Package) InterfaceByName(name string) (*Interface, bool) {
	for _, e := range p.Interfaces {
		if e.Name() == name {
			return e, true
		}
	}
	return nil, false
}

// EnsureRefType will try to get the RefType from the list by the name param. If
// the RefType exists, it will return the reference with true (second return).
// If it is does not exists in the list, the function will create a new type and
// return it with false (second return).
func (p *Package) EnsureRefType(name string) (RefType, bool) {
	refType, ok := p.RefTypeByName(name)
	if ok {
		return refType, true
	}
	refType = NewRefType(name, p, nil)
	p.AddRefType(refType)
	return refType, false
}

// RefTypeByName find RefType by name.
func (p *Package) RefTypeByName(name string) (RefType, bool) {
	refType, ok := p.refTypeMap[name]
	return refType, ok
}

func (p *Package) AddRefType(ref RefType) {
	if ref.Name() != "" {
		p.RefType = append(p.RefType, ref)
		p.refTypeMap[ref.Name()] = ref
	}
}

// AppendRefType add new RefType in Package.
func (p *Package) AppendRefType(name string) (ref RefType) {
	ref = &BaseRefType{
		pkg:  p,
		name: name,
	}
	p.AddRefType(ref)
	return ref
}

func (p *Package) AppendMethod(method *MethodDescriptor) {
	p.Methods = append(p.Methods, method)
	p.methodsMap[method.Name()] = method
}

func (p *Package) MethodByName(name string) (*MethodDescriptor, bool) {
	m, ok := p.methodsMap[name]
	return m, ok
}

func (p *Package) VariableByName(name string) (vrle *Variable) {
	for _, v := range p.Variables {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (p *Package) AppendVariable(variable *Variable) *Variable {
	p.Variables = append(p.Variables, variable)
	return variable
}

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
