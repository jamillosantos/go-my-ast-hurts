package myasthurts

import (
	"fmt"
	"go/ast"

	"github.com/fatih/structtag"
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
			declStruct := NewStruct(parent, s.Name.Name)
			refType := getRefType(parent, s.Name.Name)

			parseStruct(parent, t, declStruct)

			refType.Type = append(refType.Type, declStruct)
			parent.Structs = append(parent.Structs, declStruct)
			parent.Types = append(parent.Types, declStruct)
		}
	case *ast.ValueSpec:
		parseVariable(parent, s)
	}
}

func parseStruct(parent *Package, astStruct *ast.StructType, typeStruct *Struct) {

	for _, field := range astStruct.Fields.List {
		refType := &RefType{}

		switch t := field.Type.(type) {
		case *ast.Ident:
			refType = getRefType(parent, t.Name)
		case *ast.SelectorExpr:
			refType = getRefType(parent, t.X.(*ast.Ident).Name)
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

func parseFuncDecl(parent *Package, f *ast.FuncDecl) {
	method := NewMethodDescriptor(parent, f.Name.Name)

	if f.Recv != nil {
		structMethod := &StructMethod{}
		recvList := f.Recv.List
		for _, field := range recvList {
			recv := MethodArgument{}
			typeName := field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
			refType := getRefType(parent, typeName)
			refType.Type = append(refType.Type, method)

			recv.Name = field.Names[0].Name
			recv.Type = refType //TODO(check): I don't know if this pointers always will exist if "f.Recv" is diff than nil.
			method.Recv = append(method.Recv, recv)
		}

		structMethod.Descriptor = method
		for _, s := range parent.Structs { //TODO(enhancement): Is possible than the Struct has not been read before func.
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
		argument.Type = getRefType(parent, t.Name)
		method.Arguments = append(method.Arguments, argument)
	}
	parent.Methods = append(parent.Methods, method)
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
func getRefType(parent *Package, name string) *RefType {

	for _, pr := range parent.RefType {
		if name == pr.Name {
			for _, t := range pr.Pkg {
				if parent == t {
					return pr
				}
			}
		}
	}
	ref := &RefType{
		Name: name,
	}
	ref.Pkg = append(ref.Pkg, parent)
	return ref
}
