package bindec

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// Type can generate decoders and encoders for a given type.
type Type interface {
	// Decoder generates a decoder for the type. recv is the variable or
	// struct field the data will be decoded into.
	// If root is true, it means the decoder is being generated for a recv
	// itself.
	Decoder(recv string, root bool) string
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
	TypeName string
	Kind     BasicKind
}

// Encoder implements the Type interface.
func (t Basic) Encoder(recv string) string {
	switch t.Kind {
	case types.String:
		return fmt.Sprintf(writeString, recv)
	case types.Bool:
		return fmt.Sprintf(writeBool, recv)
	case types.Int:
		return fmt.Sprintf(writeInt, recv)
	case types.Int8:
		return fmt.Sprintf(writeInt8, recv)
	case types.Int16:
		return fmt.Sprintf(writeInt16, recv)
	case types.Int32:
		return fmt.Sprintf(writeInt32, recv)
	case types.Int64:
		return fmt.Sprintf(writeInt64, recv)
	case types.Uint:
		return fmt.Sprintf(writeUint, recv)
	case types.Uint8:
		return fmt.Sprintf(writeByte, recv)
	case types.Uint16:
		return fmt.Sprintf(writeUint16, recv)
	case types.Uint32:
		return fmt.Sprintf(writeUint32, recv)
	case types.Uint64:
		return fmt.Sprintf(writeUint64, recv)
	case types.Uintptr:
		return fmt.Sprintf(writeUintptr, recv)
	case types.Float32:
		return fmt.Sprintf(writeFloat32, recv)
	case types.Float64:
		return fmt.Sprintf(writeFloat64, recv)
	default:
		return ""
	}
}

// Decoder implements the Type interface.
func (t Basic) Decoder(recv string, root bool) string {
	prefix := recvPrefix(root)
	switch t.Kind {
	case types.String:
		return fmt.Sprintf(readString, t.TypeName, recv, prefix)
	case types.Bool:
		return fmt.Sprintf(readBool, t.TypeName, recv, prefix)
	case types.Int:
		return fmt.Sprintf(readInt, t.TypeName, recv, prefix)
	case types.Int8:
		return fmt.Sprintf(readInt8, t.TypeName, recv, prefix)
	case types.Int16:
		return fmt.Sprintf(readInt16, t.TypeName, recv, prefix)
	case types.Int32:
		return fmt.Sprintf(readInt32, t.TypeName, recv, prefix)
	case types.Int64:
		return fmt.Sprintf(readInt64, t.TypeName, recv, prefix)
	case types.Uint:
		return fmt.Sprintf(readUint, t.TypeName, recv, prefix)
	case types.Uint8:
		return fmt.Sprintf(readByte, t.TypeName, recv, prefix)
	case types.Uint16:
		return fmt.Sprintf(readUint16, t.TypeName, recv, prefix)
	case types.Uint32:
		return fmt.Sprintf(readUint32, t.TypeName, recv, prefix)
	case types.Uint64:
		return fmt.Sprintf(readUint64, t.TypeName, recv, prefix)
	case types.Uintptr:
		return fmt.Sprintf(readUintptr, t.TypeName, recv, prefix)
	case types.Float32:
		return fmt.Sprintf(readFloat32, t.TypeName, recv, prefix)
	case types.Float64:
		return fmt.Sprintf(readFloat64, t.TypeName, recv, prefix)
	default:
		return ""
	}
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
			return err
		}
	} else {
		if _, err := writer.Write([]byte{1}); err != nil {
			return err
		}

		%s
	}
}
`, recv, t.Elem.Encoder(fmt.Sprintf("(*%s)", recv)))
}

// Decoder implements the Type interface.
func (t Maybe) Decoder(recv string, root bool) string {
	tmpIdent := tmpIdent(recv)
	return fmt.Sprintf(`
{
	var v = make([]byte, 1)
	if _, err := io.ReadFull(reader, v); err != nil {
		return err
	}

	if v[0] == 0 {
		%[2]s = nil
	} else {
		var %[1]s %[3]s
		%[4]s
		%[2]s = &%[1]s
	}
}
`, tmpIdent, recv, t.ElemType, t.Elem.Decoder(tmpIdent, false))
}

// Slice type.
type Slice struct {
	TypeName string
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
		return err
	}
	
	for i := 0; i < len; i++ %s
}
`, recv, strings.TrimSpace(t.Elem.Encoder(recv+"[i]")))
}

// Decoder implements the Type interface.
func (t Slice) Decoder(recv string, root bool) string {
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

	%[1]s%[2]s = make(%[3]s, len)

	for i := 0; i < len; i++ %[4]s
}
`,
		recvPrefix(root),
		recv,
		t.TypeName,
		strings.TrimSpace(t.Elem.Decoder(
			fmt.Sprintf("(%s%s)[i]", recvPrefix(root), recv),
			false,
		)),
	)
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
func (t Array) Decoder(recv string, root bool) string {
	return fmt.Sprintf(`
{
	for i := 0; i < %d; i++ %s
}
`, t.Len, strings.TrimSpace(t.Elem.Decoder(
		fmt.Sprintf("(%s%s)[i]", recvPrefix(root), recv),
		false,
	)))
}

// Map type.
type Map struct {
	TypeName, KeyType, ElemType string
	Key                         Type
	Elem                        Type
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
		return err
	}
	
	for k, v := range %[1]s {
		%[2]s
		%[3]s
	}
}
`, recv, t.Key.Encoder("k"), t.Elem.Encoder("v"))
}

