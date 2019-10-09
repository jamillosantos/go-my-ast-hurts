package myasthurts

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/fatih/structtag"
)

func (env *environment) makeEnv() (exrr error) {

	path := ""
	if path, exrr = env.basePath(); exrr != nil {
		return exrr
	}

	builtinFile := fmt.Sprintf("%s/builtin/builtin.go", path)

	env.parse(builtinFile)
	return
}

func (p *parseContext) PackageByName(name string) (*Package, bool) {
	pkg, ok := p.PackagesMap[name]
	return pkg, ok
}

func (env *environment) parse(fileLocation string) (exrr error) {

	var (
		file *ast.File
		fset *token.FileSet
	)

	fset = token.NewFileSet()
	if file, exrr = parser.ParseFile(fset, fileLocation, nil, parser.ParseComments); exrr != nil {
		return exrr
	}

	if file.Name.Name == "builtin" {

		ast.Print(fset, file)
	}

	ctx := &parseContext{
		File:        file,
		Env:         env,
		PackagesMap: make(map[string]*Package),
	}

	parseFileName(ctx)

	decls := file.Decls
	for _, d := range decls {
		switch c := d.(type) {
		case *ast.GenDecl:
			parseGenDecl(ctx, c)
		case *ast.FuncDecl:
			parseFuncDecl(ctx, c)
		}
	}
	return
}

func parseFileName(ctx *parseContext) {

	pkge, ok := ctx.Env.PackageByName(ctx.File.Name.Name)
	if !ok {
		var comments []string
		if ctx.File.Doc != nil {
			for _, t := range ctx.File.Comments {
				parseComments(t, &comments)
			}
		}
		pkg := &Package{
			Name:    ctx.File.Name.Name,
			Comment: comments,
		}
		ctx.Env.AppendPackage(pkg)
		ctx.PackagesMap[ctx.File.Name.Name] = pkg
	} else {
		ctx.PackagesMap[ctx.File.Name.Name] = pkge
	}
}

func parseComments(doc *ast.CommentGroup, c *[]string) {
	if doc.List == nil {
		return
	}

	sizeList := len(doc.List)
	if len(*c) != 0 {
		t := make([]string, sizeList)
		for i := 0; i < sizeList; i++ {
			t[i] = doc.List[i].Text
		}
		*c = append(*c, t...)
		return
	}
	*c = make([]string, sizeList)
	for i := 0; i < sizeList; i++ {
		(*c)[i] = doc.List[i].Text
	}
}

func parseGenDecl(ctx *parseContext, s *ast.GenDecl) {
	var comments []string
	if s.Doc != nil {
		parseComments(s.Doc, &comments)
	}

	for _, spec := range s.Specs {
		parseSpec(ctx, spec, &comments)
	}
}

func parseSpec(ctx *parseContext, spec ast.Spec, comments *[]string) (exrr error) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]

	switch s := spec.(type) {
	case *ast.TypeSpec:
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := NewStruct(currentPackage, s.Name.Name)
			declStruct.Comment = *comments
			refType := getRefType(ctx, s.Name.Name)

			refType.AppendType(declStruct)

			parseStruct(ctx, t, declStruct)
			currentPackage.Types = append(currentPackage.Types, declStruct)
		case *ast.Ident:
			if ctx.File.Name.Name == "builtin" {
				currentPackage.AppendRefType(t.Name)
			}
		}

	case *ast.ImportSpec:
		namePkg := s.Path.Value[1 : len(s.Path.Value)-1]
		if s.Name != nil {
			_, ok := ctx.PackageByName(s.Name.Name)
			if !ok {
				newPkg := &Package{}
				newPkg.Name = namePkg
				if s.Name.Name != "." {
					ctx.PackagesMap[s.Name.Name] = newPkg
					ctx.Env.AppendPackage(newPkg)
				} else {

					if _, isTrue := ctx.Env.PackageByName(namePkg); isTrue {
						return
					}

					basePath := ""
					if basePath, exrr = ctx.Env.basePath(); exrr != nil {
						return exrr
					}

					basePath = fmt.Sprintf("%s/%s", basePath, namePkg)
					if exrr = ctx.Env.ParsePackage(basePath, false); exrr != nil {
						return exrr
					}
					ctx.PackagesMap[namePkg], _ = ctx.Env.PackageByName(namePkg)
				}
			}
			getRefType(ctx, s.Name.Name)
		} else {
			_, ok := ctx.PackageByName(namePkg)
			if !ok {
				newPkg := &Package{
					Name: namePkg,
				}
				ctx.PackagesMap[namePkg] = newPkg
				_, ok := ctx.Env.PackageByName(namePkg)
				if !ok {
					ctx.Env.AppendPackage(newPkg)
				}
			}
			getRefType(ctx, namePkg)
		}
	case *ast.ValueSpec:
		//parseVariable(currentPackage, s)
	}
	return
}

