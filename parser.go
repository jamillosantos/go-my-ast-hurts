package myasthurts

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
)

func (ctx *parseContext) PackageByName(name string) (*Package, bool) {
	pkg, ok := ctx.PackagesMap[name]
	return pkg, ok
}

func (ctx *parseContext) GetRefType(name string) (ref *RefType, exrr error) {

	builtin, ok := ctx.Env.PackageByName("builtin")
	if !ok {
		return nil, errors.New("Package builtin not found")
	}

	if ref = builtin.RefTypeByName(name); ref != nil {
		return ref, nil
	}

	for _, p := range ctx.PackagesMap {
		if ref = p.RefTypeByName(name); ref != nil {
			return ref, nil
		}
	}
	return nil, errors.New("RefType not found")
}

func (env *environment) makeEnv() error {
	var (
		err error
	)
	builtinPath := "builtin"
	if err = env.Parse(builtinPath); err != nil {
		return err
	}
	return nil
}

func (env *environment) parse(filePath string) error {
	var (
		file *ast.File
		fset *token.FileSet
		err  error
	)

	fset = token.NewFileSet()
	if file, err = parser.ParseFile(fset, filePath, nil, parser.ParseComments); err != nil {
		return err
	}

	ctx := &parseContext{
		File:        file,
		Env:         env,
		PackagesMap: make(map[string]*Package),
	}

	if env.Config.DevMode && env.Config.ASTI {
		ast.Print(fset, file)
	}

	err = parseFileName(ctx)
	if err != nil {
		return err
	}

	decls := file.Decls
	for _, d := range decls {
		switch c := d.(type) {
		case *ast.GenDecl:
			err = parseGenDecl(ctx, c)
			if err != nil {
				return err
			}
		case *ast.FuncDecl:
			err = parseFuncDecl(ctx, c)
			if err != nil {
				return err
			}
			// TODO(jota): We should have a default case here. Shall it return an error?
		}
	}
	return nil
}

func parseFileName(ctx *parseContext) error {
	pkg, ok := ctx.Env.PackageByName(ctx.File.Name.Name)
	if !ok {
		var comments []string
		if ctx.File.Doc != nil {
			for _, t := range ctx.File.Comments {
				rComments, err := parseComments(t)
				if err != nil {
					return err
				}
				comments = append(comments, rComments...)
			}
		}
		pkg = &Package{
			Name: ctx.File.Name.Name,
			Doc: Doc{
				Comments: comments,
			},
		}
		ctx.Env.AppendPackage(pkg)
		ctx.PackagesMap[ctx.File.Name.Name] = pkg
	} else {
		ctx.PackagesMap[ctx.File.Name.Name] = pkg
	}
	return nil
}

func parseComments(doc *ast.CommentGroup) (r []string, exrr error) {
	sizeList := len(doc.List)

	t := make([]string, sizeList)
	for i := 0; i < sizeList; i++ {
		t[i] = doc.List[i].Text
	}
	return t, nil
}

