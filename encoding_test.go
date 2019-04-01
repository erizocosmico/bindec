package bindec

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructEncodeDecode(t *testing.T) {
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

type encoder interface {
	EncodeBinary() ([]byte, error)
}

type decoder interface {
	DecodeBinaryFromBytes([]byte) error
}

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name   string
		input  encoder
		output decoder
	}{
		{
			"byte",
			ByteTestType(2),
			func() decoder {
				var v ByteTestType
				return &v
			}(),
		},
		{
			"uint16",
			Uint16TestType(2),
			func() decoder {
				var v Uint16TestType
				return &v
			}(),
		},
		{
			"uint32",
			Uint32TestType(2),
			func() decoder {
				var v Uint32TestType
				return &v
			}(),
		},
		{
			"uint64",
			Uint64TestType(2),
			func() decoder {
				var v Uint64TestType
				return &v
			}(),
		},
		{
			"uintptr",
			UintptrTestType(2),
			func() decoder {
				var v UintptrTestType
				return &v
			}(),
		},
		{
			"uint",
			UintTestType(2),
			func() decoder {
				var v UintTestType
				return &v
			}(),
		},
		{
			"int8",
			Int8TestType(2),
			func() decoder {
				var v Int8TestType
				return &v
			}(),
		},
		{
			"int16",
			Int16TestType(2),
			func() decoder {
				var v Int16TestType
				return &v
			}(),
		},
		{
			"int32",
			Int32TestType(2),
			func() decoder {
				var v Int32TestType
				return &v
			}(),
		},
		{
			"int64",
			Int64TestType(2),
			func() decoder {
				var v Int64TestType
				return &v
			}(),
		},
		{
			"int",
			IntTestType(2),
			func() decoder {
				var v IntTestType
				return &v
			}(),
		},
		{
			"float32",
			Float32TestType(3.14),
			func() decoder {
				var v Float32TestType
				return &v
			}(),
		},
		{
			"float64",
			Float64TestType(3.14),
			func() decoder {
				var v Float64TestType
				return &v
			}(),
		},
		{
			"bool",
			BoolTestType(true),
			func() decoder {
				var v BoolTestType
				return &v
			}(),
		},
		{
			"string",
			StringTestType("hello world"),
			func() decoder {
				var v StringTestType
				return &v
			}(),
		},
		{
			"bytes",
			BytesTestType{1, 2, 3},
			&BytesTestType{},
		},
		{
			"slice",
			SliceTestType{1, 2, 3},
			&SliceTestType{},
		},
		{
			"array",
			ArrayTestType{4, 2},
			&ArrayTestType{},
		},
		{
			"map",
			MapTestType{
				1: 2,
				3: 4,
			},
			&MapTestType{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			output, err := tt.input.EncodeBinary()
			require.NoError(err)

			require.NoError(tt.output.DecodeBinaryFromBytes(output))

			result := reflect.ValueOf(tt.output).Elem().Interface()
			require.Equal(tt.input, result)
		})
	}
}
