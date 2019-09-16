package myasthurts_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts", func() {

	It("should test somethhing", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models1.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		fmt.Println("---------- Reading ----------")

		ast.Inspect(f, func(x ast.Node) bool {

			switch x.(type) {
			case *ast.StructType:
				s, ok := x.(*ast.StructType)
				if !ok {
					break
				}

				for _, fieldList := range s.Fields.List {
					fmt.Printf("%s %s %s\n", fieldList.Names[0], fieldList.Type, fieldList.Tag.Value)
				}

				/*for _, field := range s.Fields.List {
					typeExpr := field.Type

					start := typeExpr.Pos() - 1
					end := typeExpr.End() - 1

					typeInSource := src[start:end]

					fmt.Println(typeInSource)
				}*/

				fmt.Println(f.Scope.String())
			case *ast.FuncDecl:
				s, ok := x.(*ast.FuncDecl)
				if !ok {
					break
				}

				params := ""
				for _, p := range s.Type.Params.List {
					for i := 0; i < len(p.Names); i++ {
						params += fmt.Sprintf("%s %s, ", p.Names[i], p.Type)
					}
				}
				params = strings.TrimSuffix(params, ", ")
				fmt.Printf("%s %s %b -> Params(%s)\n\n", s.Name.Obj.Kind, s.Name.Name, s.Type.Params.NumFields(), params)
			case *ast.TypeSpec:
				s, ok := x.(*ast.TypeSpec)
				if !ok {
					break
				}
				fmt.Printf("%s %s\n", s.Name.Obj.Kind, s.Name.String())
			}
			return true
		})

	})
})

/*
if file.Imports != nil {
		fmt.Println("---------- Reading Imports/start ----------")

		importsStr := "import (\n"
		for _, i := range file.Imports {
			importsStr += "\t"
			if i == nil {
				fmt.Println("<nil>")
			} else {
				importsStr += fmt.Sprintf("%s %s", i.Name.String(), i.Path.Value)
			}
			importsStr += "\n"
		}
		importsStr += ")"
		fmt.Println(importsStr)
		fmt.Print("---------- Reading Imports/end ----------\n\n")
	}

	if file.Imports != nil {
		fmt.Println("---------- Reading Functions/start ----------")

		for _, i := range file.Decls {
			fn, ok := i.(*ast.FuncDecl)
			if !ok {
				continue
			}

			params := ""
			for _, p := range fn.Type.Params.List {
				for i := 0; i < len(p.Names); i++ {
					params += fmt.Sprintf("%s %s, ", p.Names[i], p.Type)
				}
			}
			params = strings.TrimSuffix(params, ", ")
			fmt.Printf("%s %s %b -> Params(%s)\n", fn.Name.Obj.Kind, fn.Name.Name, fn.Type.Params.NumFields(), params)
		}

		fmt.Print("---------- Reading Functions/end ----------\n\n")
	}

	if file.Decls != nil {
		fmt.Println("---------- Reading Struct/start ----------")

		var v visitor

		ast.Walk(v, file)

		/*ast.Inspect(file, func(x ast.Node) bool {
			s, okk := x.(*ast.StructType)
			if !okk {
				return true
			}
			for _, field := range s.Fields.List {
				fmt.Printf("%s %s %s\n", field.Names[0], field.Type, field.Tag.Value)
			}
			return false
		})

		fmt.Print("---------- Reading Struct/end ----------\n\n")

	}

	/*fmt.Printf("file: %v\n\n---\n\n", file)
	fmt.Println("--------------------")
	fmt.Print("Doc:")
	if file.Doc == nil {
		fmt.Println("<nil>")
	} else {
		fmt.Println(file.Doc.Text())
	}
	fmt.Println("Name:", file.Name.String())
	fmt.Println("Declarations")
	for _, d := range file.Decls {
		// fmt.Printf("%T\n", d)
		switch dcl := d.(type) {
		case *ast.GenDecl:
			fmt.Print("  Doc: ")
			if dcl.Doc == nil {
				fmt.Println("<nil>")
			} else {
				fmt.Println(dcl.Doc.Text())
			}

			for _, s := range dcl.Specs {
				fmt.Printf("%T\n", s)
				switch spec := s.(type) {
				case *ast.ImportSpec:
					fmt.Println("    ", spec.Name.String(), spec.Path.Value)
				}
			}
			// fmt.Println("Name: ", dcl.)
		default:
			fmt.Printf("  unknown type: %T\n", dcl)
		}
	}
	fmt.Println("--------------------")
	fmt.Println()
})*/
