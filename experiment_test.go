package myasthurts_test

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/fatih/structtag"
	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts", func() {

	FContext("should parse struct", func() {

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

			Expect(pkg.Structs[0].Fields).To(HaveLen(3))
			Expect(pkg.Structs[1].Fields).To(HaveLen(2))

			Expect(pkg.Structs[0].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[0].Fields[1].Name).To(Equal("Name"))
			Expect(pkg.Structs[0].Fields[2].Name).To(Equal("Email"))

			Expect(pkg.Structs[0].Fields[0].Type).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[0].Fields[1].Type).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[0].Fields[2].Type).ToNot(Equal(BeNil()))

			Expect(pkg.Structs[0].Fields[0].Type.Name).To(Equal("int64"))
			Expect(pkg.Structs[0].Fields[1].Type.Name).To(Equal("string"))
			Expect(pkg.Structs[0].Fields[2].Type.Name).To(Equal("string"))

			Expect(pkg.Structs[0].Fields[0].Type.Type).To(BeNil())
			Expect(pkg.Structs[0].Fields[1].Type.Type).To(BeNil())
			Expect(pkg.Structs[0].Fields[2].Type.Type).To(BeNil())

			Expect(pkg.Structs[1].Fields[0].Name).To(Equal("ID"))
			Expect(pkg.Structs[1].Fields[1].Name).To(Equal("Address"))

			Expect(pkg.Structs[1].Fields[0].Type).ToNot(Equal(BeNil()))
			Expect(pkg.Structs[1].Fields[1].Type).ToNot(Equal(BeNil()))

			Expect(pkg.Structs[1].Fields[0].Type.Name).To(Equal("int64"))
			Expect(pkg.Structs[1].Fields[1].Type.Name).To(Equal("string"))

			Expect(pkg.Structs[1].Fields[0].Type.Type).To(BeNil())
			Expect(pkg.Structs[1].Fields[1].Type.Type).To(BeNil())

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

	It("should parse fields of User struct", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models2.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		// ---------- Test User struct - models2.sample ----------
		pkg, ok := env.PackageByName("models")
		fmt.Println(len(pkg.Structs))
		Expect(ok).To(BeTrue())
		Expect(pkg.Structs).To(HaveLen(1))

		s := pkg.Structs[0]

		Expect(s.Name()).To(Equal("User"))
		Expect(s.Fields).To(HaveLen(6))
		Expect(s.Fields).NotTo(BeNil())

		fields := pkg.Structs[0].Fields

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

		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		pkg, ok := env.PackageByName("models")

		// ---------- Test Function - models3.sample ----------
		Expect(ok).To(BeTrue())
		Expect(pkg.Name).To(Equal("models"))
		Expect(pkg.Methods).To(HaveLen(1))
		Expect(pkg.Methods[0].Name()).To(Equal("test_1"))
		Expect(pkg.Methods[0].Arguments).To(HaveLen(2))

		Expect(pkg.Methods[0].Arguments[0].Name).To(Equal("num"))
		Expect(pkg.Methods[0].Arguments[0].Type.Name).To(Equal("int"))

		Expect(pkg.Methods[0].Arguments[1].Name).To(Equal("text"))
		Expect(pkg.Methods[0].Arguments[1].Type.Name).To(Equal("string"))
	})

	It("should parse function in Struct", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models4.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		pkg, ok := env.PackageByName("models")

		// ---------- Tests Functions - models4.sample ----------

		Expect(ok).To(BeTrue())
		Expect(pkg.Structs).To(HaveLen(1))
		Expect(pkg.Structs[0].Methods).To(HaveLen(2))
		Expect(pkg.Structs[0].Methods[0].Descriptor.Name()).To(Equal("Address"))
		Expect(pkg.Structs[0].Methods[1].Descriptor.Name()).To(Equal("ChangePassword"))

		Expect(pkg.Methods).To(HaveLen(2))

		// Check if exist a func Address
		Expect(pkg.Methods[0].Name()).To(Equal("Address"))

		// Check if doesn't exist arguments from func Address
		Expect(pkg.Methods[0].Arguments).To(BeEmpty())

		// Check if exist all receptors from func Address
		Expect(pkg.Methods[0].Recv).To(HaveLen(1))
		Expect(pkg.Methods[0].Recv[0].Name).To(Equal("u"))
		Expect(pkg.Methods[0].Recv[0].Type.Name).To(Equal("User"))

		// ----- func ChangePassword -----

		// Check if exist a func ChangePassword
		Expect(pkg.Methods[1].Name()).To(Equal("ChangePassword"))

		// Check if exist all arguments from func ChangePassword
		Expect(pkg.Methods[1].Arguments).To(HaveLen(1))

		Expect(pkg.Methods[1].Arguments[0].Name).To(Equal("new"))
		Expect(pkg.Methods[1].Arguments[0].Type.Name).To(Equal("string"))

		// Check if exist all receptors from func ChangePassword
		Expect(pkg.Methods[1].Recv).To(HaveLen(1))

		Expect(pkg.Methods[1].Recv[0].Name).To(Equal("p"))
		Expect(pkg.Methods[1].Recv[0].Type.Name).To(Equal("User"))

	})

	PIt("should parse the variables", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models5.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		// ---------- Tests Functions - models5.sample ----------
		// TODO
	})

	/*Context("Struct Tags", func () {
		//
	})*/

	It("should parse struct tags", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models6.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		pkg, ok := env.PackageByName("models")

		// ---------- Tests struct tags - models6.sample ----------
		Expect(ok).To(BeTrue())
		Expect(pkg.Structs).To(HaveLen(2))

		Expect(pkg.Structs[0].Name()).To(Equal("User"))
		Expect(pkg.Structs[1].Name()).To(Equal("Test"))

		Expect(pkg.Structs[0].Fields).ToNot(BeNil())
		Expect(pkg.Structs[1].Fields).ToNot(BeNil())

		structUserFields := pkg.Structs[0].Fields
		structTestFields := pkg.Structs[1].Fields

		Expect(structUserFields).To(HaveLen(6))
		Expect(structTestFields).To(HaveLen(3))

		Expect(structUserFields[0].Name).To(Equal("ID"))
		Expect(structUserFields[1].Name).To(Equal("Name"))
		Expect(structUserFields[2].Name).To(Equal("Email"))
		Expect(structUserFields[3].Name).To(Equal("Password"))
		Expect(structUserFields[4].Name).To(Equal("CreatedAt"))
		Expect(structUserFields[5].Name).To(Equal("UpdatedAt"))

		Expect(structTestFields[0].Name).To(Equal("ID"))
		Expect(structTestFields[1].Name).To(Equal("Name"))
		Expect(structTestFields[2].Name).To(Equal("Email"))

		// ---------- Tests struct tags from struct User ----------
		Expect(structUserFields[0].Tag.Raw).ToNot(BeNil())
		Expect(structUserFields[1].Tag.Raw).ToNot(BeNil())
		Expect(structUserFields[2].Tag.Raw).ToNot(BeNil())
		Expect(structUserFields[3].Tag.Raw).ToNot(BeNil())
		Expect(structUserFields[4].Tag.Raw).ToNot(BeNil())
		Expect(structUserFields[5].Tag.Raw).To(BeEmpty())

		Expect(structUserFields[0].Tag.Raw).To(Equal("json:\"id,uuidTest\""))
		Expect(structUserFields[1].Tag.Raw).To(Equal("json:\"name\" bson:\"\""))
		Expect(structUserFields[2].Tag.Raw).To(Equal("json:\"email\""))
		Expect(structUserFields[3].Tag.Raw).To(Equal("json:\"password,old,newTest,moreField\""))
		Expect(structUserFields[4].Tag.Raw).To(Equal("json:\"created_at\""))

		// ----------- Test Tag with value ID
		structUserTagId, err := structtag.Parse(structUserFields[0].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagId.Tags()).To(HaveLen(1))

		getStructUserTagId, err := structUserTagId.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagId.Key).To(Equal("json"))
		Expect(getStructUserTagId.Name).To(Equal("id"))
		Expect(getStructUserTagId.Options).To(HaveLen(1))
		Expect(getStructUserTagId.Options[0]).To(Equal("uuidTest"))

		// ----------- Test Tag with value Name
		structUserTagName, err := structtag.Parse(structUserFields[1].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagName.Tags()).To(HaveLen(2))

		getStructUserTagName, err := structUserTagName.Get("json")
		Expect(err).ToNot(HaveOccurred())
		getStructUserTagNameBson, err := structUserTagName.Get("bson")
		Expect(err).ToNot(HaveOccurred())

		Expect(getStructUserTagName.Key).To(Equal("json"))
		Expect(getStructUserTagName.Name).To(Equal("name"))
		Expect(getStructUserTagName.Options).To(BeEmpty())

		Expect(getStructUserTagNameBson.Key).To(Equal("bson"))
		Expect(getStructUserTagNameBson.Name).To(BeEmpty())
		Expect(getStructUserTagNameBson.Options).To(BeEmpty())

		// ----------- Test Tag with value Email
		structUserTagEmail, err := structtag.Parse(structUserFields[2].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagEmail.Tags()).To(HaveLen(1))

		getStructUserTagEmail, err := structUserTagEmail.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagEmail.Key).To(Equal("json"))
		Expect(getStructUserTagEmail.Name).To(Equal("email"))
		Expect(getStructUserTagEmail.Options).To(BeEmpty())

		// ----------- Test Tag with value Password
		structUserTagPassword, err := structtag.Parse(structUserFields[3].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagPassword.Tags()).To(HaveLen(1))
		Expect(structUserFields[3].Tag.Params).To(HaveLen(1))

		Expect(structUserFields[3].Tag.Params[0].Options).To(HaveLen(3))

		Expect(structUserFields[3].Tag.Params[0].Options[0]).To(Equal("old"))
		Expect(structUserFields[3].Tag.Params[0].Options[1]).To(Equal("newTest"))
		Expect(structUserFields[3].Tag.Params[0].Options[2]).To(Equal("moreField"))

		getStructUserTagPassword, err := structUserTagPassword.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagPassword.Key).To(Equal("json"))
		Expect(getStructUserTagPassword.Name).To(Equal("password"))
		Expect(getStructUserTagPassword.Options).To(HaveLen(3))

		Expect(structUserFields[3].Tag.Params[0].Options[0]).To(Equal(getStructUserTagPassword.Options[0]))
		Expect(structUserFields[3].Tag.Params[0].Options[1]).To(Equal(getStructUserTagPassword.Options[1]))
		Expect(structUserFields[3].Tag.Params[0].Options[2]).To(Equal(getStructUserTagPassword.Options[2]))

		// ----------- Test Tag with value created_at
		structUserTagCreated_at, err := structtag.Parse(structUserFields[4].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagCreated_at.Tags()).To(HaveLen(1))

		getStructUserTagCreated_at, err := structUserTagCreated_at.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagCreated_at.Key).To(Equal("json"))
		Expect(getStructUserTagCreated_at.Name).To(Equal("created_at"))
		Expect(getStructUserTagCreated_at.Options).To(BeEmpty())

		// ----------- Test Tag with value updated_at
		structUserTagUpdated_at, err := structtag.Parse(structUserFields[5].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagUpdated_at.Tags()).To(BeEmpty())

		// ---------- Test Struct tag from struct Test ----------

		Expect(structTestFields[0].Tag.Raw).To(Equal("json:\"id\""))
		Expect(structTestFields[1].Tag.Raw).To(Equal("json:\"name\""))
		Expect(structTestFields[2].Tag.Raw).To(Equal("json:\"email\""))

		// ----------- Test Tag with value ID
		structTestTagId, err := structtag.Parse(structTestFields[0].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structTestTagId.Tags()).To(HaveLen(1))

		getStructTestTagId, err := structTestTagId.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructTestTagId.Key).To(Equal("json"))
		Expect(getStructTestTagId.Name).To(Equal("id"))
		Expect(getStructTestTagId.Options).To(BeEmpty())

		// ----------- Test Tag with value Name
		structTestTagName, err := structtag.Parse(structTestFields[1].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structTestTagName.Tags()).To(HaveLen(1))

		getStructTestTagName, err := structTestTagName.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructTestTagName.Key).To(Equal("json"))
		Expect(getStructTestTagName.Name).To(Equal("name"))
		Expect(getStructTestTagName.Options).To(BeEmpty())

		// ----------- Test Tag with value Email
		structTestTagEmail, err := structtag.Parse(structTestFields[2].Tag.Raw)
		Expect(err).ToNot(HaveOccurred())
		Expect(structTestTagEmail.Tags()).To(HaveLen(1))

		getStructTestTagEmail, err := structTestTagEmail.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructTestTagEmail.Key).To(Equal("json"))
		Expect(getStructTestTagEmail.Name).To(Equal("email"))
		Expect(getStructTestTagEmail.Options).To(BeEmpty())

	})

	It("should parse struct and func User", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models7.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		pkg, ok := env.PackageByName("models")

		// ---------- Tests struct tags - models7.sample ----------

		Expect(ok).To(BeTrue())
		Expect(pkg.Name).To(Equal("models"))

		Expect(pkg.RefType).To(HaveLen(4))
		Expect(pkg.RefType[0].Name).To(Equal("time"))
		Expect(pkg.RefType[1].Name).To(Equal("User"))
		Expect(pkg.RefType[2].Name).To(Equal("int64"))
		Expect(pkg.RefType[3].Name).To(Equal("string"))

		Expect(pkg.RefType[0].Pkg.Name).To(Equal("models"))
		Expect(pkg.RefType[1].Pkg.Name).To(Equal("models"))
		Expect(pkg.RefType[2].Pkg.Name).To(Equal("models"))
		Expect(pkg.RefType[3].Pkg.Name).To(Equal("models"))

		Expect(pkg.RefType[0].Type).To(BeNil())
		Expect(pkg.RefType[0].Name).To(Equal("time"))

		Expect(pkg.RefType[1].Type).ToNot(BeNil())
		Expect(pkg.RefType[1].Type.Name()).To(Equal("User"))

		Expect(pkg.RefType[2].Type).To(BeNil())
		Expect(pkg.RefType[2].Name).To(Equal("int64"))

		Expect(pkg.RefType[3].Type).To(BeNil())
		Expect(pkg.RefType[3].Name).To(Equal("string"))

		Expect(pkg.Structs).To(HaveLen(1))
		Expect(pkg.Structs[0].Name()).To(Equal("User"))
		Expect(pkg.Structs[0].Methods).To(HaveLen(1))
		Expect(pkg.Structs[0].Methods[0].Descriptor.Name()).To(Equal("getName"))

		Expect(pkg.Methods).To(HaveLen(1))
		Expect(pkg.Methods[0].Name()).To(Equal("getName"))
		Expect(pkg.Methods[0].Recv).To(HaveLen(1))
		Expect(pkg.Methods[0].Recv[0].Type.Name).To(Equal(pkg.Structs[0].Name()))

	})

	It("should parse imports", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models8.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)
		// ---------- Tests struct tags - models8.sample ----------

		pkg, ok := env.PackageByName("models")

		Expect(ok).To(BeTrue())
		Expect(pkg.Name).To(Equal("models"))

		Expect(pkg.Structs).To(HaveLen(1))
		Expect(pkg.Structs[0].Methods).To(HaveLen(1))
		Expect(pkg.Structs[0].Name()).To(Equal("User"))
		Expect(pkg.Structs[0].Fields).To(HaveLen(3))

		Expect(pkg.Structs[0].Fields[0].Name).To(Equal("ID"))
		Expect(pkg.Structs[0].Fields[1].Name).To(Equal("Name"))
		Expect(pkg.Structs[0].Fields[2].Name).To(Equal("CreatedAt"))

		Expect(pkg.Structs[0].Fields[0].Type.Name).To(Equal("int64"))
		Expect(pkg.Structs[0].Fields[1].Type.Name).To(Equal("string"))

		time, okT := env.PackageByName("time")
		Expect(okT).To(BeTrue())
		Expect(pkg.Structs[0].Fields[2].Type.Name).To(Equal("t"))
		Expect(time.Name).To(Equal("time"))

		Expect(pkg.Methods).To(HaveLen(1))
		Expect(pkg.Methods[0].Name()).To(Equal("getName"))

	})

	It("should parse types in struct", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models9.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)
		// ---------- Tests struct tags - models9.sample ----------

		pkg, ok := env.PackageByName("models")

		Expect(ok).To(BeTrue())
		Expect(pkg.Structs).To(HaveLen(2))

		Expect(pkg.Structs[0].Name()).To(Equal("Compra"))
		Expect(pkg.Structs[1].Name()).To(Equal("User"))

		Expect(pkg.Structs[0].Fields).To(HaveLen(2))
		Expect(pkg.Structs[0].Fields[1].Type).ToNot(BeNil())

		refType, _ := pkg.RefTypeByName("string")
		Expect(refType).ToNot(BeNil())

		Expect(pkg.Structs[0].Fields[1].Type.Pkg).To(Equal(refType.Pkg))

	})

	It("should parse multi line comments in struct", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models11.sample", nil, parser.ParseComments)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		pkg, ok := env.PackageByName("models")
		Expect(ok).To(BeTrue())

		Expect(pkg.Structs).To(HaveLen(1))
		Expect(pkg.Structs[0].Comment).To(HaveLen(2))

	})

	It("should parse multi line comments in struct field", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models12.sample", nil, parser.ParseComments)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)

		pkg, ok := env.PackageByName("models")
		Expect(ok).To(BeTrue())

		Expect(pkg.Structs).To(HaveLen(1))
		Expect(pkg.Structs[0].Comment).To(BeEmpty())
		Expect(pkg.Structs[0].Fields).To(HaveLen(2))
		Expect(pkg.Structs[0].Fields[1].Comment).To(HaveLen(2))

	})

	/*It("should parse comments in struct", func() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models10.sample", nil, parser.ParseComments)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := myasthurts.NewEnvironment()
		myasthurts.Parse(env, f)
		// ---------- Tests struct tags - models10.sample ----------

		pkg, ok := env.PackageByName("models")
		Expect(ok).To(BeTrue())

		Expect(pkg.Comment).To(HaveLen(1))

		Expect(pkg.Structs).To(HaveLen(2))

		Expect(pkg.Structs[0].Comment).To(HaveLen(2))
		Expect(pkg.Structs[1].Comment).To(HaveLen(1))

		Expect(pkg.Structs[0].Comment[0]).To(Equal("// Struct type Compra"))
		Expect(pkg.Structs[0].Comment[1]).To(Equal("// Test New line test"))

		Expect(pkg.Structs[1].Comment[0]).To(Equal("/** Comment for test\ntext for more test\ntesting comments "))

		Expect(pkg.Structs[1].Fields).To(HaveLen(3))
		Expect(pkg.Structs[1].Fields[0].Comment).To(HaveLen(1))

	})*/

})
