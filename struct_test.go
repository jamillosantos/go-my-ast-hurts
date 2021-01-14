package myasthurts_test

import (
	myasthurts "github.com/lab259/go-my-ast-hurts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Struct", func() {
	Describe("Parse", func() {
		It("should parse a struct property member", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).To(BeNil())

			err = env.ParseFile(newDataPackageContext(env), "data/models14.sample.go")
			Expect(err).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))

			Expect(pkg.Structs[1].Fields).To(HaveLen(5))
			Expect(pkg.Structs[1].Fields[4].RefType.Name()).To(BeEmpty())
			Expect(pkg.Structs[1].Fields[4].RefType.Type()).ToNot(BeNil())
			Expect(pkg.Structs[1].Fields[4].RefType.Type().Name()).To(BeEmpty())
			var sType *myasthurts.Struct
			Expect(pkg.Structs[1].Fields[4].RefType.Type()).To(BeAssignableToTypeOf(sType))
			sType = pkg.Structs[1].Fields[4].RefType.Type().(*myasthurts.Struct)
			Expect(sType.Methods()).To(BeEmpty())
			Expect(sType.Fields).To(HaveLen(1))
			Expect(sType.Fields[0].Name).To(Equal("Name"))
			Expect(sType.Fields[0].RefType.Name()).To(Equal("string"))
		})
	})

	Describe("Implements", func() {
		It("should find struct implements an interface", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			Expect(env.ParseFile(newDataPackageContext(env), "data/interface.go")).To(Succeed())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			s, ok := pkg.StructByName("InterfaceUser")
			Expect(ok).To(BeTrue())
			Expect(s.MethodsMap()).To(HaveKey("Name"))
			Expect(s.MethodsMap()).To(HaveKey("SetName"))

			Expect(pkg.Interfaces).To(HaveLen(3))
			Expect(pkg.Interfaces[0].Name()).To(Equal("HasName"))
			Expect(pkg.Interfaces[1].Name()).To(Equal("HasAge"))
			Expect(pkg.Interfaces[2].Name()).To(Equal("HasNameWrong"))

			Expect(s.Implements(pkg.Interfaces[0])).To(BeTrue())
		})

		It("should not recognize the interface with missing methods", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			Expect(env.ParseFile(newDataPackageContext(env), "data/interface.go")).To(Succeed())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			s, ok := pkg.StructByName("InterfaceUser")
			Expect(ok).To(BeTrue())

			Expect(s.Implements(pkg.Interfaces[1])).To(BeFalse())
		})

		It("should not recognize the interface with incompatible methods", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			Expect(env.ParseFile(newDataPackageContext(env), "data/interface.go")).To(Succeed())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			s, ok := pkg.StructByName("InterfaceUser")
			Expect(ok).To(BeTrue())

			Expect(s.Implements(pkg.Interfaces[2])).To(BeFalse())
		})
	})
})
