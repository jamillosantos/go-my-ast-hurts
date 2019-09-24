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
		ctx.Env.AppendPackage(&Package{
			Name: ctx.File.Name.Name,
		})

		ctx.PackagesMap[ctx.File.Name.Name] = &Package{
			Name: ctx.File.Name.Name,
		}

	}
}

func parseGenDecl(ctx *parseContext, s *ast.GenDecl) {
	for _, spec := range s.Specs {
		parseSpec(ctx, spec)
	}
}

func parseSpec(ctx *parseContext, spec ast.Spec) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]

	switch s := spec.(type) {
	case *ast.TypeSpec:
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := NewStruct(currentPackage, s.Name.Name)
			refType := getRefType(ctx, s.Name.Name)
			refType.Type = declStruct

			parseStruct(ctx, t, declStruct)
			currentPackage.Types = append(currentPackage.Types, declStruct)
		}
	case *ast.ImportSpec:
		pkgName := s.Path.Value[1 : len(s.Path.Value)-1]
		if s.Name != nil {
			pkgName += fmt.Sprintf(":%s", s.Name.Name)
		}
		getRefType(ctx, pkgName)
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
			//refType.Name = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
		}

		f := &Field{}
		f.Type = refType
		f.Tag.Raw = ""
		f.Comment = field.Comment.Text()

		if len(field.Names) > 0 { // TODO(jack): To check/understand multiple names.
			f.Name = field.Names[0].Name
		}

		if field.Tag != nil && field.Tag.Value != "" {
			f.Tag.Raw = field.Tag.Value[1 : len(field.Tag.Value)-1]
			tp := &TagParam{}

			structTag, err := structtag.Parse(f.Tag.Raw)
			if err != nil {
				fmt.Println("Error in format StructTag.")
				panic(err)
			}

			jsonTag, err := structTag.Get("json")
			if err != nil {
				fmt.Println("Error in parse StructTag.")
				panic(err)
			}

			size := len(jsonTag.Options)

			tp.Name = jsonTag.Key
			tp.Value = jsonTag.Name
			tp.Options = make([]string, size)

			if size != 0 {
				for i := 0; i < size; i++ {
					tp.Options[i] = jsonTag.Options[i]
				}
			}
			f.Tag.Params = append(f.Tag.Params, *tp)
		}
		typeStruct.Fields = append(typeStruct.Fields, f)
	}
}

func parseFuncDecl(ctx *parseContext, f *ast.FuncDecl) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]
	method := NewMethodDescriptor(currentPackage, f.Name.Name)

	if f.Recv != nil { //TODO(check): I don't know if this pointers always will exist if "f.Recv" is diff than nil.
		structMethod := &StructMethod{}
		recvList := f.Recv.List
		for _, field := range recvList {
			recv := MethodArgument{}
			typeName := field.Type.(*ast.StarExpr).X.(*ast.Ident).Name

			recv.Name = field.Names[0].Name
			recv.Type = getRefType(ctx, typeName)
			method.Recv = append(method.Recv, recv)
		}

		structMethod.Descriptor = method
		for _, s := range currentPackage.Structs { //TODO(enhancement): Is possible than the Struct has not been read before func.
			if s.Name() == method.Recv[0].Type.Name {
				s.Methods = append(s.Methods, structMethod)
				break
			}
		}

	}

	for _, field := range f.Type.Params.List {
		argument := MethodArgument{}
		if len(field.Names) > 0 {
			argument.Name = field.Names[0].Name
		}

		t, ok := field.Type.(*ast.Ident)
		if !ok {
			fmt.Println("Treta")
			continue
		}
		argument.Type = getRefType(ctx, t.Name)
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
func getRefType(ctx *parseContext, nameRef string) *RefType {
	packageCurrent := ctx.PackagesMap[ctx.File.Name.Name]

	for _, pr := range packageCurrent.RefType {
		if nameRef == pr.Name && packageCurrent == pr.Pkg {
			return pr
		}
	}
	return newRefType(packageCurrent, nameRef)
}

func newRefType(parent *Package, name string) *RefType {
	fmt.Println(name)
	ref := &RefType{
		Name: name,
		Pkg:  parent,
	}
	parent.RefType = append(parent.RefType, ref)
	return ref
}
