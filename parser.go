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
	currentPackage := ctx.Package
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
				Comments: docComments,
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
		// TODO(jota): To try simplifying this code.

		importPathPkg := s.Path.Value[1 : len(s.Path.Value)-1]
		if s.Name != nil { // The name is the identifier of the import. Ex: t "time", t would be the name
			// Tries to find the package on the list...
			_, ok := ctx.PackageByImportAlias(s.Name.Name)
			if !ok {
				// This checks if the import is a dot import. That means we have
				// to include this imported package into a special list for
				// prior querying. Dot imports include all types declared into
				// the same contexts for this file. So, types don't have the
				// package identification.
				if s.Name.Name == "." { // TODO(jota): This should be done before the right previous IF [ _, ok := ctx.PackageByImportAlias(s.Name.Name) ].
					// TODO(jota): Does that mean the package was already explored?
					if pkg, ok := ctx.Env.PackageByImportPath(importPathPkg); ok {
						// If the package already exists...
						ctx.dotImports = append(ctx.dotImports, pkg)
						return nil
					}

					// TODO(jota): To parametrize the import source directory.
					buildPackage, err := ctx.Env.BuildContext.Import(importPathPkg, ".", build.ImportComment)
					if err != nil {
						return err
					}

					newPkg := NewPackage(buildPackage)
					pkgCtx := NewPackageContext(newPkg, buildPackage)
					// Parse the package.
					if err = ctx.Env.parsePackage(pkgCtx); err != nil {
						return err
					}
					ctx.Env.AppendPackage(newPkg)
					// Add the package as dot imported on this context.
					ctx.dotImports = append(ctx.dotImports, newPkg)
				} else {

					// TODO(jota): To parametrize the import source directory.
					buildPackage, err := ctx.Env.BuildContext.Import(importPathPkg, ".", build.ImportComment)
					if err != nil {
						return err
					}

					newPkg := NewPackage(buildPackage)
					ctx.packageImportAliasMap[s.Name.Name] = newPkg
					ctx.Env.AppendPackage(newPkg)
				}
			}

			if refType, err = ctx.GetRefType(s.Name.Name); err != nil {
				if refType = currentPackage.AppendRefType(s.Name.Name); refType == nil {
					return errors.New("Append Reftype error")
				}
			}
		} else {
			_, ok := ctx.PackageByImportAlias(importPathPkg)
			if !ok { // The package is not in the memory yet
				// TODO(jota): To parametrize the import source directory.
				buildPackage, err := ctx.Env.BuildContext.Import(importPathPkg, ".", build.ImportComment)
				if err != nil {
					return err
				}

				// Since it is not a dot import, we don't have to explore it now.
				newPkg := NewPackage(buildPackage)
				ctx.Env.AppendPackage(newPkg)
			}

			// TODO(jota): I don't see the need of this now.
			/*
				if refType, err = ctx.GetRefType(namePkg); err != nil {
					if refType = currentPackage.AppendRefType(namePkg); refType == nil {
						return errors.New("Append Reftype error")
					}
				}
			*/
		}
	case *ast.ValueSpec:
		if currentPackage.Name != "builtin" {
			vrle := parseVariable(currentPackage, s)
			currentPackage.AppendVariable(vrle)
		}
	}
	return nil
}

func parseStruct(ctx *parseFileContext, astStruct *ast.StructType, typeStruct *Struct) error {
	currentPackage := ctx.Package

	for _, field := range astStruct.Fields.List {
		var (
			refType *RefType
			err     error
		)

		switch t := field.Type.(type) {
		case *ast.Ident:
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
		case *ast.StarExpr:
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
		f.Tag.Raw = ""
		if field.Doc != nil {
			var comments []string
			if comments, err = parseComments(field.Doc); err != nil {
				return err
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

func parseFuncDecl(ctx *parseFileContext, f *ast.FuncDecl) (err error) {
	currentPackage := ctx.Package
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
			if recv.Type, err = ctx.GetRefType(typeName); err != nil {
				if recv.Type = currentPackage.AppendRefType(typeName); recv.Type == nil {
					return errors.New("Append Reftype error")
				}
			}
			method.Recv = append(method.Recv, recv)
		}
	}

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
				if argument.Type = currentPackage.AppendRefType(t.Name); argument.Type == nil {
					return errors.New("Append Reftype error")
				}
			}
		case *ast.StarExpr:
			switch xType := t.X.(type) {
			case *ast.Ident:
				if argument.Type, err = ctx.GetRefType(xType.Name); err != nil {
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
