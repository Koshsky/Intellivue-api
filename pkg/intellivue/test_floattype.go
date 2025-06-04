package intellivue

import (
	"fmt"
	"testing"
)

func TestFloatType(t *testing.T) {
	a := float32(0xfd007d00)
	b := float32(0xff000140)
	c := float32(0x01000140)
	d := float32(0x02000020)
	fmt.Printf("%f %f %f %f", a, b, c, d)
}
