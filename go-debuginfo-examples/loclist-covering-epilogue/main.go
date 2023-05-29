package main

var global int

//go:noinline
func foo() {
	// localVar's loclist erroneously extends to the end of foo's code -
	// including the epilogue.
	var localVar int
	localVar = global
	localVar++

	// Force the stack to grow when foo() starts executing. This isn't necessary
	// for reproducing the issue, but it is useful if trying to actually
	// demonstrate that the bug is "exploitable" by trapping on the stack
	// expansion and attempting to read the variable.
	var xx [][5000]byte
	var x1 [5000]byte
	xx = append(xx, x1)

	global = localVar
}

func main() {
	global = 42
	foo()
}