func parseStruct(ctx *parseContext, astStruct *ast.StructType, typeStruct *Struct) {

	for _, field := range astStruct.Fields.List {
		refType := &RefType{}

		switch t := field.Type.(type) {
		case *ast.Ident:
			refType = getRefType(ctx, t.Name)
		case *ast.SelectorExpr:
			refType = getRefType(ctx, t.X.(*ast.Ident).Name)
		case *ast.StarExpr:
			refType = getRefType(ctx, t.X.(*ast.Ident).Name)
		} //interface conversion: ast.Expr is *ast.SelectorExpr, not *ast.Ident

		f := &Field{}
		f.Type = refType
		f.Tag.Raw = ""
		if field.Doc != nil {
			var comments []string
			parseComments(field.Doc, &comments)
			f.Comment = comments
		}

		if len(field.Names) > 0 { // TODO(jack): To check/understand multiple names.
			f.Name = field.Names[0].Name
		}

		if field.Tag != nil && field.Tag.Value != "" {
			f.Tag.Raw = field.Tag.Value[1 : len(field.Tag.Value)-1]

			structTag, err := structtag.Parse(f.Tag.Raw)
			if err != nil {
				fmt.Println("Error in format StructTag.")
				panic(err)
			}

			for _, tag := range structTag.Tags() {
				tp := &TagParam{}
				size := len(tag.Options)

				tp.Name = tag.Key
				tp.Value = tag.Name
				tp.Options = make([]string, size)

				if size != 0 {
					for i := 0; i < size; i++ {
						tp.Options[i] = tag.Options[i]
					}
				}
				f.Tag.AppendTagParam(tp)
			}
		}
		typeStruct.Fields = append(typeStruct.Fields, f)
	}
}

func parseFuncDecl(ctx *parseContext, f *ast.FuncDecl) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]
	method := NewMethodDescriptor(currentPackage, f.Name.Name)

	if f.Recv != nil { //TODO(check): I don't know if this pointers always will exist if "f.Recv" is diff than nil.
		recvList := f.Recv.List
		for _, field := range recvList {
			recv := MethodArgument{}
			typeName := field.Type.(*ast.StarExpr).X.(*ast.Ident).Name

			recv.Name = field.Names[0].Name
			recv.Type = getRefType(ctx, typeName)
			method.Recv = append(method.Recv, recv)
		}
	}

	if f.Doc != nil {
		var comments []string
		parseComments(f.Doc, &comments)
		method.Comment = comments
	}

	for _, field := range f.Type.Params.List {
		argument := MethodArgument{}
		if len(field.Names) > 0 {
			argument.Name = field.Names[0].Name
		}

		switch t := field.Type.(type) {
		case *ast.Ident:
			argument.Type = getRefType(ctx, t.Name)
		case *ast.StarExpr:
			argument.Type = getRefType(ctx, t.X.(*ast.Ident).Name)
		}

		method.Arguments = append(method.Arguments, argument)
	}
	currentPackage.Methods = append(currentPackage.Methods, method)
}

func parseVariable(parent *Package, f *ast.ValueSpec) {
	variable := &Variable{}
	varType := NewRefType(parent)

	variable.Name = f.Names[0].Name

	if f.Names == nil {
		varType.Name = f.Type.(*ast.Ident).Name
		variable.Type = varType
	} else {
		for _, value := range f.Values {
			switch v := value.(type) {
			case *ast.BasicLit:
				//varType.Name = TODO: Set values
				fmt.Printf("%T\n", v.Kind) //TODO: Convert token.token to string
			case *ast.Ident:
				//varType.Name = TODO: Set values
				fmt.Printf("%T\n", v.Name)
			}
		}
	}
}

/* 	### This method is temp ###
I don't know if always is necessary check types in builtin package and current package.
*/
func getRefType(ctx *parseContext, name string) *RefType {
	builtinPackage, ok := ctx.Env.PackageByName("builtin")
	if ok {
		refType, ok := builtinPackage.RefTypeByName(name)
		if ok {
			return refType
		}
	}

	//fmt.Println(ctx.File.Name.Name + "|" + name)

	packageCurrent := ctx.PackagesMap[ctx.File.Name.Name]

	//fmt.Println(packageCurrent.Name)

	for _, pr := range packageCurrent.RefType {
		if name == pr.Name && packageCurrent == pr.Pkg {
			return pr
		}
	}
	return packageCurrent.AppendRefType(name)
}
