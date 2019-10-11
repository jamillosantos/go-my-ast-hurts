package myasthurts

import (
	"errors"
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
	builtinPath := fmt.Sprintf("%s/builtin", path)
	env.ParsePackage(builtinPath, false)
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

func parseFileName(ctx *parseContext) (exrr error) {

	pkge, ok := ctx.Env.PackageByName(ctx.File.Name.Name)
	if !ok {
		var (
			rComments []string
			comments  []string
		)
		if ctx.File.Doc != nil {
			for _, t := range ctx.File.Comments {
				if rComments, exrr = parseComments(t); exrr != nil {
					return exrr
				}
				comments = append(comments, rComments...)
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
	return nil
}

func parseComments(doc *ast.CommentGroup) (r []string, exrr error) {
	if doc.List == nil {
		return nil, errors.New("Doc list empty")
	}

	sizeList := len(doc.List)

	t := make([]string, sizeList)
	for i := 0; i < sizeList; i++ {
		t[i] = doc.List[i].Text
	}
	return t, nil
}

func parseGenDecl(ctx *parseContext, s *ast.GenDecl) (exrr error) {
	var comments []string

	if s.Doc != nil {
		if comments, exrr = parseComments(s.Doc); exrr != nil {
			return exrr
		}
	}

	for _, spec := range s.Specs {
		parseSpec(ctx, spec, comments)
	}
	return nil
}

func parseSpec(ctx *parseContext, spec ast.Spec, comments []string) (exrr error) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]
	var refType *RefType

	switch s := spec.(type) {
	case *ast.TypeSpec:
		nameType := s.Name.Name
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := NewStruct(currentPackage, s.Name.Name)
			declStruct.Comment = comments

			if refType, exrr = ctx.GetRefType(s.Name.Name); exrr != nil {
				if refType, exrr = currentPackage.AppendRefType(s.Name.Name); exrr != nil {
					return exrr
				}
			}

			refType.AppendType(declStruct)

			parseStruct(ctx, t, declStruct)
			currentPackage.Types = append(currentPackage.Types, declStruct)
		case *ast.Ident:
			if ctx.File.Name.Name == "builtin" {
				if nameType != t.Name {
					currentPackage.AppendRefType(nameType)
				} else {
					currentPackage.AppendRefType(t.Name)
				}
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
					ctx.PackagesMap[s.Name.Name+namePkg], _ = ctx.Env.PackageByName(namePkg)
				}
			}

			if refType, exrr = ctx.GetRefType(s.Name.Name); exrr != nil {
				if refType, exrr = currentPackage.AppendRefType(s.Name.Name); exrr != nil {
					return exrr
				}
			}
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
			if refType, exrr = ctx.GetRefType(namePkg); exrr != nil {
				if refType, exrr = currentPackage.AppendRefType(namePkg); exrr != nil {
					return exrr
				}
			}
		}
	case *ast.ValueSpec:
		//parseVariable(currentPackage, s)
	}
	return
}

func parseStruct(ctx *parseContext, astStruct *ast.StructType, typeStruct *Struct) (exrr error) {

	currentPackage, ok := ctx.PackageByName(ctx.File.Name.Name)
	if !ok {
		return errors.New("Package not found during parse Struct")
	}

	for _, field := range astStruct.Fields.List {
		var refType *RefType

		switch t := field.Type.(type) {
		case *ast.Ident:
			if refType, exrr = ctx.GetRefType(t.Name); exrr != nil {
				if refType, exrr = currentPackage.AppendRefType(t.Name); exrr != nil {
					return exrr
				}
			}
		case *ast.SelectorExpr:
			if refType, exrr = ctx.GetRefType(t.X.(*ast.Ident).Name); exrr != nil {
				if refType, exrr = currentPackage.AppendRefType(t.X.(*ast.Ident).Name); exrr != nil {
					return exrr
				}
			}
		case *ast.StarExpr:
			switch xType := t.X.(type) { // TODO: Check this type
			case *ast.Ident:
				if refType, exrr = ctx.GetRefType(xType.Name); exrr != nil {
					if refType, exrr = currentPackage.AppendRefType(xType.Name); exrr != nil {
						return exrr
					}
				}
			}
		}

		f := &Field{}
		f.Type = refType
		f.Tag.Raw = ""
		if field.Doc != nil {
			var comments []string
			if comments, exrr = parseComments(field.Doc); exrr != nil {
				return exrr
			}
			f.Comment = comments
		}

		if len(field.Names) > 0 { // TODO(Jack): To check/understand multiple names.
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
	return nil
}

func parseFuncDecl(ctx *parseContext, f *ast.FuncDecl) (exrr error) {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]
	method := NewMethodDescriptor(currentPackage, f.Name.Name)

	if f.Recv != nil { //TODO(check): I don't know if this pointers always will exist if "f.Recv" is diff than nil.
		recvList := f.Recv.List
		for _, field := range recvList {
			recv := MethodArgument{}
			typeName := ""
			switch recvT := field.Type.(type) {
			case *ast.StarExpr:
				typeName = recvT.X.(*ast.Ident).Name
			}

			recv.Name = field.Names[0].Name
			if recv.Type, exrr = ctx.GetRefType(typeName); exrr != nil {
				if recv.Type, exrr = currentPackage.AppendRefType(typeName); exrr != nil {
					return exrr
				}
			}
			method.Recv = append(method.Recv, recv)
		}
	}

	if f.Doc != nil {
		var comments []string
		if comments, exrr = parseComments(f.Doc); exrr != nil {
			return exrr
		}
		method.Comment = comments
	}

	for _, field := range f.Type.Params.List {
		argument := MethodArgument{}
		if len(field.Names) > 0 {
			argument.Name = field.Names[0].Name
		}

		switch t := field.Type.(type) {
		case *ast.Ident:
			if argument.Type, exrr = ctx.GetRefType(t.Name); exrr != nil {
				if argument.Type, exrr = currentPackage.AppendRefType(t.Name); exrr != nil {
					return exrr
				}
			}
		case *ast.StarExpr:
			switch xType := t.X.(type) {
			case *ast.Ident:
				if argument.Type, exrr = ctx.GetRefType(xType.Name); exrr != nil {
					if argument.Type, exrr = currentPackage.AppendRefType(xType.Name); exrr != nil {
						return exrr
					}
				}
			case *ast.SelectorExpr:
				fmt.Println(xType.X.(*ast.Ident).Name) // TODO: Check this type
			}
		}

		method.Arguments = append(method.Arguments, argument)
	}
	currentPackage.Methods = append(currentPackage.Methods, method)
	return nil
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
