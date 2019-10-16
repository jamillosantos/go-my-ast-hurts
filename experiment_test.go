package myasthurts

import (
	//myasthurts "github.com/lab259/go-my-ast-hurts"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts - Parse simples files with tags and func from struct", func() {

	When("parsing struct", func() {

		It("should check two struct in file", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models1.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg).NotTo(BeNil())
			Expect(pkg.Structs).To(HaveLen(2))
		})

		It("should check struct fields", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models2.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields).To(HaveLen(4))

			Expect(pkg.Structs[0].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[0].Fields[1].Name).To(Equal("Name"))

			Expect(pkg.Structs[0].Fields[0].RefType).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[0].Fields[1].RefType).ToNot(Equal(BeNil()))

			Expect(pkg.Structs[0].Fields[0].RefType.Name).To(Equal("int64"))
			Expect(pkg.Structs[0].Fields[1].RefType.Name).To(Equal("string"))

			Expect(pkg.Structs[0].Fields[0].RefType.Type).To(BeNil())
			Expect(pkg.Structs[0].Fields[1].RefType.Type).To(BeNil())

			Expect(pkg.Structs[1].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[1].Fields[1].Name).To(Equal("Address"))
			Expect(pkg.Structs[1].Fields[2].Name).To(Equal("User"))

			Expect(pkg.Structs[1].Fields[0].RefType).ToNot(BeNil())
			Expect(pkg.Structs[1].Fields[1].RefType).ToNot(BeNil())
			Expect(pkg.Structs[1].Fields[2].RefType).ToNot(BeNil())

			Expect(pkg.Structs[1].Fields[0].RefType.Name).To(Equal("int64"))
			Expect(pkg.Structs[1].Fields[1].RefType.Name).To(Equal("string"))
			Expect(pkg.Structs[1].Fields[2].RefType.Name).To(Equal("User"))
			Expect(pkg.Structs[1].Fields[3].RefType.Name).To(Equal("User"))

			Expect(pkg.Structs[1].Fields[0].RefType.Type).To(BeNil())
			Expect(pkg.Structs[1].Fields[1].RefType.Type).To(BeNil())
			Expect(pkg.Structs[1].Fields[2].RefType.Type).To(Equal(pkg.Structs[0]))
			Expect(pkg.Structs[1].Fields[3].RefType.Type).To(Equal(pkg.Structs[0]))
		})

		It("should check struct tags", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models3.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Name()).To(Equal("User"))
			Expect(pkg.Structs[0].Fields).To(HaveLen(3))

			Expect(pkg.Structs[0].Fields[0].Tag).NotTo(BeNil())
			Expect(pkg.Structs[0].Fields[0].Tag.Raw).To(Equal(`json:"id"`))
			Expect(pkg.Structs[0].Fields[0].Tag.Params).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[0].Tag.Params[0].Name).To(Equal("json"))
			Expect(pkg.Structs[0].Fields[0].Tag.Params[0].Value).To(Equal("id"))
			Expect(pkg.Structs[0].Fields[0].Tag.Params[0].Options).To(BeEmpty())

			Expect(pkg.Structs[0].Fields[1].Tag.Raw).To(Equal(`json:"name,lastName"`))
			Expect(pkg.Structs[0].Fields[1].Tag.Params).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Name).To(Equal("json"))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Value).To(Equal("name"))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Options).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Tag.Params[0].Options[0]).To(Equal("lastName"))

			Expect(pkg.Structs[0].Fields[2].Tag.Raw).To(Equal(`json:"address" bson:""`))
			Expect(pkg.Structs[0].Fields[2].Tag.Params).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[0].Name).To(Equal("json"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[0].Value).To(Equal("address"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[0].Options).To(BeEmpty())

			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Name).To(Equal("bson"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Value).To(BeEmpty())
			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Options).To(BeEmpty())
		})

		It("should check struct custom field with user", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models4.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[1].RefType.Name).To(Equal("User"))
			Expect(pkg.Structs[0].Fields[1].RefType.Type).To(Equal(pkg.Structs[1]))
			Expect(pkg.Structs[0].Fields[1].RefType.Pkg).To(Equal(pkg))

		})

		It("should check struct comments", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models5.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields).To(HaveLen(4))

			Expect(pkg.Structs[0].Fields[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[2].Doc.Comments).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[3].Doc.Comments).To(BeEmpty())

		})

	})

	When("parsing variables", func() {

		It("should check variables declarated with var", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())
			env.Config = EnvConfig{
				DevMode: true,
				ASTI:    false,
			}

			exrr = env.ParsePackage("data/models11.sample", true)
			builtin, _ := env.PackageByName("builtin")
			Expect(builtin).ToNot(BeNil())
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Variables).To(HaveLen(7))

			a := pkg.VariableByName("a")
			Expect(a).ToNot(BeNil())
			Expect(a.RefType).ToNot(BeNil())

			b := pkg.VariableByName("b")
			Expect(b).ToNot(BeNil())
			Expect(b.RefType).ToNot(BeNil())

			c := pkg.VariableByName("c")
			Expect(c).ToNot(BeNil())
			Expect(c.RefType).ToNot(BeNil())

			d := pkg.VariableByName("d")
			Expect(d).ToNot(BeNil())
			Expect(d.RefType).ToNot(BeNil())

			e := pkg.VariableByName("e")
			Expect(e).ToNot(BeNil())
			Expect(e.RefType).ToNot(BeNil())

			f := pkg.VariableByName("f")
			Expect(f).ToNot(BeNil())
			Expect(f.RefType).ToNot(BeNil())

			g := pkg.VariableByName("g")
			Expect(g).ToNot(BeNil())
			Expect(g.RefType).ToNot(BeNil())

			Expect(pkg.RefTypeByName("string")).To(Equal(a.RefType))
			Expect(pkg.RefTypeByName("byte")).To(Equal(b.RefType))
			Expect(pkg.RefTypeByName("int")).To(Equal(c.RefType))
			Expect(pkg.RefTypeByName("int64")).To(Equal(d.RefType))
			Expect(pkg.RefTypeByName("float32")).To(Equal(e.RefType))
			Expect(pkg.RefTypeByName("boolean")).To(Equal(f.RefType))
			Expect(pkg.RefTypeByName("User")).To(Equal(g.RefType))
			Expect(g.RefType.Type).ToNot(BeNil())
			Expect(g.RefType.Type.Name()).To(Equal("User"))

		})

	})

	When("parsing function", func() {

		It("should check name and parameters from function", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models6.sample", true)
			Expect(exrr).To(BeNil())

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

		It("should check if func belong to struct", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models7.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Methods).To(HaveLen(2))

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Name()).To(Equal("getName"))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Arguments).To(BeEmpty())
			Expect(pkg.Structs[0].Methods[0].Descriptor.Recv).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods[0].Descriptor.Recv[0].Type.Type).To(Equal(pkg.Structs[0]))

			Expect(pkg.Methods[0].Name()).To(Equal("getName"))
			Expect(pkg.Methods[1].Name()).To(Equal("getName_"))

		})
	})

	When("parsing imports", func() {

		It("should check names all imports", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models8.sample", true)
			Expect(exrr).To(BeNil())

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

		It("should check the struct types of import package bytes with dot", func() {
			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models12.sample", true)
			Expect(exrr).To(BeNil())

			models, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			bytes, ok := env.PackageByName("bytes")
			Expect(ok).To(BeTrue())

			Expect(models.Methods).To(HaveLen(2))
			Expect(models.Methods[0].Name()).To(Equal("welcome"))
			Expect(models.Methods[0].Arguments).To(HaveLen(1))
			Expect(models.Methods[0].Arguments[0].Name).To(Equal("buf"))

			ref := bytes.RefTypeByName("Buffer")
			Expect(ref).ToNot(BeNil())
			Expect(models.Methods[0].Arguments[0].Type.Name).To(Equal(ref.Type.Name()))

			stct := bytes.StructByName("Buffer")
			Expect(stct).ToNot(BeNil())
			Expect(models.Methods[0].Arguments[0].Type.Type).To(Equal(stct))
		})
	})

	When("parsing comments", func() {

		It("should check multilines or no in struct comments", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models9.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Doc.Comments).To(HaveLen(6))
			Expect(pkg.Doc.Comments[0]).To(Equal("// Package models is a test"))

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Doc.Comments[0]).To(Equal("/* Testing comment\nnew line\n*/"))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[0].Doc.Comments[0]).To(Equal("// ID comment"))
			Expect(pkg.Structs[0].Fields[1].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Fields[1].Doc.Comments[0]).To(Equal("/* Comment with multilines\n\tTesting\n\t*/"))

			Expect(pkg.Structs[1].Fields).To(HaveLen(3))
			Expect(pkg.Structs[1].Fields[1].Doc.Comments).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields[1].Doc.Comments[0]).To(Equal("// Line 1"))
			Expect(pkg.Structs[1].Fields[1].Doc.Comments[1]).To(Equal("// Line 2"))

		})

		It("should remove /*, */ or // from comment", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models14.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Doc.Comments).To(HaveLen(6))

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Doc.FormatComment()).To(Equal("Testing comment\nnew line"))

			Expect(pkg.Structs[1].Fields).To(HaveLen(3))
			Expect(pkg.Structs[1].Fields[1].Doc.Comments).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields[1].Doc.FormatComment()).To(Equal("Line 1\nLine 2\n"))

		})

		It("should check comments in func", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models10.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Doc.Comments).To(HaveLen(7))
			Expect(pkg.Doc.Comments[0]).To(Equal("// Package models is a test"))

			Expect(pkg.Methods).To(HaveLen(3))
			Expect(pkg.Methods[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Methods[0].Doc.Comments[0]).To(Equal("// Comment here"))
			Expect(pkg.Methods[1].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Methods[1].Doc.Comments[0]).To(Equal("/** Description \n    multilines\n*/"))

		})

	})

	When("parsing builtin file", func() {

		It("should check types from builtin file", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models13.sample", true)
			Expect(exrr).To(BeNil())

			pkgM, okM := env.PackageByName("models")
			pkgB, okB := env.PackageByName("builtin")
			Expect(okM).To(BeTrue())
			Expect(okB).To(BeTrue())

			T1 := pkgB.RefTypeByName("string")
			Expect(T1).ToNot(BeNil())

			T2 := pkgB.RefTypeByName("int")
			Expect(T2).ToNot(BeNil())

			T3 := pkgB.RefTypeByName("int8")
			Expect(T3).ToNot(BeNil())

			T4 := pkgB.RefTypeByName("int16")
			Expect(T4).ToNot(BeNil())

			T5 := pkgB.RefTypeByName("int32")
			Expect(T5).ToNot(BeNil())

			T6 := pkgB.RefTypeByName("int64")
			Expect(T6).ToNot(BeNil())

			T7 := pkgB.RefTypeByName("float32")
			Expect(T7).ToNot(BeNil())

			T8 := pkgB.RefTypeByName("float64")
			Expect(T8).ToNot(BeNil())

			T9 := pkgB.RefTypeByName("byte")
			Expect(T9).ToNot(BeNil())

			Expect(pkgM.Structs).To(HaveLen(1))
			Expect(pkgM.Structs[0].Fields).To(HaveLen(9))

			Expect(pkgM.Structs[0].Fields[0].RefType).To(Equal(T1))
			Expect(pkgM.Structs[0].Fields[1].RefType).To(Equal(T2))
			Expect(pkgM.Structs[0].Fields[2].RefType).To(Equal(T3))
			Expect(pkgM.Structs[0].Fields[3].RefType).To(Equal(T4))
			Expect(pkgM.Structs[0].Fields[4].RefType).To(Equal(T5))
			Expect(pkgM.Structs[0].Fields[5].RefType).To(Equal(T6))
			Expect(pkgM.Structs[0].Fields[6].RefType).To(Equal(T7))
			Expect(pkgM.Structs[0].Fields[7].RefType).To(Equal(T8))
			Expect(pkgM.Structs[0].Fields[8].RefType).To(Equal(T9))

		})

		It("should check builtin file", func() {

			env, exrr := NewEnvironment()
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("builtin")
			Expect(ok).To(BeTrue())

			ref := pkg.RefTypeByName("string")
			Expect(ref).ToNot(BeNil())

			ref = pkg.RefTypeByName("int64")
			Expect(ref).ToNot(BeNil())

			ref = pkg.RefTypeByName("float32")
			Expect(ref).ToNot(BeNil())

			// ### WORK IN PROGRESS ###
		})

	})

})
