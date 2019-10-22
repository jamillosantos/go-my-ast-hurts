package myasthurts

import (
	"fmt"
	"go/ast"
	"go/build"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
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
		} else {
			// If it does not have any name defined, use the original package name.
			ctx.packageImportAliasMap[buildPackage.Name] = pkg
		}
	case *ast.ValueSpec:
		variable, err := parseVariable(ctx, s)
		if err != nil {
			return err
		}
		ctx.Package.AppendVariable(variable)
	}
	return nil
}

func parseStruct(ctx *parseFileContext, astStruct *ast.StructType, typeStruct *Struct) error {
	for _, field := range astStruct.Fields.List {
		refType, err := parseType(ctx, field.Type)
		if err != nil {
			return err
		}

		f := &Field{
			RefType: refType,
		}

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

func parseFuncDecl(ctx *parseFileContext, f *ast.FuncDecl) error {
	method := NewMethodDescriptor(ctx.Package, f.Name.Name)
	hasReceiver := f.Recv != nil && len(f.Recv.List) > 0
	if hasReceiver {
		field := f.Recv.List[0]

		recv := MethodArgument{}
		refType, err := parseType(ctx, field.Type)
		if err != nil {
			return err
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
		docComments, err := parseComments(f.Doc)
		if err != nil {
			return err
		}
		method.Doc = Doc{
			Comments: docComments,
		}
	}

	for _, field := range f.Type.Params.List {
		argument := MethodArgument{}
		if len(field.Names) > 0 {
			argument.Name = field.Names[0].Name
		}

		refType, err := parseType(ctx, field.Type)
		if err != nil {
			return err
		}

		argument.Type = refType

		method.Arguments = append(method.Arguments, argument)
	}
	return nil
}

// parseType will return a RefType for a given type.
func parseType(ctx *parseFileContext, t ast.Expr) (RefType, error) {
	switch recvT := t.(type) {
	// This case will cover the identifier type. This is for string, int64 and
	// types defined on the same package.
	case *ast.Ident:
		typeName := recvT.Name
		rt, ok := ctx.GetRefType(typeName)
		if !ok {
			rt = NewRefType(typeName, ctx.Package, NewBaseType(ctx.Package, typeName))
			ctx.Package.AddRefType(rt)
		}
		return rt, nil
	// This case will cover the selector type. This is for expressions like
	// time.Time or t.Time, models.User ...
	case *ast.SelectorExpr:
		pkgAliasIdent, ok := recvT.X.(*ast.Ident)
		if !ok { // We expect the recvT.X is a ast.Ident, if not...
			// TODO(jota): To formalize this error.
			return nil, errors.New("unexpected selector identifier")
		}
		pkgAlias, ok := ctx.PackageByImportAlias(pkgAliasIdent.Name)
		if !ok { // The package does not exists in the ctx?? It should not be happening...
			// TODO(jota): To formalize this error.
			return nil, fmt.Errorf("package %s was not found", pkgAliasIdent.Name)
		}
		refType, _ := pkgAlias.EnsureRefType(recvT.Sel.Name) // We don't care if the refType is created now or not.
		return refType, nil
	case *ast.InterfaceType:
		return InterfaceRefType, nil
	// This case covers pointers. It is recursive because pointers can be for
	// identifiers or selectors...
	case *ast.StarExpr:
		// TODO(jota): To create a RefType that represents a pointer and wraps the result before returning.
		return parseType(ctx, recvT.X)
	case *ast.ArrayType:
		// TODO(jota): To create a RefType that represents an array and wraps the result before returning.
		return parseType(ctx, recvT.Elt)
	case *ast.MapType:
		// TODO(jota): To create a RefType that represents a map and wraps the result before returning.
		return parseType(ctx, recvT.Value)
	case *ast.ChanType:
		// TODO(jota): To create a RefType that represents a channel and wraps the result before returning.
		return parseType(ctx, recvT.Value)
	case *ast.FuncType:
		return parseFuncType(ctx, recvT)
	case *ast.Ellipsis:
		// TODO(jota): To create a RefType that represents an Ellipsis and wraps the result before returning.
		return parseType(ctx, recvT.Elt)
	// This is a safeguard for unexpected cases.
	default:
		// TODO(jota): To formalize this error.
		return nil, fmt.Errorf("%T: unexpected expression type", t)
	}
}

func parseFuncType(ctx *parseFileContext, f *ast.FuncType) (RefType, error) {
	md := &MethodDescriptor{
		baseType: *NewBaseType(ctx.Package, ""),
	}

	for _, p := range f.Params.List {
		refType, err := parseType(ctx, p.Type)
		if err != nil {
			return nil, err
		}

		methodArg := MethodArgument{
			Type: refType,
		}

		if p.Names != nil {
			methodArg.Name = p.Names[0].Name
		}

		if p.Doc != nil {
			docComments, err := parseComments(p.Doc)
			if err != nil {
				return nil, err
			}
			methodArg.Doc = Doc{
				Comments: docComments,
			}
		}

		md.Arguments = append(md.Arguments, methodArg)
	}

	return NewRefType("", ctx.Package, md), nil
}

func parseVariable(ctx *parseFileContext, vValue *ast.ValueSpec) (*Variable, error) {
	variable := &Variable{
		Name: vValue.Names[0].Name,
	}

	// Defines the variable documentation...
	if vValue.Doc != nil {
		docComments, err := parseComments(vValue.Doc)
		if err != nil {
			return nil, err
		}
		variable.Doc = Doc{
			Comments: docComments,
		}
	}

	if vValue.Type == nil {
		variable.RefType = NullRefType
	} else {
		// Define and set the RefType of the variable.
		refType, err := parseType(ctx, vValue.Type)
		if err != nil {
			return nil, err
		}
		variable.RefType = refType
	}

	return variable, nil
}
