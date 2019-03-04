package bench

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		require := require.New(b)
		for i := 0; i < b.N; i++ {
			_, err := input.EncodeBinary()
			require.NoError(err)
		}
	})

	b.Run("gob", func(b *testing.B) {
		require := require.New(b)
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			err := gob.NewEncoder(&buf).Encode(input)
			require.NoError(err)
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	bindecEncoded, err := input.EncodeBinary()
	require.NoError(b, err)

	var buf bytes.Buffer
	require.NoError(b, gob.NewEncoder(&buf).Encode(input))
	gobEncoded := buf.Bytes()

	b.Run("bindec", func(b *testing.B) {
		require := require.New(b)
		for i := 0; i < b.N; i++ {
			var out Foo
			err := out.DecodeBinary(bindecEncoded)
			require.NoError(err)
		}
	})

	b.Run("gob", func(b *testing.B) {
		require := require.New(b)
		for i := 0; i < b.N; i++ {
			var out Foo
			err := gob.NewDecoder(bytes.NewReader(gobEncoded)).Decode(&out)
			require.NoError(err)
		}
	})
}

func TestSize(t *testing.T) {
	bindecEncoded, err := input.EncodeBinary()
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, gob.NewEncoder(&buf).Encode(input))
	gobEncoded := buf.Bytes()

	fmt.Println("BINDEC:", len(bindecEncoded), "bytes")
	fmt.Println("GOB:", len(gobEncoded), "bytes")

	require.True(t, len(bindecEncoded) < len(gobEncoded))
}
