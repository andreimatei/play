package main

type T struct {
	// The size of the array matters; only 1 element is not enough for the demo.
	x [2]byte
}

//go:noinline
func foo() T {
	// t1 has a loclist before a74e5f584e
	var t1, t2 T
	_ = t2 // t2 is necessary for the demo; without it t1 is completely optimized out?
	return t1
}

func main() {
	foo()
}
