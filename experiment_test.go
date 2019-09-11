package myasthurts_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts", func() {
	It("should test somethhing", func() {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "data/models1.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(file).ToNot(BeNil())

		// fmt.Printf("file: %v\n\n---\n\n", file)
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
	})
})
