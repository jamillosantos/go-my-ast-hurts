package myasthurts

import (
	"fmt"
	"go/ast"
)

func Parse(file *ast.File, definitions *Environment) {

	p := &Package{}
	p.Name = file.Name.Name

	decls := file.Decls
	for _, d := range decls {
		switch c := d.(type) {
		case *ast.GenDecl:
			parseGenDecl(p, c)
		case *ast.FuncDecl:
			parseFuncDecl(p, c)
		}
	}

	definitions.Packages = append(definitions.Packages, p)
}

func parseGenDecl(parent *Package, s *ast.GenDecl) {
	for _, spec := range s.Specs {
		parseSpec(parent, spec)
	}
}

func parseSpec(parent *Package, spec ast.Spec) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := &Struct{}
			declStruct.Name = s.Name.Name

			fmt.Println(s.Comment.Text())

			parseStruct(parent, t, declStruct)
			parent.Structs = append(parent.Structs, declStruct)
		}
	}
}

func parseStruct(parent *Package, astStruct *ast.StructType, typeStruct *Struct) {

	for _, field := range astStruct.Fields.List {
		f := &Field{}
		//tp := &TagParam{}

		if len(field.Names) > 0 { // TODO(jack): To check/understand multiple names.
			f.Name = field.Names[0].Name
		}

		// TODO(jack): To parse the type.
		// f.Type

		// TODO(jack): To parse the tag.
		f.Tag.Raw = field.Tag.Value

		f.Comment = field.Comment.Text()

		typeStruct.Fields = append(typeStruct.Fields, f)
	}
}

func parseFuncDecl(parent *Package, f *ast.FuncDecl) {
	method := &MethodDescriptor{}
	method.Name = f.Name.Name

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
		argument.Type = t.Name
		method.Arguments = append(method.Arguments, argument)
	}
	parent.Methods = append(parent.Methods, method)
}
