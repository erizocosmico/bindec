package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/erizocosmico/bindec"
)

func main() {
	var fs flag.FlagSet
	var recv, path, typ, output string
	fs.StringVar(&recv, "recv", "t", "Name given to the receiver type on the generated methods. For multiple types, separate with commas e.g. -recv=t,x,c.")
	fs.StringVar(&typ, "type", "", "Type/s to generate encoder and decoder for. Separate with commas for more than one e.g. -type=A,B,C.")
	fs.StringVar(&output, "o", "", "Generated file name, by default TYPE_bindec.go.")
	fs.Parse(os.Args[1:])

	if typ == "" {
		fmt.Println("-type argument is mandatory!")
		os.Exit(1)
	}

	args := fs.Args()
	if l := len(args); l > 1 {
		fs.Usage()
		os.Exit(1)
	} else if l == 1 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		assert(err)
	}

	recvs := strings.Split(recv, ",")
	types := strings.Split(typ, ",")

	if len(recvs) == 0 || recv == "t" {
		recvs = make([]string, len(types))
		for i := range types {
			recvs[i] = "t"
		}
	}

	if len(recvs) != len(types) {
		assert(fmt.Errorf(
			"got %d receivers, but %d types to generate",
			len(recvs), len(types),
		))
	}

	filename := strings.ToLower(strings.Join(types, "_")) + "_bindec.go"
	content, err := bindec.Generate(bindec.Options{
		Path:  path,
		Recvs: recvs,
		Types: types,
	})
	assert(err)

	file := output
	if file == "" {
		file = filepath.Join(path, filename)
	}

	f, err := os.Create(file)
	assert(err)

	_, err = f.Write(content)
	assert(err)
	assert(f.Close())
}

func assert(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
