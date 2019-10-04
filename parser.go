package myasthurts

import (
	"fmt"
	"go/ast"

	"github.com/fatih/structtag"
)

type parseContext struct {
	File        *ast.File
	Env         *Environment
	PackagesMap map[string]*Package
}

func (p *parseContext) PackageByName(name string) (*Package, bool) {
	pkg, ok := p.PackagesMap[name]
	return pkg, ok
}

func Parse(env *Environment, file *ast.File) {

	ctx := &parseContext{
		File:        file,
		Env:         env,
		PackagesMap: make(map[string]*Package),
	}

	parsePackage(ctx)

	decls := file.Decls
	for _, d := range decls {
		switch c := d.(type) {
		case *ast.GenDecl:
			parseGenDecl(ctx, c)
		case *ast.FuncDecl:
			parseFuncDecl(ctx, c)
		}
	}
}

func parsePackage(ctx *parseContext) {

	_, ok := ctx.Env.PackageByName(ctx.File.Name.Name)
	if !ok {
		var comments []string
		if ctx.File.Doc != nil {
			for _, t := range ctx.File.Comments {
				parseComments(t, &comments)
			}
		}
		pkg := &Package{
			Name:    ctx.File.Name.Name,
			Comment: comments,
		}
		ctx.Env.AppendPackage(pkg)
		ctx.PackagesMap[ctx.File.Name.Name] = pkg
	}
}

func parseComments(doc *ast.CommentGroup, c *[]string) {
	if doc.List == nil {
		return
	}

	sizeList := len(doc.List)
	if len(*c) != 0 {
		t := make([]string, sizeList)
		for i := 0; i < sizeList; i++ {
			t[i] = doc.List[i].Text
		}
		*c = append(*c, t...)
		return
	}
	*c = make([]string, sizeList)
	for i := 0; i < sizeList; i++ {
		(*c)[i] = doc.List[i].Text
	}
}

func parseGenDecl(ctx *parseContext, s *ast.GenDecl) {
	var comments []string
	if s.Doc != nil {
		parseComments(s.Doc, &comments)
	}

	for _, spec := range s.Specs {
		parseSpec(ctx, spec, &comments)
	}
}

func parseSpec(ctx *parseContext, spec ast.Spec, comments *[]string) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]

	switch s := spec.(type) {
	case *ast.TypeSpec:
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := NewStruct(currentPackage, s.Name.Name)
			declStruct.Comment = *comments
			refType := getRefType(ctx, s.Name.Name)
			refType.AppendType(declStruct)

			parseStruct(ctx, t, declStruct)
			currentPackage.Types = append(currentPackage.Types, declStruct)
		}
	case *ast.ImportSpec:
		namePkg := s.Path.Value[1 : len(s.Path.Value)-1]
		if s.Name != nil {
			_, ok := ctx.PackageByName(s.Name.Name)
			if !ok {
				newPkg := &Package{
					Name: namePkg,
				}
				ctx.PackagesMap[s.Name.Name] = newPkg
				_, ok := ctx.Env.PackageByName(namePkg)
				if !ok {
					ctx.Env.AppendPackage(newPkg)
				}
			}
			getRefType(ctx, s.Name.Name)
		} else {
			_, ok := ctx.PackageByName(namePkg)
			if !ok {
				newPkg := &Package{
					Name: namePkg,
				}
				ctx.PackagesMap[namePkg] = newPkg
				_, ok := ctx.Env.PackageByName(namePkg)
				if !ok {
					ctx.Env.AppendPackage(newPkg)
				}
			}
			getRefType(ctx, namePkg)
		}
	case *ast.ValueSpec:
		parseVariable(currentPackage, s)
	}
}

