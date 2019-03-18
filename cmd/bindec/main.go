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
	fs.StringVar(&recv, "recv", "t", "Name given to the receiver type on the generated methods.")
	fs.StringVar(&typ, "type", "", "Type to generate encoder and decoder for.")
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

	filename := strings.ToLower(typ) + "_bindec.go"
	content, err := bindec.Generate(bindec.Options{
		Path: path,
		Recv: recv,
		Type: typ,
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
