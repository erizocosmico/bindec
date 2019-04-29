package bindec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type encoderDecoder interface {
	encoder
	decoder
}

func TestAlpha(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"skjhfdkjfhsdkjf", true},
		{"sdfsdfs.", false},
		{"slkfjslkfjdss3fdksljf", false},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &AlphaTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestAlphanum(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"skjhfdkjfhsdkjf", true},
		{"sdfsdfs.", false},
		{"sdfs3dfs.", false},
		{"slkfjslkfjdss3fdksljf", true},
		{"2189372137", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &AlphanumTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestNumeric(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"12345.6", false},
		{"12xd.", false},
		{"0xbeef.", false},
		{"32490823", true},
		{"1", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &NumericTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestHexadecimal(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", true},
		{"1239801.", false},
		{"beef", true},
		{"32490823", true},
		{"avcdejkh", false},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &HexadecimalTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestEmail(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", false},
		{"1239801.", false},
		{"foo@foo", false},
		{"foo@foo.bar", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &EmailTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestURL(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", false},
		{"1239801.", false},
		{"example.com/foo.bar", false},
		{"https://example.com/foo.bar", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &URLTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestBase64(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"ljkfdasdja9039482", false},
		{"Zm9vYmFy", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &Base64TestType{S: tt.str}, tt.ok)
		})
	}
}

func TestIPv4(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", false},
		{"1239801.", false},
		{"foo@foo", false},
		{"127.0.0", false},
		{"127.0.0.1", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &IPv4TestType{S: tt.str}, tt.ok)
		})
	}
}

func TestIPv6(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", false},
		{"1239801.", false},
		{"foo@foo", false},
		{"127.0.0", false},
		{"127.0.0.1", false},
		{"2001:0db8:85a3:08d3:1319:8a2e:0370", false},
		{"2001:0db8:85a3:08d3:1319:8a2e:0370:7334", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &IPv6TestType{S: tt.str}, tt.ok)
		})
	}
}

