package myasthurts_test

import (
	"go/parser"
	"go/token"
	"strings"
	"time"

	"github.com/fatih/structtag"
	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("My AST Hurts", func() {

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
		Expect(s.Fields).NotTo(BeNil())

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

		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Tests Functions - models5.sample ----------
		// TODO
	})

	type User struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	It("should parse struct tags", func() {

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "data/models6.sample", nil, parser.AllErrors)
		Expect(err).ToNot(HaveOccurred())
		Expect(f).ToNot(BeNil())
		Expect(f.Decls).ToNot(BeNil())

		//ast.Print(fset, f)
		env := &myasthurts.Environment{}
		myasthurts.Parse(f, env)

		// ---------- Tests struct tags - models5.sample ----------
		Expect(env.Packages).To(HaveLen(1))
		Expect(env.Packages[0].Structs).To(HaveLen(2))

		Expect(env.Packages[0].Structs[0].Name).To(Equal("User"))
		Expect(env.Packages[0].Structs[1].Name).To(Equal("Test"))

		Expect(env.Packages[0].Structs[0].Fields).ToNot(BeNil())
		Expect(env.Packages[0].Structs[1].Fields).ToNot(BeNil())

		structUser := env.Packages[0].Structs[0].Fields
		structTest := env.Packages[0].Structs[1].Fields

		Expect(structUser).To(HaveLen(6))
		Expect(structTest).To(HaveLen(3))

		Expect(structUser[0].Name).To(Equal("ID"))
		Expect(structUser[1].Name).To(Equal("Name"))
		Expect(structUser[2].Name).To(Equal("Email"))
		Expect(structUser[3].Name).To(Equal("Password"))
		Expect(structUser[4].Name).To(Equal("CreatedAt"))
		Expect(structUser[5].Name).To(Equal("UpdatedAt"))

		Expect(structTest[0].Name).To(Equal("ID"))
		Expect(structTest[1].Name).To(Equal("Name"))
		Expect(structTest[2].Name).To(Equal("Email"))

		// ---------- Tests Tags tags ----------
		Expect(structUser[0].Tag.Raw).ToNot(BeNil())
		Expect(structUser[1].Tag.Raw).ToNot(BeNil())
		Expect(structUser[2].Tag.Raw).ToNot(BeNil())
		Expect(structUser[3].Tag.Raw).ToNot(BeNil())
		Expect(structUser[4].Tag.Raw).ToNot(BeNil())
		Expect(structUser[5].Tag.Raw).ToNot(BeNil())

		Expect(structUser[0].Tag.Raw).To(Equal("`json:\"id,uuidTest\"`"))
		Expect(structUser[1].Tag.Raw).To(Equal("`json:\"name\"`"))
		Expect(structUser[2].Tag.Raw).To(Equal("`json:\"email\"`"))
		Expect(structUser[3].Tag.Raw).To(Equal("`json:\"password,old,newTest,moreField\"`"))
		Expect(structUser[4].Tag.Raw).To(Equal("`json:\"created_at\"`"))
		Expect(structUser[5].Tag.Raw).To(Equal("`json:\"updated_at\"`"))

		// ----------- Test Tag with value ID
		structUserTagId, err := structtag.Parse(strings.ReplaceAll(string(structUser[0].Tag.Raw), "`", ""))
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagId.Tags()).To(HaveLen(1))

		getStructUserTagId, err := structUserTagId.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagId.Key).To(Equal("json"))
		Expect(getStructUserTagId.Name).To(Equal("id"))
		Expect(getStructUserTagId.Options).To(HaveLen(1))
		Expect(getStructUserTagId.Options[0]).To(Equal("uuidTest"))

		// ----------- Test Tag with value Name
		structUserTagName, err := structtag.Parse(strings.ReplaceAll(string(structUser[1].Tag.Raw), "`", ""))
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagName.Tags()).To(HaveLen(1))

		getStructUserTagName, err := structUserTagName.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagName.Key).To(Equal("json"))
		Expect(getStructUserTagName.Name).To(Equal("name"))
		Expect(getStructUserTagName.Options).To(HaveLen(0))

		// ----------- Test Tag with value Email
		structUserTagEmail, err := structtag.Parse(strings.ReplaceAll(string(structUser[2].Tag.Raw), "`", ""))
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagEmail.Tags()).To(HaveLen(1))

		getStructUserTagEmail, err := structUserTagEmail.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagEmail.Key).To(Equal("json"))
		Expect(getStructUserTagEmail.Name).To(Equal("email"))
		Expect(getStructUserTagEmail.Options).To(HaveLen(0))

		// ----------- Test Tag with value Password
		structUserTagPassword, err := structtag.Parse(strings.ReplaceAll(string(structUser[3].Tag.Raw), "`", ""))
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagPassword.Tags()).To(HaveLen(1))

		getStructUserTagPassword, err := structUserTagPassword.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagPassword.Key).To(Equal("json"))
		Expect(getStructUserTagPassword.Name).To(Equal("password"))
		Expect(getStructUserTagPassword.Options).To(HaveLen(3))

		Expect(getStructUserTagPassword.Options[0]).To(Equal("old"))
		Expect(getStructUserTagPassword.Options[1]).To(Equal("newTest"))
		Expect(getStructUserTagPassword.Options[2]).To(Equal("moreField"))

		// ----------- Test Tag with value created_at
		structUserTagCreated_at, err := structtag.Parse(strings.ReplaceAll(string(structUser[4].Tag.Raw), "`", ""))
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagCreated_at.Tags()).To(HaveLen(1))

		getStructUserTagCreated_at, err := structUserTagCreated_at.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagCreated_at.Key).To(Equal("json"))
		Expect(getStructUserTagCreated_at.Name).To(Equal("created_at"))
		Expect(getStructUserTagCreated_at.Options).To(HaveLen(0))

		// ----------- Test Tag with value updated_at
		structUserTagUpdated_at, err := structtag.Parse(strings.ReplaceAll(string(structUser[5].Tag.Raw), "`", ""))
		Expect(err).ToNot(HaveOccurred())
		Expect(structUserTagUpdated_at.Tags()).To(HaveLen(1))

		getStructUserTagUpdated_at, err := structUserTagUpdated_at.Get("json")
		Expect(err).ToNot(HaveOccurred())
		Expect(getStructUserTagUpdated_at.Key).To(Equal("json"))
		Expect(getStructUserTagUpdated_at.Name).To(Equal("updated_at"))
		Expect(getStructUserTagUpdated_at.Options).To(HaveLen(0))
	})

})
