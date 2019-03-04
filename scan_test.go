package bindec

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPackage(t *testing.T) {
	require := require.New(t)
	pkg, err := getPackage("")
	require.NoError(err)
	require.Equal("bindec", pkg.Name())

	pkg, err = getPackage("./cmd/bindec")
	require.NoError(err)
	require.Equal("main", pkg.Name())
}

func TestFindType(t *testing.T) {
	require := require.New(t)
	pkg, err := getPackage("")
	require.NoError(err)

	typ, err := findType(pkg, "Options")
	require.NoError(err)
	_, ok := typ.(*types.Named)
	require.True(ok)
}
