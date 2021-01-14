package myasthurts

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path"
)

// ParsePackageContext keeps all information needed for parsing a package.
type ParsePackageContext struct {
	// BuildPackage holds the build information about the package.
	BuildPackage *build.Package

	// Package is what is being parsed. It will keep all information for what is
	// being parsed.
	Package *Package
}

func NewPackageContext(pkg *Package, buildPackage *build.Package) *ParsePackageContext {
	return &ParsePackageContext{
		BuildPackage: buildPackage,
		Package:      pkg,
	}
}

// ParseFileContext will keep all information needed for parsing a single file.
type ParseFileContext struct {
	File                  *ast.File
	FSet                  *token.FileSet
	Env                   *Environment
	Package               *Package
	dotImports            []*Package
	packageImportAliasMap map[string]*Package
}

func (ctx *ParseFileContext) PackageByImportAlias(name string) (*Package, bool) {
	pkg, ok := ctx.packageImportAliasMap[name]
	return pkg, ok
}

// GetRefType will return a type defined on the context or in the dot imported
// libraries. If no file exists, it will return an `ErrTypeNotFound`.
func (ctx *ParseFileContext) GetRefType(name string) (RefType, bool) {
	// First, it tries to find the type on its own package.
	if t, ok := ctx.Package.RefTypeByName(name); ok {
		return t, true
	}

	// Now, it tries to find into the dot imported libraries.
	for _, pkg := range ctx.dotImports {
		if t, ok := pkg.RefTypeByName(name); ok {
			return t, true
		}
	}

	// Givin' up, never so easy...
	return nil, false
}

// Environment is the virtual representation of a Go environment.
type Environment struct {
	// BuildContext is the reference of the build context used to extract
	// information about the packages.
	BuildContext build.Context

	// BuiltIn is the builtin package reference, already explored.
	BuiltIn *Package

	// packages is the list of packages inside of this environment.
	packages []*Package

	// packageMap stores the *Packages reference by its import name.
	packageMap map[string]*Package

	Config EnvConfig

	// Listener implement a set of interfaces that let the developer to have
	// some control over how files are parsed.
	Listener interface{}
}

func NewEnvironment() (*Environment, error) {
	env := &Environment{
		packages:     make([]*Package, 0, 5),
		packageMap:   make(map[string]*Package, 5),
		BuildContext: build.Default,
	}

	if err := env.makeEnv(); err != nil {
		return nil, err
	}
	return env, nil
}

func NewEnvironmentWithListener(listener interface{}) (*Environment, error) {
	env, err := NewEnvironment()
	if err != nil {
		return nil, err
	}
	env.Listener = listener
	return env, nil
}

func (env *Environment) Import(importPathPkg string) (*build.Package, error) {
	d := "."
	if env.Config.CurrentDir != "" {
		d = env.Config.CurrentDir
	}
	buildPkg, err := env.BuildContext.Import(importPathPkg, d, build.ImportComment)
	if err != nil {
		return nil, err
	}
	return buildPkg, nil
}

func (env *Environment) ImportDir(importDir string) (*build.Package, error) {
	buildPkg, err := env.BuildContext.ImportDir(importDir, build.ImportComment)
	if err != nil {
		return nil, err
	}
	return buildPkg, nil
}

// PackageByImportPath find Package by name in Environment.
func (env *Environment) PackageByImportPath(importPath string) (*Package, bool) {
	pkg, ok := env.packageMap[importPath]
	return pkg, ok
}

// AppendPackage add new Package in Environment.
func (env *Environment) AppendPackage(pkg *Package) {
	env.packages = append(env.packages, pkg)
	env.packageMap[pkg.ImportPath] = pkg
}

