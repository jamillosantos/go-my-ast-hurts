package myasthurts_test

import (
	"go/parser"
	"go/token"

	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts", func() {

	Context("should parse struct", func() {

		It("should parse struct", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models1.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg).NotTo(BeNil())
			Expect(pkg.Structs).To(HaveLen(2))

		})

		It("should parse struct fields", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models2.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields).To(HaveLen(4))

			Expect(pkg.Structs[0].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[0].Fields[1].Name).To(Equal("Name"))

			Expect(pkg.Structs[0].Fields[0].Type).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[0].Fields[1].Type).ToNot(Equal(BeNil()))

			Expect(pkg.Structs[0].Fields[0].Type.Name).To(Equal("int64"))
			Expect(pkg.Structs[0].Fields[1].Type.Name).To(Equal("string"))

			Expect(pkg.Structs[0].Fields[0].Type.Type).To(BeNil())
			Expect(pkg.Structs[0].Fields[1].Type.Type).To(BeNil())

			Expect(pkg.Structs[1].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[1].Fields[1].Name).To(Equal("Address"))
			Expect(pkg.Structs[1].Fields[2].Name).To(Equal("User"))

			Expect(pkg.Structs[1].Fields[0].Type).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[1].Fields[1].Type).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[1].Fields[2].Type).ToNot(Equal(BeNil()))

			Expect(pkg.Structs[1].Fields[0].Type.Name).To(Equal("int64"))
			Expect(pkg.Structs[1].Fields[1].Type.Name).To(Equal("string"))
			Expect(pkg.Structs[1].Fields[2].Type.Name).To(Equal("User"))
			Expect(pkg.Structs[1].Fields[3].Type.Name).To(Equal("User"))

			Expect(pkg.Structs[1].Fields[0].Type.Type).To(BeNil())
			Expect(pkg.Structs[1].Fields[1].Type.Type).To(BeNil())
			Expect(pkg.Structs[1].Fields[2].Type.Type).To(Equal(pkg.Structs[0]))
			Expect(pkg.Structs[1].Fields[3].Type.Type).To(Equal(pkg.Structs[0]))

		})

		It("should parse struct tags", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models3.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Name()).To(Equal("User"))
			Expect(pkg.Structs[0].Fields).To(HaveLen(3))

			Expect(pkg.Structs[0].Fields[0].Tag).NotTo(BeNil())
			Expect(pkg.Structs[0].Fields[0].Tag.Raw).To(Equal("json:\"id\""))
			Expect(pkg.Structs[0].Fields[0].Tag.Params).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[0].Tag.Params[0].Name).To(Equal("json"))
			Expect(pkg.Structs[0].Fields[0].Tag.Params[0].Value).To(Equal("id"))
			Expect(pkg.Structs[0].Fields[0].Tag.Params[0].Options).To(BeEmpty())

			Expect(pkg.Structs[0].Fields[1].Tag.Raw).To(Equal("json:\"name,lastName\""))
			Expect(pkg.Structs[0].Fields[1].Tag.Params).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Name).To(Equal("json"))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Value).To(Equal("name"))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Options).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Options[0]).To(Equal("lastName"))

			Expect(pkg.Structs[0].Fields[2].Tag.Raw).To(Equal("json:\"address\" bson:\"\""))
			Expect(pkg.Structs[0].Fields[2].Tag.Params).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[0].Name).To(Equal("json"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[0].Value).To(Equal("address"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[0].Options).To(BeEmpty())

			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Name).To(Equal("bson"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Value).To(BeEmpty())
			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Options).To(BeEmpty())

		})

		It("should parse struct custom field", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models4.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[1].Type.Name).To(Equal("User"))
			Expect(pkg.Structs[0].Fields[1].Type.Type).To(Equal(pkg.Structs[1]))
			Expect(pkg.Structs[0].Fields[1].Type.Pkg).To(Equal(pkg))

		})

		It("should parse struct comments", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models5.sample", nil, parser.ParseComments)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Comment).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields).To(HaveLen(4))

			Expect(pkg.Structs[0].Fields[0].Comment).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Comment).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[2].Comment).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[3].Comment).To(BeEmpty())

		})

	})

	Context("should parse function", func() {

		It("should parse function", func() {

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models6.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg.Name).To(Equal("models"))

			Expect(pkg.Methods).To(HaveLen(3))
			Expect(pkg.Methods[0].Name()).To(Equal("getName"))
			Expect(pkg.Methods[0].Arguments).To(BeEmpty())
			Expect(pkg.Methods[0].Recv).To(HaveLen(1))

			Expect(pkg.Methods[1].Name()).To(Equal("show"))
			Expect(pkg.Methods[1].Arguments).To(HaveLen(2))
			Expect(pkg.Methods[1].Recv).To(BeEmpty())

			Expect(pkg.Methods[1].Arguments[0].Name).To(Equal("name"))
			Expect(pkg.Methods[1].Arguments[1].Name).To(Equal("age"))

			Expect(pkg.Methods[1].Arguments[0].Type.Type).To(BeNil())
			Expect(pkg.Methods[1].Arguments[1].Type.Type).To(BeNil())

			Expect(pkg.Methods[1].Arguments[0].Type.Name).To(Equal("string"))
			Expect(pkg.Methods[1].Arguments[1].Type.Name).To(Equal("int64"))

			Expect(pkg.Methods[2].Name()).To(Equal("welcome"))
			Expect(pkg.Methods[2].Arguments).To(BeEmpty())
			Expect(pkg.Methods[2].Recv).To(BeEmpty())

		})

		It("should parse function in Struct", func() {

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models7.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Methods).To(HaveLen(2))

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Name()).To(Equal("getName"))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Arguments).To(HaveLen(0))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Recv).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Recv[0].Type.Type).To(Equal(pkg.Structs[0]))

			Expect(pkg.Methods[0].Name()).To(Equal("getName"))
			Expect(pkg.Methods[1].Name()).To(Equal("getName_"))

		})
	})

	Context("should parse imports", func() {

		It("should parse imports", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models8.sample", nil, parser.AllErrors)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			models, _ := env.PackageByName("models")
			fmt, _ := env.PackageByName("fmt")
			time, _ := env.PackageByName("time")

			Expect(models).ToNot(BeNil())
			Expect(fmt).ToNot(BeNil())
			Expect(time).ToNot(BeNil())

			Expect(models.Name).To(Equal("models"))
			Expect(fmt.Name).To(Equal("fmt"))
			Expect(time.Name).To(Equal("time"))

		})

	})

	Context("should parse comments", func() {

		It("should parse multilines or no in struct comments", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models9.sample", nil, parser.ParseComments)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Comment).To(HaveLen(8))
			Expect(pkg.Comment[0]).To(Equal("// Package models is a test"))

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Comment).To(HaveLen(1))
			Expect(pkg.Structs[0].Comment[0]).To(Equal("/* Testing comment\nnew line\n*/"))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[0].Comment).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[0].Comment[0]).To(Equal("// ID comment"))
			Expect(pkg.Structs[0].Fields[1].Comment).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Comment[0]).To(Equal("/* Comment with multilines\n\t\tTesting\n\t*/"))

			Expect(pkg.Structs[1].Fields).To(HaveLen(3))
			Expect(pkg.Structs[1].Fields[1].Comment).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields[1].Comment[0]).To(Equal("// Line 1"))
			Expect(pkg.Structs[1].Fields[1].Comment[1]).To(Equal("// Line 2"))

		})

		It("should parse comments in func", func() {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "data/models10.sample", nil, parser.ParseComments)
			Expect(err).ToNot(HaveOccurred())
			Expect(f).ToNot(BeNil())
			Expect(f.Decls).ToNot(BeNil())

			//ast.Print(fset, f)
			env := myasthurts.NewEnvironment()
			myasthurts.Parse(env, f)

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Comment).To(HaveLen(7))
			Expect(pkg.Comment[0]).To(Equal("// Package models is a test"))

			Expect(pkg.Methods).To(HaveLen(3))
			Expect(pkg.Methods[0].Comment).To(HaveLen(1))
			Expect(pkg.Methods[0].Comment[0]).To(Equal("// Comment here"))
			Expect(pkg.Methods[1].Comment).To(HaveLen(1))
			Expect(pkg.Methods[1].Comment[0]).To(Equal("/** Description \n    multilines\n*/"))

			//Expect(pkg.Comment).To(HaveLen(6))

		})

	})

})
