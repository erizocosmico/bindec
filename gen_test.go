package bindec

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	path, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	data, err := Generate(Options{Path: path, Type: "StructTestType", Recv: "t"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "file.go", string(data)+typeTestDefs, 0)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	cfg := &types.Config{
		FakeImportC:              true,
		DisableUnusedImportCheck: true,
		Importer:                 importer.For("source", nil),
	}

	_, err = cfg.Check(path, fset, []*ast.File{file}, nil)
	if err != nil {
		t.Errorf("expected generated file to type check, got: %s", err)
	}
}

func TestGenerateCyclic(t *testing.T) {
	path, err := filepath.Abs(".")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = Generate(Options{Path: path, Type: "StructCyclic", Recv: "t"})
	if err == nil {
		t.Errorf("expected error")
	}
}