// parsePackage will list all files for a package and
func (env *Environment) parsePackage(pkgCtx *ParsePackageContext) error {
	for _, file := range pkgCtx.BuildPackage.GoFiles {
		filePath := path.Join(pkgCtx.Package.RealPath, file)

		if beforeFileListener, ok := env.Listener.(ListenerBeforeFile); ok {
			err := beforeFileListener.BeforeFile(pkgCtx, filePath)
			if err == Skip { // Shall the file be skipped?
				continue
			} else if err != nil { // This is an actual error...
				return err
			}
		}

		err := env.ParseFile(pkgCtx, filePath)
		if fileListener, ok := env.Listener.(ListenerAfterFile); ok {
			errAfterFile := fileListener.AfterFile(pkgCtx, filePath, err)
			if err != nil {
				return errAfterFile
			}
		}
		if err != nil {
			return err
		}
	}
	pkgCtx.Package.Explored = true
	return nil
}

func (env *Environment) ParseDir(dir string) (*Package, error) {
	// Find the path of the package.
	buildPkg, err := env.BuildContext.ImportDir(dir, build.ImportComment)
	if err != nil {
		return nil, err
	}

	p, ok := env.packageMap[buildPkg.ImportPath]
	if ok { // If the package exists in the environment.
		if p.Explored { // If the package is already explored.
			return p, nil // just return it, no need to do anything.
		}
		// The package will be explored down function.
	} else {
		p = NewPackage(buildPkg)
	}

	err = env.parsePackage(NewPackageContext(p, buildPkg))
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Parse checks if the parse was already done, if not, it parses the package.
func (env *Environment) Parse(packageName string) (*Package, error) {
	p, ok := env.packageMap[packageName]
	if ok && p.Explored { // If the package exists in the environment and it was explored.
		return p, nil // just return it, no need to do anything.
	}

	ctx := &env.BuildContext

	// Find the path of the package.
	buildPkg, err := ctx.Import(packageName, ".", build.ImportComment)
	if err != nil {
		return nil, err
	}

	newPkg := p
	if newPkg == nil {
		newPkg = NewPackage(buildPkg)
	}

	pkgCtx := NewPackageContext(newPkg, buildPkg)
	if err = env.parsePackage(pkgCtx); err != nil {
		return nil, err
	}

	if p == nil { // If it was not defined before
		env.AppendPackage(pkgCtx.Package) // define it now
	}

	return pkgCtx.Package, nil
}

func (env *Environment) gorootSourceDir() (rtn string, exrr error) {
	if rtn = os.Getenv("GOROOT"); rtn == "" {
		return "", errors.New("GOROOT environment variable not found or is empty")
	}
	return fmt.Sprintf("%s/src", rtn), nil
}

func (env *Environment) makeEnv() error {
	pkg, err := env.Parse("builtin")
	if err != nil {
		return err
	}
	env.BuiltIn = pkg
	return nil
}

func (env *Environment) ParseFile(pkgCtx *ParsePackageContext, filePath string) error {
	var (
		file *ast.File
		fset *token.FileSet
		err  error
	)

	fset = token.NewFileSet()
	if file, err = parser.ParseFile(fset, filePath, nil, parser.ParseComments); err != nil {
		return err
	}

	dotImports := make([]*Package, 0, 1)
	if env.BuiltIn != nil {
		// Adds the built in as a default dot imported package. If it is defined.
		dotImports = append(dotImports, env.BuiltIn)
	} // If it is not defined, it means we are parsing the builtin package.

	// Create the context of the file parse.
	fileCtx := &ParseFileContext{
		File:                  file,
		FSet:                  fset,
		Env:                   env,
		Package:               pkgCtx.Package,
		dotImports:            dotImports,
		packageImportAliasMap: make(map[string]*Package),
	}

	// Prints the AST, if configured.
	if env.Config.DevMode && env.Config.ASTI {
		ast.Print(fset, file)
	}

	err = parseFileName(fileCtx)
	if err != nil {
		return err
	}

	decls := file.Decls
	for _, d := range decls {
		switch c := d.(type) {
		case *ast.GenDecl:
			err = parseGenDecl(fileCtx, c)
			if err != nil {
				return err
			}
		case *ast.FuncDecl:
			err = parseFuncDecl(fileCtx, c)
			if err != nil {
				return err
			}

		default:
			return errors.New("Decl not found")
		}

	}
	return nil
}
