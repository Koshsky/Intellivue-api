package main

import (
	"fmt"
)

func main() {
	a := float32(0xfd007d00) // must be 32.000
	b := float32(0xff000140) // must be 32.0
	c := float32(0x01000140) // must be 3200
	d := float32(0x02000020) // must be 3200
	fmt.Printf("%f %f %f %f\n\n", a, b, c, d)
}
