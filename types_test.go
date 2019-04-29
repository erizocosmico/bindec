package bindec

//go:generate ./bindec_bin -type=StructTestType,MapTestType,ArrayTestType,SliceTestType,ByteTestType,Uint16TestType,Uint32TestType,Uint64TestType,UintTestType,Int8TestType,Int16TestType,Int32TestType,Int64TestType,IntTestType,UintptrTestType,Float32TestType,Float64TestType,StringTestType,BytesTestType,BoolTestType,AlphaTestType,AlphanumTestType,NumericTestType,HexadecimalTestType,EmailTestType,URLTestType,Base64TestType,ContainsTestType,StartsWithTestType,EndsWithTestType,EqTestType,NeqTestType,UUIDTestType,IPTestType,IPv4TestType,IPv6TestType,OneOfTestType,MaxTestType,MinTestType,MaxLenTestType,MinLenTestType -o bindec_test.go

type (
	MapTestType   map[byte]uint16
	SliceTestType []uint16
	ArrayTestType [2]byte

	ByteTestType    byte
	Uint16TestType  uint16
	Uint32TestType  uint32
	Uint64TestType  uint64
	UintptrTestType uintptr
	UintTestType    uint

	Int8TestType  int8
	Int16TestType int16
	Int32TestType int32
	Int64TestType int64
	IntTestType   int

	Float32TestType float32
	Float64TestType float64

	StringTestType string
	BytesTestType  []byte

	BoolTestType bool
)

type Struct2 struct {
	Field1  int
	Flield2 string
}

type StructTestType struct {
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Int        int
	Byte       byte
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uint       uint
	String     string
	Float32    float32
	Float64    float64
	Bool       bool
	Pointer    *bool
	NilPointer *bool
	Slice      []int16
	Bytes      []byte
	Array      [4]int16
	Struct     struct {
		Field1  int
		Flield2 string
	}
	NamedStruct   Struct2
	StructPointer *Struct2
	Ignored       int `bindec:"-"`
}

const typeTestDefs = `
type Struct2 struct {
	Field1  int
	Flield2 string
}

type StructTestType struct {
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Int        int
	Byte       byte
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uint       uint
	String     string
	Float32    float32
	Float64    float64
	Bool       bool
	Pointer    *bool
	NilPointer *bool
	Slice      []int16
	Bytes      []byte
	Array      [4]int16
	Struct     struct {
		Field1  int
		Flield2 string
	}
	NamedStruct   Struct2
	StructPointer *Struct2
	Ignored int ` + "`bindec:\"-\"`" + `
}
`

type StructCyclic struct {
	Cycle *StructCyclic
}

type AlphaTestType struct {
	S string `bindec:"alpha"`
}

type AlphanumTestType struct {
	S string `bindec:"alphanum"`
}

type NumericTestType struct {
	S string `bindec:"numeric"`
}

type HexadecimalTestType struct {
	S string `bindec:"hexadecimal"`
}

type EmailTestType struct {
	S string `bindec:"email"`
}

type URLTestType struct {
	S string `bindec:"url"`
}

type Base64TestType struct {
	S string `bindec:"base64"`
}

type ContainsTestType struct {
	S string `bindec:"contains=worl"`
}

type StartsWithTestType struct {
	S string `bindec:"startswith=Hello"`
}

type EndsWithTestType struct {
	S string `bindec:"endswith=world"`
}

type IPv6TestType struct {
	S string `bindec:"ipv6"`
}

type IPv4TestType struct {
	S string `bindec:"ipv4"`
}

type IPTestType struct {
	S string `bindec:"ip"`
}

type UUIDTestType struct {
	S string `bindec:"uuid"`
}

type EqTestType struct {
	Uint8   uint8   `bindec:"eq=6"`
	Int8    int8    `bindec:"eq=6"`
	Uint16  uint16  `bindec:"eq=6"`
	Int16   int16   `bindec:"eq=6"`
	Uint32  uint32  `bindec:"eq=6"`
	Int32   int32   `bindec:"eq=6"`
	Uint64  uint64  `bindec:"eq=6"`
	Int64   int64   `bindec:"eq=6"`
	Uint    uint    `bindec:"eq=6"`
	Int     int     `bindec:"eq=6"`
	Uintptr uintptr `bindec:"eq=6"`
	String  string  `bindec:"eq=hello"`
	Bool    bool    `bindec:"eq=true"`
	Float32 float32 `bindec:"eq=3.14"`
	Float64 float64 `bindec:"eq=3.14"`
}

