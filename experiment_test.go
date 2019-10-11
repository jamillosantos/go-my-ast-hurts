package myasthurts_test

import (
	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts - Parse simples files with tags and func from struct", func() {

	When("parsing struct", func() {

		It("should check two struct in file", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models1.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg).NotTo(BeNil())
			Expect(pkg.Structs).To(HaveLen(2))

		})

		It("should check struct fields", func() {
			env, exrr := myasthurts.NewEnvironment()
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

		It("should check struct tags", func() {
			env, exrr := myasthurts.NewEnvironment()
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

		It("should check struct custom field with user", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models4.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[1].Type.Name).To(Equal("User"))
			Expect(pkg.Structs[0].Fields[1].Type.Type).To(Equal(pkg.Structs[1]))
			Expect(pkg.Structs[0].Fields[1].Type.Pkg).To(Equal(pkg))

		})

		It("should check struct comments", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models5.sample", true)
			Expect(exrr).To(BeNil())

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

	When("parsing variables", func() {

		PIt("should check variables names", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models11.sample", true)
			Expect(exrr).To(BeNil())

			// TODO(Jeconias):
		})

	})

	When("parsing function", func() {

		It("should check name and parameters from function", func() {

			env, exrr := myasthurts.NewEnvironment()
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
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models7.sample", true)
			Expect(exrr).To(BeNil())

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

	When("parsing imports", func() {

		It("should check names all imports", func() {

			env, exrr := myasthurts.NewEnvironment()
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
			env, exrr := myasthurts.NewEnvironment()
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

			ref, _ := bytes.RefTypeByName("Buffer")
			Expect(ref).ToNot(BeNil())
			Expect(models.Methods[0].Arguments[0].Type.Name).To(Equal(ref.Type.Name()))

			stct := bytes.StructByName("Buffer")
			Expect(stct).ToNot(BeNil())
			Expect(models.Methods[0].Arguments[0].Type.Type).To(Equal(stct))
		})
	})

	When("parsing comments", func() {

		It("should check multilines or no in struct comments", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models9.sample", true)
			Expect(exrr).To(BeNil())

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

		It("should check comments in func", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models10.sample", true)
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("models")
			Expect(ok).To(BeTrue())

			Expect(pkg.Comment).To(HaveLen(7))
			Expect(pkg.Comment[0]).To(Equal("// Package models is a test"))

			Expect(pkg.Methods).To(HaveLen(3))
			Expect(pkg.Methods[0].Comment).To(HaveLen(1))
			Expect(pkg.Methods[0].Comment[0]).To(Equal("// Comment here"))
			Expect(pkg.Methods[1].Comment).To(HaveLen(1))
			Expect(pkg.Methods[1].Comment[0]).To(Equal("/** Description \n    multilines\n*/"))

		})

	})

	When("parsing builtin file", func() {

		It("should check types from builtin file", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParsePackage("data/models13.sample", true)
			Expect(exrr).To(BeNil())

			pkgM, okM := env.PackageByName("models")
			pkgB, okB := env.PackageByName("builtin")
			Expect(okM).To(BeTrue())
			Expect(okB).To(BeTrue())

			T1, _ := pkgB.RefTypeByName("string")
			Expect(T1).ToNot(BeNil())

			T2, _ := pkgB.RefTypeByName("int")
			Expect(T2).ToNot(BeNil())

			T3, _ := pkgB.RefTypeByName("int8")
			Expect(T3).ToNot(BeNil())

			T4, _ := pkgB.RefTypeByName("int16")
			Expect(T4).ToNot(BeNil())

			T5, _ := pkgB.RefTypeByName("int32")
			Expect(T5).ToNot(BeNil())

			T6, _ := pkgB.RefTypeByName("int64")
			Expect(T6).ToNot(BeNil())

			T7, _ := pkgB.RefTypeByName("float32")
			Expect(T7).ToNot(BeNil())

			T8, _ := pkgB.RefTypeByName("float64")
			Expect(T8).ToNot(BeNil())

			T9, _ := pkgB.RefTypeByName("byte")
			Expect(T9).ToNot(BeNil())

			Expect(pkgM.Structs).To(HaveLen(1))
			Expect(pkgM.Structs[0].Fields).To(HaveLen(9))

			Expect(pkgM.Structs[0].Fields[0].Type).To(Equal(T1))
			Expect(pkgM.Structs[0].Fields[1].Type).To(Equal(T2))
			Expect(pkgM.Structs[0].Fields[2].Type).To(Equal(T3))
			Expect(pkgM.Structs[0].Fields[3].Type).To(Equal(T4))
			Expect(pkgM.Structs[0].Fields[4].Type).To(Equal(T5))
			Expect(pkgM.Structs[0].Fields[5].Type).To(Equal(T6))
			Expect(pkgM.Structs[0].Fields[6].Type).To(Equal(T7))
			Expect(pkgM.Structs[0].Fields[7].Type).To(Equal(T8))
			Expect(pkgM.Structs[0].Fields[8].Type).To(Equal(T9))

		})

		It("should check builtin file", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByName("builtin")
			Expect(ok).To(BeTrue())

			_, exrr = pkg.RefTypeByName("string")
			Expect(exrr).To(BeNil())

			_, exrr = pkg.RefTypeByName("int64")
			Expect(exrr).To(BeNil())

			_, exrr = pkg.RefTypeByName("float32")
			Expect(exrr).To(BeNil())

			// ### WORK IN PROGRESS ###

		})

	})

})
