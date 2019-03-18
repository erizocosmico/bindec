package bench

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

var input = Foo{
	A: 1,
	B: "foooo",
	C: []byte("baaar"),
	D: struct {
		A uint64
		B string
	}{
		A: 6,
		B: "baaaz",
	},
	E: []int{1, 2, 3, 4},
	F: [2]int{5, 6},
	G: true,
}

func init() {
	gob.Register(Foo{})
}

func BenchmarkEncode(b *testing.B) {
	b.Run("bindec", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := input.EncodeBinary()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gob", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			err := gob.NewEncoder(&buf).Encode(input)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	bindecEncoded, err := input.EncodeBinary()
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(input)
	if err != nil {
		b.Fatal(err)
	}
	gobEncoded := buf.Bytes()

	b.Run("bindec", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var out Foo
			err := out.DecodeBinaryFromBytes(bindecEncoded)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("gob", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var out Foo
			err := gob.NewDecoder(bytes.NewReader(gobEncoded)).Decode(&out)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestSize(t *testing.T) {
	bindecEncoded, err := input.EncodeBinary()
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(input)
	if err != nil {
		t.Fatal(err)
	}
	gobEncoded := buf.Bytes()

	fmt.Println("BINDEC:", len(bindecEncoded), "bytes")
	fmt.Println("GOB:", len(gobEncoded), "bytes")

	if len(bindecEncoded) >= len(gobEncoded) {
		t.Errorf(
			"expected bindec (%d bytes) to be smaller than gob (%d bytes)",
			len(bindecEncoded),
			len(gobEncoded),
		)
	}
}
