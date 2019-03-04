package bindec

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
	"strings"
)

// Type can generate decoders and encoders for a given type.
type Type interface {
	// Decoder generates a decoder for the type. recv is the variable or
	// struct field the data will be decoded into.
	Decoder(recv string) string
	// Encoder generates an encoder for the type. recv is the variable or
	// struct field that will be encoded.
	Encoder(recv string) string
}

// BasicKind is the kind of basic type.
type BasicKind = types.BasicKind

const (
	// String type.
	String = types.String
	// Bool type.
	Bool = types.Bool
	// Int type.
	Int = types.Int
	// Int8 type.
	Int8 = types.Int8
	// Int16 type.
	Int16 = types.Int16
	// Int32 type.
	Int32 = types.Int32
	// Int64 type.
	Int64 = types.Int64
	// Uint type.
	Uint = types.Uint
	// Uint8 type.
	Uint8 = types.Uint8
	// Uint16 type.
	Uint16 = types.Uint16
	// Uint32 type.
	Uint32 = types.Uint32
	// Uint64 type.
	Uint64 = types.Uint64
	// Uintptr type.
	Uintptr = types.Uintptr
	// Float32 type.
	Float32 = types.Float32
	// Float64 type.
	Float64 = types.Float64
)

// Basic type.
type Basic struct {
	Kind BasicKind
}

// Encoder implements the Type interface.
func (t Basic) Encoder(recv string) string {
	return writeBasic(recv, t.Kind)
}

// Decoder implements the Type interface.
func (t Basic) Decoder(recv string) string {
	return readBasic(recv, t.Kind)
}

// Maybe is a type whose value can not be present.
type Maybe struct {
	ElemType string
	Elem     Type
}

// Encoder implements the Type interface.
func (t Maybe) Encoder(recv string) string {
	return fmt.Sprintf(`
{
	if x := %s; x == nil {
		if _, err := writer.Write([]byte{0}); err != nil {
			return nil, err
		}
	} else {
		if _, err := writer.Write([]byte{1}); err != nil {
			return nil, err
		}

		%s
	}
}
`, recv, t.Elem.Encoder(recv))
}

// Decoder implements the Type interface.
func (t Maybe) Decoder(recv string) string {
	return fmt.Sprintf(`
{
	var v = make([]byte, 1)
	if _, err := io.ReadFull(reader, v); err != nil {
		return err
	}

	if v[0] == 0 {
		%[1]s = nil
	} else {
		var tmp_%[1]s %[2]s
		%[3]s
		%[1]s = &tmp_%[1]s
	}
}
`, recv, t.ElemType, t.Elem.Decoder("tmp_"+recv))
}

// Slice type.
type Slice struct {
	ElemType string
	Elem     Type
}

// Encoder implements the Type interface.
func (t Slice) Encoder(recv string) string {
	return fmt.Sprintf(`
{
	len := len(%s)
	ux := uint64(len) << 1
	if len < 0 {
		ux = ^ux
	}
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, ux)
	_, err := writer.Write(bs)
	if err != nil {
		return nil, err
	}
	
	for i := 0; i < len; i++ %s
}
`, recv, strings.TrimSpace(t.Elem.Encoder(recv+"[i]")))
}

// Decoder implements the Type interface.
func (t Slice) Decoder(recv string) string {
	return fmt.Sprintf(`
{
	var bs = make([]byte, 8)
	if _, err := io.ReadFull(reader, bs); err != nil {
		return err
	}

	ux := binary.LittleEndian.Uint64(bs)
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}

	len := int(x)

	%s = make([]%s, len)

	for i := 0; i < len; i++ %s
}
`, recv, t.ElemType, strings.TrimSpace(t.Elem.Decoder(recv+"[i]")))
}

// Array type of fixed size.
type Array struct {
	Len  int64
	Elem Type
}

// Encoder implements the Type interface.
func (t Array) Encoder(recv string) string {
	return fmt.Sprintf(`
{
	for i := 0; i < %d; i++ %s
}
`, t.Len, strings.TrimSpace(t.Elem.Encoder(recv+"[i]")))
}

// Decoder implements the Type interface.
func (t Array) Decoder(recv string) string {
	return fmt.Sprintf(`
{
	for i := 0; i < %d; i++ %s
}
`, t.Len, strings.TrimSpace(t.Elem.Decoder(recv+"[i]")))
}

// Map type.
type Map struct {
	KeyType  string
	ElemType string
	Key      Type
	Elem     Type
}

// Encoder implements the Type interface.
func (t Map) Encoder(recv string) string {
	return fmt.Sprintf(`
{
	len := len(%[1]s)
	ux := uint64(len) << 1
	if len < 0 {
		ux = ^ux
	}
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, ux)
	_, err := writer.Write(bs)
	if err != nil {
		return nil, err
	}
	
	for k, v := range %[1]s {
		%[2]s
		%[3]s
	}
}
`, recv, t.Key.Encoder("k"), t.Elem.Encoder("v"))
}