func TestIP(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", false},
		{"1239801.", false},
		{"foo@foo", false},
		{"127.0.0", false},
		{"127.0.0.1", true},
		{"2001:0db8:85a3:08d3:1319:8a2e:0370", false},
		{"2001:0db8:85a3:08d3:1319:8a2e:0370:7334", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &IPTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestUUID(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"beefad", false},
		{"1239801.", false},
		{"8418c509-c51b-4afd-9c6d-adffb010", false},
		{"8418c509-c51b-4afd-9c6d-adffb010e6d3", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &UUIDTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"hello", false},
		{"hello world.", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &ContainsTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestStartsWith(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"hello", false},
		{"Hello world.", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &StartsWithTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestEndsWith(t *testing.T) {
	testCases := []struct {
		str string
		ok  bool
	}{
		{"hello", false},
		{"hello world", true},
	}

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			assertConstraints(t, &EndsWithTestType{S: tt.str}, tt.ok)
		})
	}
}

func TestMinLen(t *testing.T) {
	assertConstraints(t, &MinLenTestType{
		String: "",
		Bytes:  []byte{1, 2, 3, 4, 5},
		Slice:  []int{1, 2, 3, 4, 5},
	}, false)

	assertConstraints(t, &MinLenTestType{
		String: "12345",
		Bytes:  []byte{1, 2, 3, 4},
		Slice:  []int{1, 2, 3, 4, 5},
	}, false)

	assertConstraints(t, &MinLenTestType{
		String: "12345",
		Bytes:  []byte{1, 2, 3, 4, 5},
		Slice:  []int{1, 2, 3, 4},
	}, false)

	assertConstraints(t, &MinLenTestType{
		String: "12345",
		Bytes:  []byte{1, 2, 3, 4, 5},
		Slice:  []int{1, 2, 3, 4, 5},
	}, true)
}

func TestMaxLen(t *testing.T) {
	assertConstraints(t, &MaxLenTestType{
		String: "123456",
		Bytes:  []byte{1, 2, 3, 4, 5},
		Slice:  []int{1, 2, 3, 4, 5},
	}, false)

	assertConstraints(t, &MaxLenTestType{
		String: "12345",
		Bytes:  []byte{1, 2, 3, 4, 5, 6},
		Slice:  []int{1, 2, 3, 4, 5},
	}, false)

	assertConstraints(t, &MaxLenTestType{
		String: "12345",
		Bytes:  []byte{1, 2, 3, 4, 5},
		Slice:  []int{1, 2, 3, 4, 5, 6},
	}, false)

	assertConstraints(t, &MaxLenTestType{
		String: "12345",
		Bytes:  []byte{1, 2, 3, 4, 5},
		Slice:  []int{1, 2, 3, 4, 5},
	}, true)
}

func TestOneOf(t *testing.T) {
	assertConstraints(t, &OneOfTestType{
		Uint8:   1,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    1,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  1,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   1,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  1,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   3,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   1,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  1,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   1,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   6,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    1,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   6,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     1,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   3,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 1,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hallo",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    false,
		Float32: 3.14,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.,
		Float64: 1.1,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.,
	}, false)

	assertConstraints(t, &OneOfTestType{
		Uint8:   2,
		Int8:    2,
		Uint16:  3,
		Int16:   6,
		Uint32:  6,
		Int32:   3,
		Uint64:  3,
		Int64:   2,
		Uint:    6,
		Int:     6,
		Uintptr: 2,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 1.1,
	}, true)
}

func TestEq(t *testing.T) {
	assertConstraints(t, &EqTestType{
		Uint8:   0,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    0,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  0,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   0,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  0,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   0,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  0,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   0,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    0,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     0,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 0,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "foo",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    false,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.,
	}, false)

	assertConstraints(t, &EqTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		String:  "hello",
		Bool:    true,
		Float32: 3.14,
		Float64: 3.14,
	}, true)
}
func TestNeq(t *testing.T) {
	assertConstraints(t, &NeqTestType{
		Uint8:   6,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    6,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  6,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   6,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  6,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   6,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  6,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   6,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    6,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     6,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 6,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "hello",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    true,
		Float32: 2.,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 3.14,
		Float64: 2.,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &NeqTestType{
		Uint8:   0,
		Int8:    0,
		Uint16:  0,
		Int16:   0,
		Uint32:  0,
		Int32:   0,
		Uint64:  0,
		Int64:   0,
		Uint:    0,
		Int:     0,
		Uintptr: 0,
		String:  "bar",
		Bool:    false,
		Float32: 2.,
		Float64: 2.,
	}, true)
}

func TestMin(t *testing.T) {
	assertConstraints(t, &MinTestType{
		Uint8:   0,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    0,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  0,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   0,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  0,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   0,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  0,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   0,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    0,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     0,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 0,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.,
	}, false)

	assertConstraints(t, &MinTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, true)
}

func TestMax(t *testing.T) {
	assertConstraints(t, &MaxTestType{
		Uint8:   7,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    7,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  7,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   7,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  7,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   7,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  7,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   7,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    7,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     7,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 7,
		Float32: 3.14,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 4.,
		Float64: 3.14,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 4.,
	}, false)

	assertConstraints(t, &MaxTestType{
		Uint8:   6,
		Int8:    6,
		Uint16:  6,
		Int16:   6,
		Uint32:  6,
		Int32:   6,
		Uint64:  6,
		Int64:   6,
		Uint:    6,
		Int:     6,
		Uintptr: 6,
		Float32: 3.14,
		Float64: 3.14,
	}, true)
}

func assertConstraints(t *testing.T, in encoderDecoder, ok bool) {
	t.Helper()

	bs, err := in.EncodeBinary()
	require.NoError(t, err)

	err = in.DecodeBinaryFromBytes(bs)
	if ok {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
	}
}
