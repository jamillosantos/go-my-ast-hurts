package myasthurts_test

import (
	//myasthurts "github.com/jamillosantos/go-my-ast-hurts"

	"fmt"
	"go/build"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
)

func newDataPackageContext(env *myasthurts.Environment) *myasthurts.ParsePackageContext {
	buildPkg := &build.Package{
		Dir:        "data",
		Name:       "models",
		ImportPath: "data",
	}
	pkg := myasthurts.NewPackage(buildPkg)
	env.AppendPackage(pkg)
	return myasthurts.NewPackageContext(pkg, buildPkg)
}

var _ = Describe("My AST Hurts - Parse simples files with tags and func from struct", func() {

	Context("Initialization", func() {

		var (
			originalGOROOT string
		)

		BeforeSuite(func() {
			originalGOROOT = build.Default.GOROOT
		})

		JustAfterEach(func() {
			build.Default.GOROOT = originalGOROOT
		})

		It("should show error if GOROOT not exist", func() {
			Expect(build.Default.GOROOT).ToNot(BeEmpty())

			build.Default.GOROOT = path.Join("not", "existing", "path")

			_, err := myasthurts.NewEnvironment()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("go/build"))
			Expect(err.Error()).To(ContainSubstring("cannot find GOROOT directory"))
		})

		It("should show error if file not found", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			err = env.ParseFile(newDataPackageContext(env), "data/models259.sample.go")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no such file or directory"))
		})

	})

	When("parsing struct", func() {
		It("should check two struct in file and if anyName struct no exist", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			err = env.ParseFile(newDataPackageContext(env), "data/models1.sample.go")
			Expect(err).ToNot(HaveOccurred())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg).NotTo(BeNil())
			Expect(pkg.Structs).To(HaveLen(2))

			Expect(pkg.StructByName("anyName")).To(BeNil())
		})

		It("should check struct fields", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			err = env.ParseFile(newDataPackageContext(env), "data/models2.sample.go")
			Expect(err).ToNot(HaveOccurred())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields).To(HaveLen(4))

			Expect(pkg.Structs[0].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[0].Fields[1].Name).To(Equal("Name"))

			Expect(pkg.Structs[0].Fields[0].RefType).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[0].Fields[1].RefType).ToNot(Equal(BeNil()))

			Expect(pkg.Structs[0].Fields[0].RefType.Name()).To(Equal("int64"))
			Expect(pkg.Structs[0].Fields[1].RefType.Name()).To(Equal("string"))

			Expect(pkg.Structs[0].Fields[0].RefType.Type()).To(BeNil())
			Expect(pkg.Structs[0].Fields[1].RefType.Type()).To(BeNil())

			Expect(pkg.Structs[1].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[1].Fields[1].Name).To(Equal("Address"))
			Expect(pkg.Structs[1].Fields[2].Name).To(Equal("User"))

			Expect(pkg.Structs[1].Fields[0].RefType).ToNot(BeNil())
			Expect(pkg.Structs[1].Fields[1].RefType).ToNot(BeNil())
			Expect(pkg.Structs[1].Fields[2].RefType).ToNot(BeNil())

			Expect(pkg.Structs[1].Fields[0].RefType.Name()).To(Equal("int64"))
			Expect(pkg.Structs[1].Fields[1].RefType.Name()).To(Equal("string"))
			Expect(pkg.Structs[1].Fields[2].RefType.Name()).To(Equal("User"))
			Expect(pkg.Structs[1].Fields[3].RefType.Name()).To(Equal("User"))

			Expect(pkg.Structs[1].Fields[0].RefType.Type()).To(BeNil())
			Expect(pkg.Structs[1].Fields[1].RefType.Type()).To(BeNil())
			Expect(pkg.Structs[1].Fields[2].RefType.Type()).To(Equal(pkg.Structs[0]))
			Expect(pkg.Structs[1].Fields[3].RefType.Type()).To(Equal(pkg.Structs[0]))
		})

		It("should check struct tags", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).To(BeNil())

			err = env.ParseFile(newDataPackageContext(env), "data/models3.sample.go")
			Expect(err).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
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
			Expect(pkg.Structs[0].Fields[0].Tag.TagParamByName("test")).To(BeNil())

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
			Expect(pkg.Structs[0].Fields[2].Tag.TagParamByName("bson")).ToNot(BeNil())

			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Name).To(Equal("bson"))
			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Value).To(BeEmpty())
			Expect(pkg.Structs[0].Fields[2].Tag.Params[1].Options).To(BeEmpty())

		})

		It("should check struct custom field with user", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models4.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields).To(HaveLen(2))
			Expect(pkg.Structs[0].Fields[1].RefType.Name()).To(Equal("User"))
			Expect(pkg.Structs[0].Fields[1].RefType.Type()).To(Equal(pkg.Structs[1]))
			Expect(pkg.Structs[0].Fields[1].RefType.Pkg()).To(Equal(pkg))

		})

		It("should check struct comments", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models5.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
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

		It("should parse a struct with a interface{} member", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models14.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())
			Expect(pkg).NotTo(BeNil())

			Expect(pkg.Structs).To(HaveLen(2))

			Expect(pkg.Structs[1].Fields).To(HaveLen(5))
			Expect(pkg.Structs[1].Fields[3].RefType.Name()).To(BeEmpty())
			Expect(pkg.Structs[1].Fields[3].RefType.Type()).ToNot(BeNil())
			var iType *myasthurts.Interface
			Expect(pkg.Structs[1].Fields[3].RefType.Type()).To(BeAssignableToTypeOf(iType))
			iType = pkg.Structs[1].Fields[3].RefType.Type().(*myasthurts.Interface)
			Expect(iType.Methods()).To(BeEmpty())
		})
	})

	When("parsing variables", func() {

		It("should check variables declarated with var", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())
			env.Config = myasthurts.EnvConfig{
				DevMode: true,
				ASTI:    false,
			}

			exrr = env.ParseFile(newDataPackageContext(env), "data/models11.sample.go")
			builtin, _ := env.PackageByImportPath("builtin")
			Expect(builtin).ToNot(BeNil())
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg.Variables).To(HaveLen(8))

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

			h := pkg.VariableByName("h")
			Expect(h).ToNot(BeNil())
			Expect(h.RefType).ToNot(BeNil())

			x := pkg.VariableByName("x")
			Expect(x).To(BeNil())

			ref, ok := pkg.RefTypeByName("User")
			Expect(ok).To(BeTrue())
			Expect(ref).To(Equal(g.RefType))

			Expect(g.RefType.Type).ToNot(BeNil())
			Expect(g.RefType.Type().Name()).To(Equal("User"))

			/* Obs: At the moment it is not possible to identify if
			 * 		the variable is array or no. (This is necessary?)
			 */
		})

		It("should check const", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())
			env.Config = myasthurts.EnvConfig{
				DevMode: true,
				ASTI:    false,
			}

			exrr = env.ParseFile(newDataPackageContext(env), "data/models15.sample.go")
			builtin, _ := env.PackageByImportPath("builtin")
			Expect(builtin).ToNot(BeNil())
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			a := pkg.VariableByName("PI")
			Expect(a).ToNot(BeNil())
			Expect(a.RefType).ToNot(BeNil())

			b := pkg.VariableByName("OLM")
			Expect(b).ToNot(BeNil())
			Expect(b.RefType).ToNot(BeNil())
		})
	})

	When("parsing function", func() {

		It("should check name and parameters from function", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models6.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())
			Expect(pkg.Name).To(Equal("models"))

			Expect(pkg.Methods).To(HaveLen(2))
			Expect(pkg.Methods[0].Name()).To(Equal("show"))
			Expect(pkg.Methods[0].Arguments).To(HaveLen(2))
			Expect(pkg.Methods[0].Recv).To(BeEmpty())

			Expect(pkg.Methods[0].Arguments[0].Name).To(Equal("name"))
			Expect(pkg.Methods[0].Arguments[1].Name).To(Equal("age"))

			Expect(pkg.Methods[0].Arguments[0].Type.Type()).To(BeNil())
			Expect(pkg.Methods[0].Arguments[1].Type.Type()).To(BeNil())

			Expect(pkg.Methods[0].Arguments[0].Type.Name()).To(Equal("string"))
			Expect(pkg.Methods[0].Arguments[1].Type.Name()).To(Equal("int64"))

			Expect(pkg.Methods[1].Name()).To(Equal("welcome"))
			Expect(pkg.Methods[1].Arguments).To(BeEmpty())
			Expect(pkg.Methods[1].Recv).To(BeEmpty())
			Expect(pkg.Methods[1].Result).To(HaveLen(1))
			Expect(pkg.Methods[1].Result[0].Type.Name()).To(Equal("string"))

			Expect(pkg.MethodsMap).To(HaveLen(2))
			Expect(pkg.MethodsMap).To(HaveKey("show"))
			Expect(pkg.MethodsMap).To(HaveKey("welcome"))

			Expect(pkg.Methods[0].Name()).To(Equal("show"))

		})

		It("should check if func belong to struct", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models7.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg.Methods).To(HaveLen(1))
			Expect(pkg.Methods[0].Name()).To(Equal("getName_"))

			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods()).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods()[0].Descriptor.Name()).To(Equal("getName"))
			Expect(pkg.Structs[0].Methods()[0].Descriptor.Arguments).To(BeEmpty())
			Expect(pkg.Structs[0].Methods()[0].Descriptor.Recv).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods()[0].Descriptor.Recv[0].Type.Type()).To(Equal(pkg.Structs[0]))
			Expect(pkg.Structs[0].Methods()[0].Descriptor.Result).To(HaveLen(1))
			Expect(pkg.Structs[0].Methods()[0].Descriptor.Result[0].Type.Name()).To(Equal("string"))
			Expect(pkg.Structs[0].MethodsMap()).To(HaveLen(1))
			Expect(pkg.Structs[0].MethodsMap()).To(HaveKey("getName"))
		})

		It("should check package of func", func() {
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models7.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg.Methods).To(HaveLen(1))
			Expect(pkg.Methods[0].Package()).To(Equal(pkg))
		})
	})

	When("parsing imports", func() {

		It("should check names all imports", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models8.sample.go")
			Expect(exrr).To(BeNil())

			models, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())
			fmt, ok := env.PackageByImportPath("fmt")
			Expect(ok).To(BeTrue())
			time, ok := env.PackageByImportPath("time")
			Expect(ok).To(BeTrue())

			Expect(models).ToNot(BeNil())
			Expect(fmt).ToNot(BeNil())
			Expect(time).ToNot(BeNil())

			Expect(models.Name).To(Equal("models"))
			Expect(fmt.Name).To(Equal("fmt"))
			Expect(time.Name).To(Equal("time"))
		})

		It("should explore a dot imported package", func() {
			env, err := myasthurts.NewEnvironment()
			Expect(err).ToNot(HaveOccurred())

			dataPkgCtx := newDataPackageContext(env)
			Expect(env.ParseFile(dataPkgCtx, "data/models12.sample.go")).To(Succeed())

			models, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			bytes, ok := env.PackageByImportPath("bytes")
			Expect(ok).To(BeTrue())

			Expect(models.Methods).To(HaveLen(2))
			Expect(models.Methods[0].Name()).To(Equal("welcome"))
			Expect(models.Methods[0].Arguments).To(HaveLen(1))
			Expect(models.Methods[0].Arguments[0].Name).To(Equal("buf"))

			ref, ok := bytes.RefTypeByName("Buffer")
			Expect(ok).To(BeTrue())
			Expect(ref).ToNot(BeNil())
			Expect(models.Methods[0].Arguments[0].Type.Name()).To(Equal(ref.Type().Name()))

			stct, ok := bytes.StructByName("Buffer")
			Expect(ok).To(BeTrue())
			Expect(stct).ToNot(BeNil())
			Expect(fmt.Sprintf("%p", models.Methods[0].Arguments[0].Type.Type())).To(Equal(fmt.Sprintf("%p", stct)))
		})
	})

	When("parsing comments", func() {

		It("should check multilines or no in struct comments", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models9.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
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

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models14.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg.Doc.Comments).To(HaveLen(6))

			Expect(pkg.Structs).To(HaveLen(2))
			Expect(pkg.Structs[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Structs[0].Doc.FormatComment()).To(Equal("Testing comment\nnew line"))

			Expect(pkg.Structs[1].Fields).To(HaveLen(5))
			Expect(pkg.Structs[1].Fields[1].Doc.Comments).To(HaveLen(2))
			Expect(pkg.Structs[1].Fields[1].Doc.FormatComment()).To(Equal("Line 1\nLine 2\n"))

		})

		It("should check comments in func", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models10.sample.go")
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("data")
			Expect(ok).To(BeTrue())

			Expect(pkg.Doc.Comments).To(HaveLen(7))
			Expect(pkg.Doc.Comments[0]).To(Equal("// Package models is a test"))

			Expect(pkg.Methods).To(HaveLen(2))
			Expect(pkg.Methods[0].Doc.Comments).To(HaveLen(1))
			Expect(pkg.Methods[0].Doc.Comments[0]).To(Equal("/** Description\n  multilines\n*/"))
			Expect(pkg.Methods[1].Doc.Comments).To(BeEmpty())
		})
	})

	When("parsing builtin file", func() {

		It("should check types from builtin file", func() {

			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			exrr = env.ParseFile(newDataPackageContext(env), "data/models13.sample.go")
			Expect(exrr).To(BeNil())

			pkgM, okM := env.PackageByImportPath("data")
			Expect(okM).To(BeTrue())
			pkgB, okB := env.PackageByImportPath("builtin")
			Expect(okB).To(BeTrue())

			T1, ok := pkgB.RefTypeByName("string")
			Expect(ok).To(BeTrue())
			Expect(T1).ToNot(BeNil())

			T2, ok := pkgB.RefTypeByName("int")
			Expect(ok).To(BeTrue())
			Expect(T2).ToNot(BeNil())

			T3, ok := pkgB.RefTypeByName("int8")
			Expect(ok).To(BeTrue())
			Expect(T3).ToNot(BeNil())

			T4, ok := pkgB.RefTypeByName("int16")
			Expect(ok).To(BeTrue())
			Expect(T4).ToNot(BeNil())

			T5, ok := pkgB.RefTypeByName("int32")
			Expect(ok).To(BeTrue())
			Expect(T5).ToNot(BeNil())

			T6, ok := pkgB.RefTypeByName("int64")
			Expect(ok).To(BeTrue())
			Expect(T6).ToNot(BeNil())

			T7, ok := pkgB.RefTypeByName("float32")
			Expect(ok).To(BeTrue())
			Expect(T7).ToNot(BeNil())

			T8, ok := pkgB.RefTypeByName("float64")
			Expect(ok).To(BeTrue())
			Expect(T8).ToNot(BeNil())

			T9, ok := pkgB.RefTypeByName("byte")
			Expect(ok).To(BeTrue())
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
			env, exrr := myasthurts.NewEnvironment()
			Expect(exrr).To(BeNil())

			pkg, ok := env.PackageByImportPath("builtin")
			Expect(ok).To(BeTrue())

			ref, ok := pkg.RefTypeByName("string")
			Expect(ok).To(BeTrue())
			Expect(ref).ToNot(BeNil())

			ref, ok = pkg.RefTypeByName("int64")
			Expect(ok).To(BeTrue())
			Expect(ref).ToNot(BeNil())

			ref, ok = pkg.RefTypeByName("float32")
			Expect(ok).To(BeTrue())
			Expect(ref).ToNot(BeNil())

			// ### WORK IN PROGRESS ###
		})

	})
})