type NeqTestType struct {
	Uint8   uint8   `bindec:"neq=6"`
	Int8    int8    `bindec:"neq=6"`
	Uint16  uint16  `bindec:"neq=6"`
	Int16   int16   `bindec:"neq=6"`
	Uint32  uint32  `bindec:"neq=6"`
	Int32   int32   `bindec:"neq=6"`
	Uint64  uint64  `bindec:"neq=6"`
	Int64   int64   `bindec:"neq=6"`
	Uint    uint    `bindec:"neq=6"`
	Int     int     `bindec:"neq=6"`
	Uintptr uintptr `bindec:"neq=6"`
	String  string  `bindec:"neq=hello"`
	Bool    bool    `bindec:"neq=true"`
	Float32 float32 `bindec:"neq=3.14"`
	Float64 float64 `bindec:"neq=3.14"`
}

type MinTestType struct {
	Uint8   uint8   `bindec:"min=6"`
	Int8    int8    `bindec:"min=6"`
	Uint16  uint16  `bindec:"min=6"`
	Int16   int16   `bindec:"min=6"`
	Uint32  uint32  `bindec:"min=6"`
	Int32   int32   `bindec:"min=6"`
	Uint64  uint64  `bindec:"min=6"`
	Int64   int64   `bindec:"min=6"`
	Uint    uint    `bindec:"min=6"`
	Int     int     `bindec:"min=6"`
	Uintptr uintptr `bindec:"min=6"`
	Float32 float32 `bindec:"min=3.14"`
	Float64 float64 `bindec:"min=3.14"`
}

type MaxTestType struct {
	Uint8   uint8   `bindec:"max=6"`
	Int8    int8    `bindec:"max=6"`
	Uint16  uint16  `bindec:"max=6"`
	Int16   int16   `bindec:"max=6"`
	Uint32  uint32  `bindec:"max=6"`
	Int32   int32   `bindec:"max=6"`
	Uint64  uint64  `bindec:"max=6"`
	Int64   int64   `bindec:"max=6"`
	Uint    uint    `bindec:"max=6"`
	Int     int     `bindec:"max=6"`
	Uintptr uintptr `bindec:"max=6"`
	Float32 float32 `bindec:"max=3.14"`
	Float64 float64 `bindec:"max=3.14"`
}

type MinLenTestType struct {
	String string `bindec:"minlen=5"`
	Bytes  []byte `bindec:"minlen=5"`
	Slice  []int  `bindec:"minlen=5"`
}

type MaxLenTestType struct {
	String string `bindec:"maxlen=5"`
	Bytes  []byte `bindec:"maxlen=5"`
	Slice  []int  `bindec:"maxlen=5"`
}

type OneOfTestType struct {
	Uint8   uint8   `bindec:"oneof=6 2 3"`
	Int8    int8    `bindec:"oneof=6 2 3"`
	Uint16  uint16  `bindec:"oneof=6 2 3"`
	Int16   int16   `bindec:"oneof=6 2 3"`
	Uint32  uint32  `bindec:"oneof=6 2 3"`
	Int32   int32   `bindec:"oneof=6 2 3"`
	Uint64  uint64  `bindec:"oneof=6 2 3"`
	Int64   int64   `bindec:"oneof=6 2 3"`
	Uint    uint    `bindec:"oneof=6 2 3"`
	Int     int     `bindec:"oneof=6 2 3"`
	Uintptr uintptr `bindec:"oneof=6 2 3"`
	String  string  `bindec:"oneof=hello world foo"`
	Bool    bool    `bindec:"oneof=true"`
	Float32 float32 `bindec:"oneof=3.14 1.1 2.2"`
	Float64 float64 `bindec:"oneof=3.14 1.1 2.2"`
}
