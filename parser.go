package myasthurts

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"strings"

	"github.com/fatih/structtag"
)

func parseFileName(ctx *parseFileContext) error {
	pkg := ctx.Package
	if pkg.explored { // If the package is already explored, ignore this.
		return nil
	}
	var comments []string
	if ctx.File.Doc != nil {
		for _, t := range ctx.File.Comments {
			rComments, err := parseComments(t)
			if err != nil {
				return err
			}
			comments = append(comments, rComments...)
		}

		pkg.Doc = Doc{
			Comments: comments,
		}
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

func parseGenDecl(ctx *parseFileContext, s *ast.GenDecl) error {
	var (
		docs []string
		err  error
	)

	if s.Doc != nil {
		if docs, err = parseComments(s.Doc); err != nil {
			return err
		}
	}

	for _, spec := range s.Specs {
		err = parseSpec(ctx, spec, docs)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseSpec(ctx *parseFileContext, spec ast.Spec, docComments []string) error {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		nameType := s.Name.Name
		switch t := s.Type.(type) {
		case *ast.StructType:
			declStruct := NewStruct(ctx.Package, s.Name.Name)
			declStruct.Doc = Doc{
				Comments: docComments,
			}

			// Get the refType from the package.
			refType, ok := ctx.Package.RefTypeByName(s.Name.Name)
			if ok {
				// If the refType exists...
				if refType.Type() != nil { // if the refType is already resolved
					bt, ok := refType.Type().(*baseType)
					if !ok { // That means a double declaration or some unexpected error...
						return fmt.Errorf("type %t was not expected", refType.Type())
					}
					// Since it is a baseType, we should make it specific
					// (struct) and use its already defined methods ...
					declStruct.baseType = *bt
				}
				refType.AppendType(declStruct) // Realizes the refType
			} else {
				// If the ref type does not exists, creates and registers it.
				refType = NewRefType(s.Name.Name, ctx.Package, declStruct)
				ctx.Package.AddRefType(refType)
			}

			err := parseStruct(ctx, t, declStruct)
			if err != nil {
				return err
			}
			ctx.Package.Structs = append(ctx.Package.Structs, declStruct)
			ctx.Package.Types = append(ctx.Package.Types, declStruct)
		case *ast.Ident:
			if ctx.File.Name.Name == "builtin" {
				if nameType != t.Name {
					ctx.Package.AppendRefType(nameType)
				} else {
					ctx.Package.AppendRefType(t.Name)
				}
			}
		}

	case *ast.ImportSpec:
		importPathPkg := s.Path.Value[1 : len(s.Path.Value)-1]

		// Tries to find the package on the list...
		pkg, pkgExists := ctx.Env.PackageByImportPath(importPathPkg)

		// TODO(jota): To parametrize the import source directory.
		buildPackage, err := ctx.Env.BuildContext.Import(importPathPkg, ".", build.ImportComment)
		if err != nil {
			return err
		}

		if !pkgExists {
			pkg = NewPackage(buildPackage)
			ctx.Env.AppendPackage(pkg)
		}

		if s.Name != nil { // The name is the identifier of the import. Ex: t "time", t would be the name
			// This checks if the import is a dot import. That means we have
			// to include this imported package into a special list for
			// prior querying. Dot imports include all types declared into
			// the same contexts for this file. So, types don't have the
			// package identification.
			if s.Name.Name == "." {
				pkgCtx := NewPackageContext(pkg, buildPackage)
				if err = ctx.Env.parsePackage(pkgCtx); err != nil {
					return err
				}
				ctx.dotImports = append(ctx.dotImports, pkg) // If we do explore, it means the package is dot imported.
			} else {
				// Sets the alias of the package for this file context.
				ctx.packageImportAliasMap[s.Name.Name] = pkg
			}
		}
	case *ast.ValueSpec:
		if ctx.Package.Name != "builtin" {
			vrle := parseVariable(ctx.Package, s)
			ctx.Package.AppendVariable(vrle)
		}
	}
	return nil
}

func parseStruct(ctx *parseFileContext, astStruct *ast.StructType, typeStruct *Struct) error {
	currentPackage := ctx.Package

	for _, field := range astStruct.Fields.List {
		var (
			refType RefType
			err     error
		)

		// TODO(jota): Add a default case returning an error.
		switch t := field.Type.(type) {
		case *ast.Ident: // this is a type like int64
			if refType, err = ctx.GetRefType(t.Name); err != nil {
				if refType = currentPackage.AppendRefType(t.Name); refType == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.SelectorExpr:
			if refType, err = ctx.GetRefType(t.X.(*ast.Ident).Name); err != nil {
				if refType = currentPackage.AppendRefType(t.X.(*ast.Ident).Name); refType == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.StarExpr: // This is a pointer declaration.
			// TODO(jota): A new RefType type should be created to express pointers.
			// TODO(jota): Add a default case returning an error.
			switch xType := t.X.(type) { // TODO: Check this type
			case *ast.Ident:
				if refType, err = ctx.GetRefType(xType.Name); err != nil {
					if refType = currentPackage.AppendRefType(xType.Name); refType == nil {
						return errors.New("Append Reftype error")
					}
				}
			}
		}

		f := &Field{}
		f.RefType = refType
		if field.Doc != nil {
			var docComments []string
			if docComments, err = parseComments(field.Doc); err != nil {
				return err
			}
			f.Doc = Doc{
				Comments: docComments,
			}
		}

		if len(field.Names) > 0 {
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

func parseFuncDecl(ctx *parseFileContext, f *ast.FuncDecl) (err error) {
	method := NewMethodDescriptor(ctx.Package, f.Name.Name)
	hasReceiver := f.Recv != nil && len(f.Recv.List) > 0
	if hasReceiver {
		field := f.Recv.List[0]

		recv := MethodArgument{}
		var refType RefType
		switch recvT := field.Type.(type) {
		case *ast.Ident:
			// TODO(jota): To check for non pointer types.
		case *ast.StarExpr:
			typeName := recvT.X.(*ast.Ident).Name
			rt, ok := ctx.Package.RefTypeByName(typeName)
			if !ok {
				rt = NewRefType(typeName, ctx.Package, NewBaseType(ctx.Package, typeName))
				ctx.Package.AddRefType(rt)
			}
			refType = rt
		}

		recv.Name = field.Names[0].Name
		recv.Type = refType
		method.Recv = append(method.Recv, recv)

		// Add method to the type...
		refType.Type().AddMethod(&TypeMethod{
			Name:       field.Names[0].Name,
			Descriptor: method,
		})
	} else {
		ctx.Package.Methods = append(ctx.Package.Methods, method) // TODO(jota): This might not be a package method.
	}

	// Set the method documentation.
	if f.Doc != nil {
		var comments []string
		if comments, err = parseComments(f.Doc); err != nil {
			return err
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
			if argument.Type, err = ctx.GetRefType(t.Name); err != nil {
				if argument.Type = ctx.Package.AppendRefType(t.Name); argument.Type == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.StarExpr:
			switch xType := t.X.(type) {
			case *ast.Ident:
				if argument.Type, err = ctx.GetRefType(xType.Name); err != nil {
					if argument.Type = ctx.Package.AppendRefType(xType.Name); argument.Type == nil {
						return errors.New("Append Reftype error")
					}
				}
			case *ast.SelectorExpr:
				fmt.Println(xType.X.(*ast.Ident).Name) // TODO: Check this type
			}
		}

		method.Arguments = append(method.Arguments, argument)
	}
	return nil
}

func parseVariable(parent *Package, f *ast.ValueSpec) (vrle *Variable) {
	vrle = &Variable{
		Name: f.Names[0].Name,
	}

	switch t := f.Type.(type) {
	case *ast.Ident:
		if refType, ok := parent.RefTypeByName(t.Name); ok {
			vrle.RefType = refType
			return
		}
		variable.RefType = parent.AppendRefType(t.Name)
	case *ast.ArrayType:
		n := t.Elt.(*ast.Ident).Name
		if refType, ok := parent.RefTypeByName(n); ok {
			vrle.RefType = refType
			return
		}
	}

	for _, value := range f.Values {
		switch v := value.(type) {
		case *ast.BasicLit:
			typeName := strings.ToLower(v.Kind.String())
			if refType, ok := parent.RefTypeByName(typeName); ok {
				vrle.RefType = refType
				return
			}
			variable.RefType = parent.AppendRefType(typeName)
		}
	}

	return variable
}
