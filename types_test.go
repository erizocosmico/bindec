package bindec

//go:generate ./bindec_bin -type=StructTestType,MapTestType,ArrayTestType,SliceTestType,ByteTestType,Uint16TestType,Uint32TestType,Uint64TestType,UintTestType,Int8TestType,Int16TestType,Int32TestType,Int64TestType,IntTestType,UintptrTestType,Float32TestType,Float64TestType,StringTestType,BytesTestType,BoolTestType -o bindec_test.go

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
