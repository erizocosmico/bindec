package bindec

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"
)

func getPackage(path string) (*types.Package, error) {
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var files []*ast.File
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}

		for _, f := range pkg.Files {
			files = append(files, f)
		}
		break
	}

	return typeCheck(path, fset, files)
}

func typeCheck(path string, fset *token.FileSet, files []*ast.File) (*types.Package, error) {
	cfg := &types.Config{
		IgnoreFuncBodies:         true,
		FakeImportC:              true,
		DisableUnusedImportCheck: true,
		Importer:                 importer.Default(),
	}

	return cfg.Check(path, fset, files, nil)
}

func findType(pkg *types.Package, typ string) (types.Type, error) {
	obj := pkg.Scope().Lookup(typ)
	if obj == nil {
		return nil, fmt.Errorf("type %s not found in %s", typ, pkg.Path())
	}
	return obj.Type(), nil
}
