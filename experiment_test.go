package myasthurts_test

import (
	"go/ast"
	"go/parser"
	"go/token"

	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts", func() {

	It("should parse a User struct", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models2.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Test User struct - models2.sample ----------
		Expect(env.Packages).To(HaveLen(1))
		Expect(env.Packages[0].Structs).To(HaveLen(1))

		s := env.Packages[0].Structs[0]

		Expect(s.Name).To(Equal("User"))
		Expect(s.Fields).To(HaveLen(6))
		//Expect(s.Comment).To(Equal("User is a model.")) TODO
	})

	It("should parse fields of User struct", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models2.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Test User struct - models2.sample ----------
		Expect(env.Packages).To(HaveLen(1))
		Expect(env.Packages[0].Structs).To(HaveLen(1))

		s := env.Packages[0].Structs[0]

		Expect(s.Name).To(Equal("User"))
		Expect(s.Fields).To(HaveLen(6))

		fields := env.Packages[0].Structs[0].Fields

		Expect(fields[0].Name).To(Equal("ID"))
		Expect(fields[1].Name).To(Equal("Name"))
		Expect(fields[2].Name).To(Equal("Email"))
		Expect(fields[3].Name).To(Equal("Password"))
		Expect(fields[4].Name).To(Equal("CreatedAt"))
		Expect(fields[5].Name).To(Equal("UpdatedAt"))

	})

	It("should parse function", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models3.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Test Function - models3.sample ----------
		Expect(env.Packages).To(HaveLen(1))
		Expect(env.Packages[0].Methods).To(HaveLen(1))
		Expect(env.Packages[0].Methods[0].Name).To(Equal("test_1"))
		Expect(env.Packages[0].Methods[0].Arguments).To(HaveLen(2))

		Expect(env.Packages[0].Methods[0].Arguments[0].Name).To(Equal("num"))
		Expect(env.Packages[0].Methods[0].Arguments[0].Type).To(Equal("int"))

		Expect(env.Packages[0].Methods[0].Arguments[1].Name).To(Equal("text"))
		Expect(env.Packages[0].Methods[0].Arguments[1].Type).To(Equal("string"))
	})

	It("should parse function in Struct", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models4.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Tests Functions - models4.sample ----------

		Expect(env.Packages).To(HaveLen(1))
		Expect(env.Packages[0].Structs).To(HaveLen(1))
		Expect(env.Packages[0].Structs[0].Methods).To(HaveLen(2))
		Expect(env.Packages[0].Structs[0].Methods[0].Descriptor.Name).To(Equal("Address"))
		Expect(env.Packages[0].Structs[0].Methods[1].Descriptor.Name).To(Equal("ChangePassword"))

		Expect(env.Packages[0].Methods).To(HaveLen(2))

		// Check if exist a func Address
		Expect(env.Packages[0].Methods[0].Name).To(Equal("Address"))

		// Check if doesn't exist arguments from func Address
		Expect(env.Packages[0].Methods[0].Arguments).To(HaveLen(0))

		// Check if exist all receptors from func Address
		Expect(env.Packages[0].Methods[0].Recv).To(HaveLen(1))
		Expect(env.Packages[0].Methods[0].Recv[0].Name).To(Equal("u"))
		Expect(env.Packages[0].Methods[0].Recv[0].Type).To(Equal("User"))

		// ----- func ChangePassword -----

		// Check if exist a func ChangePassword
		Expect(env.Packages[0].Methods[1].Name).To(Equal("ChangePassword"))

		// Check if exist all arguments from func ChangePassword
		Expect(env.Packages[0].Methods[1].Arguments).To(HaveLen(1))

		Expect(env.Packages[0].Methods[1].Arguments[0].Name).To(Equal("new"))
		Expect(env.Packages[0].Methods[1].Arguments[0].Type).To(Equal("string"))

		// Check if exist all receptors from func ChangePassword
		Expect(env.Packages[0].Methods[1].Recv).To(HaveLen(1))

		Expect(env.Packages[0].Methods[1].Recv[0].Name).To(Equal("p"))
		Expect(env.Packages[0].Methods[1].Recv[0].Type).To(Equal("User"))

	})

	It("should parse the variables", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models5.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		ast.Print(fset, f)

		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Tests Functions - models5.sample ----------

	})

})
