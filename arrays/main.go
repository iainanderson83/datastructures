package main

// All of these benchmarks come from the following blog and
// are reproduced here for my own testing:
// https://philpearl.github.io/

type A struct {
	B int
	C int
}

const bigStructSize = 10

type bigStruct struct {
	a [bigStructSize]int
}
