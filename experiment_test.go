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
		file, err := parser.ParseFile(fset, "data/models1.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(file).ToNot(BeNil())

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
						if len(p.Names) == i {
							params = strings.TrimSuffix(params, ", ")
						}
					}
				}

				fmt.Printf("%s %s %b -> Params(%s)\n", fn.Name.Obj.Kind, fn.Name.Name, fn.Type.Params.NumFields(), params)

			}

			fmt.Print("---------- Reading Functions/end ----------\n\n")
		}

		if file.Decls != nil {
			fmt.Println("---------- Reading Decls/start ----------")
			for _, i := range file.Decls {
				fn, ok := i.(*ast.GenDecl)
				if !ok {
					continue
				}
				fmt.Println(fn.Tok)

			}
			fmt.Println("---------- Reading Decls/end ----------")
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
		fmt.Println()*/
	})
})
