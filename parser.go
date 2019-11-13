package myasthurts

import (
	"fmt"
	"go/ast"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
)

func parseFileName(ctx *ParseFileContext) error {
	pkg := ctx.Package
	if pkg.Explored { // If the package is already explored, ignore this.
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

func parseGenDecl(ctx *ParseFileContext, s *ast.GenDecl) error {
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

func parseInterface(ctx *ParseFileContext, name string, spec *ast.InterfaceType, docComments []string) (*Interface, error) {
	i := NewInterface(ctx.Package, name)
	for _, m := range spec.Methods.List {
		switch t := m.Type.(type) {
		case *ast.Ident:
			// This case is for composing interfaces for this same package.
			// TODO(jota): Parse interfaces and create a mechanism for "including methods".
		case *ast.SelectorExpr:
			// This case is for composing interfaces from other packages.
			// TODO(jota): Parse interfaces and create a mechanism for "including methods".
		case *ast.FuncType:
			name := ""
			if len(m.Names) > 0 {
				name = m.Names[0].Name
			}
			md, err := parseFuncType(ctx, name, t)
			if err != nil {
				return nil, err
			}
			i.AddMethod(&TypeMethod{
				Name:       md.Name(),
				Descriptor: md,
			})
		default:
			pos := ctx.FSet.Position(spec.Pos())
			return nil, errors.Wrapf(ErrUnexpectedExpressionType, "%T found while parsing %s (%s)", m.Type, name, pos.String())
		}
	}
	return i, nil
}

func parseSpec(ctx *ParseFileContext, spec ast.Spec, docComments []string) error {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		nameType := s.Name.Name
		switch t := s.Type.(type) {
		case *ast.InterfaceType:
			i, err := parseInterface(ctx, nameType, t, docComments)
			if err != nil {
				return err
			}
			ctx.Package.AppendInterface(i)
		case *ast.StructType:
			declStruct := NewStruct(ctx.Package, nameType)
			declStruct.Doc = Doc{
				Comments: docComments,
			}

			// Get the refType from the package.
			refType, ok := ctx.Package.RefTypeByName(nameType)
			if ok {
				// If the refType exists...
				if refType.Type() != nil { // if the refType is already resolved
					bt, ok := refType.Type().(*BaseType)
					if !ok { // That means a double declaration or some unexpected error...
						return fmt.Errorf("type %t was not expected", refType.Type())
					}
					// Since it is a baseType, we should make it specific
					// (struct) and use its already defined methods ...
					declStruct.BaseType = *bt
				}
				refType.AppendType(declStruct) // Realizes the refType
			} else {
				// If the ref type does not exists, creates and registers it.
				refType = NewRefType(nameType, ctx.Package, declStruct)
				ctx.Package.AddRefType(refType)
			}

			err := parseStruct(ctx, t, declStruct)
			if err != nil {
				return err
			}
			ctx.Package.AppendStruct(declStruct)
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

		buildPackage, err := ctx.Env.Import(importPathPkg)
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

func parseStruct(ctx *ParseFileContext, astStruct *ast.StructType, typeStruct *Struct) error {
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

func parseFuncDecl(ctx *ParseFileContext, f *ast.FuncDecl) error {
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
		ctx.Package.AppendMethod(method)
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

	if f.Type.Results != nil {
		for _, field := range f.Type.Results.List {
			r := MethodResult{}
			if len(field.Names) > 0 {
				r.Name = field.Names[0].Name
			}

			refType, err := parseType(ctx, field.Type)
			if err != nil {
				return err
			}

			r.Type = refType

			method.Result = append(method.Result, r)
		}
	}

	return nil
}

// parseType will return a RefType for a given type.
func parseType(ctx *ParseFileContext, t ast.Expr) (RefType, error) {
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
			return nil, errors.Wrapf(ErrUnexpectedSelector, "%T", recvT.X)
		}
		pkgAlias, ok := ctx.PackageByImportAlias(pkgAliasIdent.Name)
		if !ok { // The package does not exists in the ctx?? It should not be happening...
			return nil, errors.Wrap(ErrPackageAliasNotFound, pkgAliasIdent.Name)
		}
		refType, _ := pkgAlias.EnsureRefType(recvT.Sel.Name) // We don't care if the refType is created now or not.
		return refType, nil
	case *ast.InterfaceType:
		return InterfaceRefType, nil
	// This case covers pointers. It is recursive because pointers can be for
	// identifiers or selectors...
	case *ast.StarExpr:
		refType, err := parseType(ctx, recvT.X)
		if err != nil {
			return nil, err
		}
		return NewStarRefType(refType), nil
	case *ast.ArrayType: // TODO(Jeconias): Shall ArrayType be represented as a Type not a RefType?
		refType, err := parseType(ctx, recvT.Elt)
		if err != nil {
			return nil, err
		}
		return NewArrayRefType(refType), nil
	case *ast.MapType:

		keyRefType, err := parseType(ctx, recvT.Key)
		if err != nil {
			return nil, err
		}

		valueRefType, err := parseType(ctx, recvT.Value)
		if err != nil {
			return nil, err
		}

		mType := NewMap(ctx.Package, keyRefType, valueRefType)
		return NewRefType(mType.Name(), ctx.Package, mType), nil
	case *ast.ChanType: // TODO(Jeconias): Shall ChanType be represented as a Type not a RefType?
		refType, err := parseType(ctx, recvT.Value)
		if err != nil {
			return nil, err
		}
		return NewChanRefType(refType), nil
	case *ast.FuncType:
		md, err := parseFuncType(ctx, "", recvT)
		if err != nil {
			return nil, err
		}
		refType := NewRefType("", ctx.Package, md)
		return refType, nil
	case *ast.Ellipsis: // TODO(Jeconias): Shall Ellipsis be represented as a Type not a RefType?
		refType, err := parseType(ctx, recvT.Elt)
		if err != nil {
			return nil, err
		}
		return NewEllipsisRefType(refType), nil
	// This is a safeguard for unexpected cases.
	default:
		return nil, errors.Wrapf(ErrUnexpectedExpressionType, "%T", t)
	}
}

func parseFuncType(ctx *ParseFileContext, name string, f *ast.FuncType) (*MethodDescriptor, error) {
	md := &MethodDescriptor{
		BaseType: *NewBaseType(ctx.Package, name),
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

	if f.Results != nil {
		for _, r := range f.Results.List {
			name := ""
			if len(r.Names) > 0 {
				name = r.Names[0].Name
			}
			refType, err := parseType(ctx, r.Type)
			if err != nil {
				return nil, err
			}
			md.Result = append(md.Result, MethodResult{
				Name: name,
				Type: refType,
			})
		}
	}

	return md, nil
}

func parseVariable(ctx *ParseFileContext, vValue *ast.ValueSpec) (*Variable, error) {
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