func parseGenDecl(ctx *parseContext, s *ast.GenDecl) error {
	var (
		comments []string
		err      error
	)

	if s.Doc != nil {
		if comments, err = parseComments(s.Doc); err != nil {
			return err
		}
	}

	for _, spec := range s.Specs {
		err = parseSpec(ctx, spec, comments)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseSpec(ctx *parseContext, spec ast.Spec, comments []string) error {
	currentPackage := ctx.PackagesMap[ctx.File.Name.Name]
	var (
		refType *RefType
		err     error
	)

	switch s := spec.(type) {
	case *ast.TypeSpec:
		nameType := s.Name.Name
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := NewStruct(currentPackage, s.Name.Name)
			declStruct.Doc = Doc{
				Comments: comments,
			}

			if refType, err = ctx.GetRefType(s.Name.Name); err != nil {
				if refType = currentPackage.AppendRefType(s.Name.Name); refType == nil {
					return errors.New("Append Reftype error")
				}
			}

			refType.AppendType(declStruct)

			err = parseStruct(ctx, t, declStruct)
			if err != nil {
				return err
			}
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

				// This checks if the import is a dot import. That means we have
				// to include this imported package into a special list for
				// prior querying. Dot imports include all types declared into
				// the same contexts for this file. So, types don't have the
				// package identification.
				if s.Name.Name == "." {
					// TODO(jota): Does that mean the package was already explored?
					if _, ok := ctx.Env.PackageByName(namePkg); ok {
						return nil
					}

					// TODO(jota): To parametrize the import source directory.
					pkg, err := ctx.Env.BuildContext.Import(namePkg, ".", build.FindOnly)
					if err != nil {
						return err
					}

					// Parse the package.
					if err = ctx.Env.parsePackage(pkg); err != nil {
						return err
					}
					ctx.PackagesMap[pkg.ImportPath], _ = ctx.Env.PackageByName(namePkg)
				} else {
					ctx.PackagesMap[s.Name.Name] = newPkg
					ctx.Env.AppendPackage(newPkg)
				}
			}

			if refType, err = ctx.GetRefType(s.Name.Name); err != nil {
				if refType = currentPackage.AppendRefType(s.Name.Name); refType == nil {
					return errors.New("Append Reftype error")
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
			if refType, err = ctx.GetRefType(namePkg); err != nil {
				if refType = currentPackage.AppendRefType(namePkg); refType == nil {
					return errors.New("Append Reftype error")
				}
			}
		}
	case *ast.ValueSpec:
		if currentPackage.Name != "builtin" {
			vrle := parseVariable(currentPackage, s)
			currentPackage.AppendVariable(vrle)
		}
	}
	return nil
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
				if refType = currentPackage.AppendRefType(t.Name); refType == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.SelectorExpr:
			if refType, exrr = ctx.GetRefType(t.X.(*ast.Ident).Name); exrr != nil {
				if refType = currentPackage.AppendRefType(t.X.(*ast.Ident).Name); refType == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.StarExpr:
			switch xType := t.X.(type) { // TODO: Check this type
			case *ast.Ident:
				if refType, exrr = ctx.GetRefType(xType.Name); exrr != nil {
					if refType = currentPackage.AppendRefType(xType.Name); refType == nil {
						return errors.New("Append Reftype error")
					}
				}
			}
		}

		f := &Field{}
		f.RefType = refType
		f.Tag.Raw = ""
		if field.Doc != nil {
			var comments []string
			if comments, exrr = parseComments(field.Doc); exrr != nil {
				return exrr
			}
			f.Doc = Doc{
				Comments: comments,
			}
		}

		if len(field.Names) > 0 { // TODO(Jack): To check/understand multiple names.
			f.Name = field.Names[0].Name
		}

		if field.Tag != nil && field.Tag.Value != "" {
			f.Tag.Raw = field.Tag.Value[1 : len(field.Tag.Value)-1]

			structTag, err := structtag.Parse(f.Tag.Raw)
			if err != nil {
				return err
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

	if f.Recv != nil {
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
				if recv.Type = currentPackage.AppendRefType(typeName); recv.Type == nil {
					return errors.New("Append Reftype error")
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
		method.Doc = Doc{
			Comments: comments,
		}
	}

	for _, field := range f.Type.Params.List {
		argument := MethodArgument{}
		if len(field.Names) > 0 {
			argument.Name = field.Names[0].Name
		}

		switch t := field.Type.(type) {
		case *ast.Ident:
			if argument.Type, exrr = ctx.GetRefType(t.Name); exrr != nil {
				if argument.Type = currentPackage.AppendRefType(t.Name); argument.Type == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.StarExpr:
			switch xType := t.X.(type) {
			case *ast.Ident:
				if argument.Type, exrr = ctx.GetRefType(xType.Name); exrr != nil {
					if argument.Type = currentPackage.AppendRefType(xType.Name); argument.Type == nil {
						return errors.New("Append Reftype error")
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

func parseVariable(parent *Package, f *ast.ValueSpec) (vrle *Variable) {
	var (
		refType *RefType
	)
	variable := &Variable{
		Name: f.Names[0].Name,
	}

	switch t := f.Type.(type) {
	case *ast.Ident:
		refType = parent.RefTypeByName(t.Name)
		if refType != nil {
			variable.RefType = refType
			return
		}
		variable.RefType = parent.AppendRefType(t.Name)
	case *ast.ArrayType:
		n := t.Elt.(*ast.Ident).Name
		refType = parent.RefTypeByName(n)
		if refType != nil {
			variable.RefType = refType
			//I don't know why the "return" causes error here
		}
	}

	for _, value := range f.Values {
		switch v := value.(type) {
		case *ast.BasicLit:
			typeName := strings.ToLower(v.Kind.String())
			refType = parent.RefTypeByName(typeName)
			if refType != nil {
				variable.RefType = refType
				return
			}
			variable.RefType = parent.AppendRefType(typeName)
		}
	}

	return variable
}
