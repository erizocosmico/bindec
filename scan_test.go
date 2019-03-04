package bindec

import (
	"go/types"
	"testing"
)

func TestGetPackage(t *testing.T) {
	pkg, err := getPackage("")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Name() != "bindec" {
		t.Errorf("expected package name to be bindec, is %q", pkg.Name())
	}

	pkg, err = getPackage("./cmd/bindec")
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Name() != "main" {
		t.Errorf("expected package name to be main, is %q", pkg.Name())
	}
}

func TestFindType(t *testing.T) {
	pkg, err := getPackage("")
	if err != nil {
		t.Fatal(err)
	}

	typ, err := findType(pkg, "Options")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := typ.(*types.Named)
	if !ok {
		t.Errorf("expecting Options to be *types.Named, is %T", typ)
	}
}