// Decoder implements the Type interface.
func (t Map) Decoder(recv string) string {
	return fmt.Sprintf(`
{
	var bs = make([]byte, 8)
	if _, err := io.ReadFull(reader, bs); err != nil {
		return err
	}

	ux := binary.LittleEndian.Uint64(bs)
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}

	len := int(x)

	%[1]s = make(map[%[2]s]%[3]s, len)

	for i := 0; i < len; i++ {
		var key %[2]s
		var value %[3]s
		%[4]s
		%[5]s
		%[1]s[key] = value
	}
}
`,
		recv,
		t.KeyType,
		t.ElemType,
		t.Key.Decoder("key"),
		t.Elem.Decoder("value"))
}

// Struct type.
type Struct struct {
	Fields []StructField
}

// Encoder implements the Type interface.
func (t Struct) Encoder(recv string) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for _, f := range t.Fields {
		buf.WriteString(f.Type.Encoder(recv + "." + f.Name))
	}
	buf.WriteString("}\n")
	return buf.String()
}

// Decoder implements the Type interface.
func (t Struct) Decoder(recv string) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for _, f := range t.Fields {
		buf.WriteString(f.Type.Decoder(recv + "." + f.Name))
	}
	buf.WriteString("}\n")
	return buf.String()
}

// StructField is a field in a struct.
type StructField struct {
	Name string
	Type Type
	// TODO: Validations
}

// Bytes is a special type for []byte.
type Bytes struct{}

// Encoder implements the Type interface.
func (t Bytes) Encoder(recv string) string {
	return fmt.Sprintf(writeBytes, recv)
}

// Decoder implements the Type interface.
func (t Bytes) Decoder(recv string) string {
	return fmt.Sprintf(readBytes, recv)
}

func parseType(t types.Type) (Type, error) {
	switch t := t.(type) {
	case *types.Named:
		return parseType(t.Underlying())
	case *types.Struct:
		return parseStruct(t)
	case *types.Pointer:
		elem, err := parseType(t.Elem())
		if err != nil {
			return nil, err
		}

		return Maybe{t.Elem().String(), elem}, nil
	case *types.Array:
		elem, err := parseType(t.Elem())
		if err != nil {
			return nil, err
		}

		return Array{t.Len(), elem}, nil
	case *types.Map:
		key, err := parseType(t.Key())
		if err != nil {
			return nil, err
		}

		elem, err := parseType(t.Elem())
		if err != nil {
			return nil, err
		}

		return Map{t.Key().String(), t.Elem().String(), key, elem}, nil
	case *types.Slice:
		if t.Elem().String() == "byte" {
			return Bytes{}, nil
		}

		elem, err := parseType(t.Elem())
		if err != nil {
			return nil, err
		}

		return Slice{t.Elem().String(), elem}, nil
	case *types.Basic:
		switch t.Kind() {
		case types.String,
			types.Bool,
			types.Int,
			types.Int8,
			types.Int16,
			types.Int32,
			types.Int64,
			types.Uint,
			types.Uint8,
			types.Uint16,
			types.Uint32,
			types.Uint64,
			types.Uintptr,
			types.Float32,
			types.Float64:
			return Basic{t.Kind()}, nil
		default:
			return nil, fmt.Errorf("type contains a basic type which cannot be serialized (unsafe pointer or complex number)")
		}
	case *types.Chan:
		return nil, fmt.Errorf("type contains a channel type which cannot be serialized")
	case *types.Signature:
		return nil, fmt.Errorf("type contains a function type which cannot be serialized")
	case *types.Interface:
		return nil, fmt.Errorf("type contains an interface type which cannot be serialized")
	default:
		return nil, fmt.Errorf("invalid type received: %T", t)
	}
}

func parseStruct(t *types.Struct) (Type, error) {
	var s Struct
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		cfg, err := parseTag(t.Tag(i))
		if err != nil {
			return nil, fmt.Errorf("error parsing tag of field %s: %s", f.Name(), err)
		}

		if cfg.ignore {
			continue
		}

		ft, err := parseType(f.Type())
		if err != nil {
			return nil, fmt.Errorf("on field %s: %s", f.Name(), err)
		}

		s.Fields = append(s.Fields, StructField{f.Name(), ft})
	}
	return s, nil
}

type fieldConfig struct {
	ignore bool
	// TODO: validations
}

func parseTag(tag string) (*fieldConfig, error) {
	tag = reflect.StructTag(tag).Get("bindec")
	var tags []string
	for _, p := range strings.Split(tag, ",") {
		tags = append(tags, strings.TrimSpace(p))
	}

	var cfg fieldConfig
	for _, t := range tags {
		if t == "-" {
			cfg.ignore = true
		}

		// TODO: parse validations
	}

	return &cfg, nil
}
