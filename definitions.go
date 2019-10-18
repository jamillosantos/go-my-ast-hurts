package myasthurts

import (
	"errors"
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
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

type environment struct {
	packages    []*Package
	packagesMap map[string]*Package
	Config      EnvConfig
}

// Field is utilized in Struct type in the present moment.
type Field struct {
	Name    string
	RefType *RefType
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
	Type *RefType
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
	Directory   string
	Doc         Doc
	Variables   []*Variable
	Constants   []*constant
	Methods     []*MethodDescriptor
	Structs     []*Struct
	Interfaces  []*Interface
	RefType     []*RefType
	Types       []Type
	Files       []*File
	Parent      *Package
	Subpackages []*Package
}

type parseContext struct {
	File        *ast.File
	Env         *environment
	PackagesMap map[string]*Package
}

type RefType struct {
	Name string
	Pkg  *Package
	Type Type
}

type Struct struct {
	pkg        *Package
	name       string
	Doc        Doc
	Fields     []*Field
	Methods    []*StructMethod
	Interfaces []*Interface
}

type StructMethod struct {
	Descriptor *MethodDescriptor
	// TODO
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

type Type interface {
	Package() *Package
	Name() string
}

type Variable struct {
	Name string
	Type *RefType
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

// NewEnvironment is the method allow start parse in file
func NewEnvironment() (env *environment, exrr error) {

	env = &environment{
		packages:    []*Package{},
		packagesMap: map[string]*Package{},
	}

	if exrr = env.makeEnv(); exrr != nil {
		return nil, exrr
	}
	return env, nil
}

// PackageByName find Package by name in Environment.
func (env *environment) PackageByName(name string) (*Package, bool) {
	pkg, ok := env.packagesMap[name]
	return pkg, ok
}

// AppendPackage add new Pacakge in Environment.
func (env *environment) AppendPackage(pkg *Package) {
	env.packages = append(env.packages, pkg)
	env.packagesMap[pkg.Name] = pkg
}

// ParsePackage checks if the Path ou File exist before of parse.
func (env *environment) ParsePackage(pathOrName string, isFile bool) (exrr error) {

	if isFile {
		if _, ok := os.Stat(pathOrName); os.IsNotExist(ok) {
			return errors.New("File not found")
		}
		env.parse(pathOrName)
	} else {
		files, err := ioutil.ReadDir(pathOrName)
		if err != nil {
			return err
		}

		for _, file := range files {
			fileName := file.Name()
			if s, _ := regexp.MatchString(`(?ms)test\b`, fileName); s {
				continue
			}
			fileLocation := fmt.Sprintf("%s/%s", pathOrName, fileName)
			env.parse(fileLocation)
		}
	}

	return
}

func (env *environment) basePath() (rtn string, exrr error) {

	if rtn = os.Getenv("GOROOT"); rtn == "" {
		return "", errors.New("GOROOT environment variable not found or is empty")
	}
	return fmt.Sprintf("%s/src", rtn), nil
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
func (p *Package) RefTypeByName(name string) (ref *RefType) {
	for _, pp := range p.RefType {
		if name == pp.Name {
			return pp
		}
	}
	return nil
}

// AppendRefType add new RefType in Package.
func (p *Package) AppendRefType(name string) (ref *RefType, exrr error) {
	ref = &RefType{
		Pkg:  p,
		Name: name,
	}
	p.RefType = append(p.RefType, ref)
	return ref, nil
}

// NewRefType return new pointer RefType
func NewRefType(pkg *Package) *RefType {
	return &RefType{
		Pkg: pkg,
	}
}

// AppendType add Type in RefType
func (rt *RefType) AppendType(tp Type) {

	if rt.Type == nil {
		switch t := tp.(type) {
		case *Struct:
			rt.Type = tp
			var mDescriptor *StructMethod
			for _, s := range t.Package().Methods {
				if len(s.Recv) > 0 && s.Recv[0].Type != nil && s.Recv[0].Type.Name == t.Name() {
					mDescriptor = &StructMethod{
						Descriptor: s,
					}
					t.Methods = append(t.Methods, mDescriptor)
				}
			}
		}
	} else {
		rt.Type = tp
	}
}

// NewStruct return new pointer Struct
func NewStruct(pkg *Package, name string) *Struct {
	srct := &Struct{
		pkg:  pkg,
		name: name,
	}
	pkg.AppendStruct(srct)
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
func (s *Struct) FormatComments() string {
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