// Decoder implements the Type interface.
func (t Map) Decoder(recv string, root bool) string {
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

	%[7]s%[1]s = make(%[2]s, len)

	for i := 0; i < len; i++ {
		var key %[3]s
		var value %[4]s
		%[5]s
		%[6]s
		(%[7]s%[1]s)[key] = value
	}
}
`,
		recv,
		t.TypeName,
		t.KeyType,
		t.ElemType,
		t.Key.Decoder("key", false),
		t.Elem.Decoder("value", false),
		recvPrefix(root),
	)
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
func (t Struct) Decoder(recv string, root bool) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for _, f := range t.Fields {
		buf.WriteString(f.Type.Decoder(recv+"."+f.Name, false))
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
type Bytes struct {
	TypeName string
}

// Encoder implements the Type interface.
func (t Bytes) Encoder(recv string) string {
	return fmt.Sprintf(writeBytes, recv)
}

// Decoder implements the Type interface.
func (t Bytes) Decoder(recv string, root bool) string {
	return fmt.Sprintf(readBytes, recv, t.TypeName, recvPrefix(root))
}

func typeName(ctx *parseContext, typ types.Type) string {
	if named, ok := typ.(*types.Named); ok {
		pkgName := named.Obj().Pkg().Name()
		pkgPath := named.Obj().Pkg().Path()
		typeName := named.Obj().Name()

		// If the path of the package starts with a / means it's an absolute
		// path, so it must be the current package. In that case we need to
		// ignore the package name in the type name and not import it.
		// TODO: fix this for windows
		if !strings.HasPrefix(pkgPath, "/") {
			ctx.addImport(pkgPath)
			return pkgName + "." + typeName
		}

		return typeName
	}

	return typ.String()
}

type parseContext struct {
	imports map[string]struct{}
	seen    []string
}

func newParseContext() *parseContext {
	return &parseContext{make(map[string]struct{}), nil}
}

func (ctx *parseContext) clone() *parseContext {
	seen := make([]string, len(ctx.seen))
	copy(seen, ctx.seen)

	return &parseContext{ctx.imports, seen}
}

func (ctx *parseContext) markSeen(typ types.Type) {
	named, ok := typ.(*types.Named)
	if ok && !stringContains(ctx.seen, named.String()) {
		ctx.seen = append(ctx.seen, named.String())
	}
}

func (ctx *parseContext) isSeen(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if ok && stringContains(ctx.seen, named.String()) {
		return true
	}
	return false
}

func (ctx *parseContext) addImport(pkg string) {
	ctx.imports[pkg] = struct{}{}
}

func (ctx *parseContext) getImports() []string {
	var result []string
	for i := range ctx.imports {
		result = append(result, i)
	}

	sort.Strings(result)
	return result
}

func stringContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func parseType(ctx *parseContext, t types.Type) (Type, error) {
	switch t := t.(type) {
	case *types.Named:
		if ctx.isSeen(t) {
			return nil, fmt.Errorf("there is a cyclic reference, type %s was already seen", t)
		}

		ctx.markSeen(t)
		typ, err := parseType(ctx, t.Underlying())
		if err != nil {
			return nil, err
		}

		return replaceTypeName(typ, typeName(ctx, t)), nil
	case *types.Struct:
		return parseStruct(ctx, t)
	case *types.Pointer:
		elem, err := parseType(ctx, t.Elem())
		if err != nil {
			return nil, err
		}

		return Maybe{
			typeName(ctx, t.Elem()),
			elem,
		}, nil
	case *types.Array:
		elem, err := parseType(ctx, t.Elem())
		if err != nil {
			return nil, err
		}

		return Array{t.Len(), elem}, nil
	case *types.Map:
		key, err := parseType(ctx.clone(), t.Key())
		if err != nil {
			return nil, err
		}

		elem, err := parseType(ctx.clone(), t.Elem())
		if err != nil {
			return nil, err
		}

		return Map{
			TypeName: typeName(ctx, t),
			KeyType:  typeName(ctx, t.Key()),
			ElemType: typeName(ctx, t.Elem()),
			Key:      key,
			Elem:     elem,
		}, nil
	case *types.Slice:
		tn := typeName(ctx, t)
		if t.Elem().String() == "byte" {
			return Bytes{tn}, nil
		}

		elem, err := parseType(ctx, t.Elem())
		if err != nil {
			return nil, err
		}

		return Slice{tn, elem}, nil
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
			return Basic{typeName(ctx, t), t.Kind()}, nil
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

func parseStruct(ctx *parseContext, t *types.Struct) (Type, error) {
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

		ft, err := parseType(ctx.clone(), f.Type())
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

func tmpIdent(recv string) string {
	var runes []rune
	for _, ru := range recv {
		if unicode.IsDigit(ru) || unicode.IsLetter(ru) || ru == '_' {
			runes = append(runes, ru)
		} else {
			runes = append(runes, '_')
		}
	}
	return "tmp_" + string(runes)
}

func recvPrefix(root bool) string {
	if root {
		return "*"
	}
	return ""
}

func replaceTypeName(t Type, name string) Type {
	switch t := t.(type) {
	case Slice:
		t.TypeName = name
		return t
	case Map:
		t.TypeName = name
		return t
	case Basic:
		t.TypeName = name
		return t
	case Bytes:
		t.TypeName = name
		return t
	default:
		return t
	}
}
