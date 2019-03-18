package bindec

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
}
`

type StructCyclic struct {
	Cycle *StructCyclic
}
