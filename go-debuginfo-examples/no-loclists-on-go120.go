package main

//go:noinline
func foo() {
	var xx [][5000]byte
	var x1 [5000]byte
	x1[4000] = 0xff
	xx = append(xx, x1)
}

//go:noinline
func main() {
	foo()
}
