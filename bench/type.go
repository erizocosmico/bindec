package bench

//go:generate bindec -type=Foo

// Foo is the type to be benchmarked.
type Foo struct {
	A int
	B string
	C []byte
	D struct {
		A uint64
		B string
	}
	E []int
	F [2]int
	G bool
}
