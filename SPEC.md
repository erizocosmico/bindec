# bindec spec

bindec encodes a Go type as bytes, with a certain shape for each type. It does not store any information about the fields a struct may have or type information. That's up for the decoder to know.

All numbers are encoded in little endian.

## Integers

- **`uint8`/`byte`**: 1 byte.
- **`uint16`**: 2 byte.
- **`uint32`**: 4 bytes.
- **`uint64`**: 8 bytes.
- **`uintptr`**: 8 bytes.
- **`uint`**: 8 bytes.
- **`int8`**: 1 byte.
- **`int16`**: 2 bytes.
- **`int32`**: 4 bytes.
- **`int64`**: 8 bytes.
- **`int`**: 8 bytes.

## Floats

- **`float32`**: 4 bytes.
- **`float64`**: 8 bytes.

## Booleans

1 byte with `1` as value for `true` or `0` for `false`.

## Strings

- 8 bytes unsigned 64 bits integer with the number of bytes of the string.
- N bytes, where N is the string length in bytes.

```
[ 8 bytes (size) ][ N bytes ]
```

## Slices

- 8 bytes unsigned 64 bits integer with the number of elements in the slice.
- N elements, where N is the number of elements in the map.

```
[ 8 bytes (size) ][ Element 1 ][ Element 2 ] ... [ Element N ]
```

## Arrays

Unlike slices, arrays have fixed length, so there is no mention about the size in the binary representation, since that is known by the decoder.

```
[ Element 1 ][ Element 2 ] ... [ Element N ]
```

## Maps

**IMPORTANT:** due to the iteration order randomization of maps in Go, the same map is not guaranteed to produce the exact same binary representation, but it will be correctly decoded regardless of the order.

- 8 bytes unsigned 64 bits integer with the number of key-value pairs in the map.
- N key-value pairs, where N is the number of key-value pairs in the map.

```
[ 8 bytes (size) ][ Key 1 ][ Value 1 ] ... [ Key N ][ Value N ]
```

## Structs

**IMPORTANT:** struct fields are stored in the same order in which they appear in the struct, so if fields are reordered, previously encoded data will not be correctly decoded anymore.

For each field that is not ignored using the `bindec:"-"` struct tag, the representation of the field is written on the output. That means, if a field is also a struct, all fields in the field struct will be written before the following fields of the current struct.

```
[ Field 1 ][ Field 2 ] ... [ Field N ]
```

## Maybes

Maybe's can either be empty or contain a value. This is equivalent to pointers in Go.

- 1 byte with `1` as value if it's not empty, `0` otherwise`.
- If first byte was `0`, nothing else will be written. Otherwise, the element representation will be written.

Empty:

```
[ 0 ]
```

Full:
```
[ 1 ][ Element ]
```