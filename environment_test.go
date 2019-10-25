package myasthurts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	Describe("ParseDir", func() {
		It("it should fail parsing a non existing directory", func() {
			env, err := NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			pkg, err := env.ParseDir("./data/non-exiting-directory")
			Expect(err).To(HaveOccurred())
			Expect(pkg).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("cannot find package"))
		})

		It("should parse the package", func() {
			env, err := NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Name()).To(Equal("Home"))
			Expect(pkg.Structs[0].Fields).To(HaveLen(4))
			Expect(pkg.Structs[1].Name()).To(Equal("User"))
			Expect(pkg.Structs[1].Fields).To(HaveLen(3))
		})

		It("should explore a package already registered but not explored", func() {
			env, err := NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			buildPkg, err := env.ImportDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())

			unexploredPkg := NewPackage(buildPkg)
			env.AppendPackage(unexploredPkg)

			Expect(unexploredPkg.explored).To(BeFalse())

			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())

			Expect(pkg.explored).To(BeTrue())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Name()).To(Equal("Home"))
			Expect(pkg.Structs[0].Fields).To(HaveLen(4))
			Expect(pkg.Structs[1].Name()).To(Equal("User"))
			Expect(pkg.Structs[1].Fields).To(HaveLen(3))
		})

		It("should skip a package already explored", func() {
			env, err := NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			exploredPkg := &Package{
				Name:       "models",
				ImportPath: ".",
				explored:   true, // Faking an explored package.
			}

			env.AppendPackage(exploredPkg)

			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())

			Expect(pkg.explored).To(BeTrue())

			// If it is not explored by the ParseDir, it will be empty.
			Expect(pkg.Structs).To(BeEmpty())
		})
	})
})
