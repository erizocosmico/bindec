package bindec

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	trueVal := true
	input := StructTestType{
		Int8:       math.MaxInt8,
		Int16:      math.MaxInt16,
		Int32:      math.MaxInt32,
		Int64:      math.MaxInt64,
		Int:        math.MaxInt32,
		Byte:       math.MaxUint8,
		Uint8:      math.MaxUint8,
		Uint16:     math.MaxUint16,
		Uint32:     math.MaxUint32,
		Uint64:     math.MaxUint64,
		Uint:       math.MaxUint32,
		String:     "cowabunga",
		Float32:    3.123201938,
		Float64:    6.329840298429048,
		Bool:       true,
		Pointer:    &trueVal,
		NilPointer: nil,
		Slice:      []int16{0, 1, -5, -math.MaxInt16, math.MaxInt16},
		Bytes:      []byte("ay caramba"),
		Array:      [4]int16{0, 1, -5, 4},
		Struct: struct {
			Field1  int
			Flield2 string
		}{6, "foo"},
		NamedStruct:   Struct2{4, "bar"},
		StructPointer: &Struct2{8, "baz"},
		Ignored:       42,
	}

	output, err := input.EncodeBinary()
	if err != nil {
		t.Errorf("unexpected error encoding: %s", err)
	}

	var result StructTestType
	if err := result.DecodeBinaryFromBytes(output); err != nil {
		t.Errorf("unexpected error decoding: %s", err)
	}

	if result.Ignored != 0 {
		t.Errorf("expected Ignored field to be 0")
	}

	input.Ignored = 0 // just to be able to equal them
	require.Equal(t, input, result)
}