func parseStruct(ctx *parseContext, astStruct *ast.StructType, typeStruct *Struct) {

	for _, field := range astStruct.Fields.List {
		refType := &RefType{}

		switch t := field.Type.(type) {
		case *ast.Ident:
			refType = getRefType(ctx, t.Name)
		case *ast.SelectorExpr:
			refType = getRefType(ctx, t.X.(*ast.Ident).Name)
		case *ast.StarExpr:
			refType = getRefType(ctx, t.X.(*ast.Ident).Name)
		}

		f := &Field{}
		f.Type = refType
		f.Tag.Raw = ""
		if field.Doc != nil {
			var comments []string
			parseComments(field.Doc, &comments)
			f.Comment = comments
		}

		if len(field.Names) > 0 { // TODO(jack): To check/understand multiple names.
			f.Name = field.Names[0].Name
		}

		if field.Tag != nil && field.Tag.Value != "" {
			f.Tag.Raw = field.Tag.Value[1 : len(field.Tag.Value)-1]

			structTag, err := structtag.Parse(f.Tag.Raw)
			if err != nil {
				fmt.Println("Error in format StructTag.")
				panic(err)
			}

			for _, tag := range structTag.Tags() {
				tp := &TagParam{}
				size := len(tag.Options)

				tp.Name = tag.Key
				tp.Value = tag.Name
				tp.Options = make([]string, size)

				if size != 0 {
					for i := 0; i < size; i++ {
						tp.Options[i] = tag.Options[i]
					}
				}
				f.Tag.AppendTagParam(tp)
			}
		}
		typeStruct.Fields = append(typeStruct.Fields, f)
	}
}

func parseFuncDecl(ctx *parseContext, f *ast.FuncDecl) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]
	method := NewMethodDescriptor(currentPackage, f.Name.Name)

	if f.Recv != nil { //TODO(check): I don't know if this pointers always will exist if "f.Recv" is diff than nil.
		recvList := f.Recv.List
		for _, field := range recvList {
			recv := MethodArgument{}
			typeName := field.Type.(*ast.StarExpr).X.(*ast.Ident).Name

			recv.Name = field.Names[0].Name
			recv.Type = getRefType(ctx, typeName)
			method.Recv = append(method.Recv, recv)
		}
	}

	if f.Doc != nil {
		var comments []string
		parseComments(f.Doc, &comments)
		method.Comment = comments
	}

	for _, field := range f.Type.Params.List {
		argument := MethodArgument{}
		if len(field.Names) > 0 {
			argument.Name = field.Names[0].Name
		}

		switch t := field.Type.(type) {
		case *ast.Ident:
			argument.Type = getRefType(ctx, t.Name)
		case *ast.StarExpr:
			argument.Type = getRefType(ctx, t.X.(*ast.Ident).Name)
		}

		method.Arguments = append(method.Arguments, argument)
	}
	currentPackage.Methods = append(currentPackage.Methods, method)
}

func parseVariable(parent *Package, f *ast.ValueSpec) {
	variable := &Variable{}
	varType := NewRefType(parent)

	variable.Name = f.Names[0].Name

	if f.Names == nil {
		varType.Name = f.Type.(*ast.Ident).Name
		variable.Type = varType
	} else {
		for _, value := range f.Values {
			switch v := value.(type) {
			case *ast.BasicLit:
				//varType.Name = TODO: Set values
				fmt.Printf("%T\n", v.Kind) //TODO: Convert token.token to string
			case *ast.Ident:
				//varType.Name = TODO: Set values
				fmt.Printf("%T\n", v.Name)
			}
		}
	}
}

//This method is temp.
func getRefType(ctx *parseContext, name string) *RefType {
	packageCurrent := ctx.PackagesMap[ctx.File.Name.Name]

	for _, pr := range packageCurrent.RefType {
		if name == pr.Name && packageCurrent == pr.Pkg {
			return pr
		}
	}
	return newRefType(packageCurrent, name)
}

func newRefType(parent *Package, name string) *RefType {
	ref := &RefType{
		Pkg:  parent,
		Name: name,
	}
	parent.RefType = append(parent.RefType, ref)
	return ref
}
