package myasthurts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interface", func() {
	Describe("Parse", func() {
		It("should parse interfaces", func() {
			env, err := NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			Expect(env.parseFile(newDataPackageContext(env), "data/interface.go")).To(Succeed())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg.Interfaces).To(HaveLen(3))
			Expect(pkg.Interfaces[0].Name()).To(Equal("HasName"))
			Expect(pkg.Interfaces[0].Methods()).To(HaveLen(2))
			Expect(pkg.Interfaces[0].Methods()[0].Descriptor.Name()).To(Equal("Name"))
			Expect(pkg.Interfaces[0].Methods()[0].Descriptor.Arguments).To(BeEmpty())
			Expect(pkg.Interfaces[0].Methods()[1].Descriptor.Name()).To(Equal("SetName"))
			Expect(pkg.Interfaces[0].Methods()[1].Descriptor.Arguments).To(HaveLen(1))
			Expect(pkg.Interfaces[0].Methods()[1].Descriptor.Arguments[0].Name).To(Equal("value"))
			Expect(pkg.Interfaces[0].Methods()[1].Descriptor.Arguments[0].Type.Name()).To(Equal("string"))
			Expect(pkg.Interfaces[0].methodsMap).To(HaveLen(2))
			Expect(pkg.Interfaces[0].methodsMap).To(HaveKey("Name"))
			Expect(pkg.Interfaces[0].methodsMap).To(HaveKey("SetName"))
			Expect(pkg.Interfaces[0].Methods()[1].Descriptor.Arguments[0].Type.Name()).To(Equal("string"))
			Expect(pkg.Interfaces[1].Name()).To(Equal("HasAge"))
			Expect(pkg.Interfaces[1].Methods()).To(HaveLen(2))
			Expect(pkg.Interfaces[1].Methods()[0].Descriptor.Name()).To(Equal("Age"))
			Expect(pkg.Interfaces[1].Methods()[0].Descriptor.Arguments).To(BeEmpty())
			Expect(pkg.Interfaces[1].Methods()[1].Descriptor.Name()).To(Equal("SetAge"))
			Expect(pkg.Interfaces[1].Methods()[1].Descriptor.Arguments).To(HaveLen(1))
			Expect(pkg.Interfaces[1].Methods()[1].Descriptor.Arguments[0].Name).To(Equal("value"))
			Expect(pkg.Interfaces[1].Methods()[1].Descriptor.Arguments[0].Type.Name()).To(Equal("int"))
			Expect(pkg.Interfaces[1].methodsMap).To(HaveLen(2))
			Expect(pkg.Interfaces[1].methodsMap).To(HaveKey("Age"))
			Expect(pkg.Interfaces[1].methodsMap).To(HaveKey("SetAge"))
			Expect(pkg.Interfaces[2].Name()).To(Equal("HasNameWrong"))
			Expect(pkg.Interfaces[2].Methods()).To(HaveLen(2))
			Expect(pkg.Interfaces[2].Methods()[0].Descriptor.Name()).To(Equal("Name"))
			Expect(pkg.Interfaces[2].Methods()[0].Descriptor.Arguments).To(BeEmpty())
			Expect(pkg.Interfaces[2].Methods()[1].Descriptor.Name()).To(Equal("SetName"))
			Expect(pkg.Interfaces[2].Methods()[1].Descriptor.Arguments).To(HaveLen(1))
			Expect(pkg.Interfaces[2].Methods()[1].Descriptor.Arguments[0].Name).To(Equal("value"))
			Expect(pkg.Interfaces[2].Methods()[1].Descriptor.Arguments[0].Type.Name()).To(Equal("int"))
			Expect(pkg.Interfaces[2].methodsMap).To(HaveLen(2))
			Expect(pkg.Interfaces[2].methodsMap).To(HaveKey("Name"))
			Expect(pkg.Interfaces[2].methodsMap).To(HaveKey("SetName"))
		})
	})
})
