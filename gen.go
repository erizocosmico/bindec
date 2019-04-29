package bindec

import (
	"fmt"
	"go/format"
	"strings"
)

// Options to configure generation.
type Options struct {
	// Path of the package in which the type is located.
	Path string
	// Types to generate encoder and decoder for.
	Types []string
	// Recvs are the receiver names for the generated methods.
	Recvs []string
}

// Generate a file of source code containing an encoder and a decoder to
// encode and decode a given type to and from a binary representation of
// itself.
func Generate(opts Options) ([]byte, error) {
	pkg, err := getPackage(opts.Path)
	if err != nil {
		return nil, err
	}

	ctx := newParseContext()
	ctx.addImport("encoding/binary")
	ctx.addImport("bytes")
	ctx.addImport("io")
	ctx.addImport("math")

	var methods = make([]string, len(opts.Types))
	for i, tName := range opts.Types {
		recv := opts.Recvs[i]

		typ, err := findType(pkg, tName)
		if err != nil {
			return nil, err
		}

		t, err := parseType(ctx, typ)
		if err != nil {
			return nil, err
		}

		methods[i] = generateMethods(recv, tName, t)
	}

	src := []byte(generateFile(
		pkg.Name(),
		strings.Join(methods, "\n"),
		ctx.getImports(),
		ctx.getDecls(),
	))

	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("error formatting code: %s\n\n%s", err, prettySource(src))
	}

	return formatted, nil
}

func generateMethods(recv, typeName string, typ Type) string {
	return fmt.Sprintf(
		methodsTpl,
		recv,
		typeName,
		typ.Encoder(recv),
		typ.Decoder(recv, true),
	)
}

func generateFile(
	pkgName string,
	methods string,
	imports []string,
	decls []string,
) string {
	var deps = make([]string, len(imports))
	for i, x := range imports {
		deps[i] = fmt.Sprintf("%q", x)
	}

	return fmt.Sprintf(
		fileTpl,
		pkgName,
		methods,
		strings.Join(deps, "\n"),
		strings.Join(decls, "\n"),
	)
}

func prettySource(src []byte) string {
	lines := strings.Split(string(src), "\n")
	maxDigits := len(fmt.Sprint(len(lines)))
	format := fmt.Sprintf("%%%dd | %%s", maxDigits)

	for i, line := range lines {
		lines[i] = fmt.Sprintf(format, i+1, line)
	}

	return strings.Join(lines, "\n")
}
