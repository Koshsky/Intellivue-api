package base

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
)

const (
	mantissaMax    = 1<<23 - 3 // 2^23 - 3
	mantissaMin    = -(1 << 23) + 3
	mantissaNaN    = 1<<23 - 1      // 0x7FFFFF
	mantissaNRes   = -(1 << 23)     // 0x800000
	mantissaInfPos = 1<<23 - 2      // 0x7FFFFE
	mantissaInfNeg = -(1 << 23) + 2 // 0x800002
	exponentMin    = -128
	exponentMax    = 127
)

type FLOATType struct {
	Exponent int8
	Mantissa int32
}

func (f FLOATType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *FLOATType) Size() uint16 {
	return 4
}

func (f *FLOATType) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, f.Exponent); err != nil {
		return nil, fmt.Errorf("failed to marshal Exponent: %w", err)
	}

	mantissa24 := f.Mantissa & 0xFFFFFF

	mantissaBytes := []byte{
		byte((mantissa24 >> 16) & 0xFF),
		byte((mantissa24 >> 8) & 0xFF),
		byte(mantissa24 & 0xFF),
	}
	if _, err := buf.Write(mantissaBytes); err != nil {
		return nil, fmt.Errorf("failed to marshal Mantissa: %w", err)
	}

	return buf.Bytes(), nil
}

func (f *FLOATType) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &f.Exponent); err != nil {
		return fmt.Errorf("failed to unmarshal Exponent: %w", err)
	}

	var mantissaBytes [3]byte
	if err := binary.Read(r, binary.BigEndian, &mantissaBytes); err != nil {
		return fmt.Errorf("failed to unmarshal Mantissa: %w", err)
	}

	mantissa24 := int32(mantissaBytes[0])<<16 | int32(mantissaBytes[1])<<8 | int32(mantissaBytes[2])

	if mantissa24&0x800000 != 0 {
		mantissa24 |= ^0xFFFFFF
	}

	f.Mantissa = mantissa24
	return nil
}

func (f *FLOATType) Float64() float64 {
	switch {
	case f.IsNaN():
		return math.NaN()
	case f.IsInfPos():
		return math.Inf(1)
	case f.IsInfNeg():
		return math.Inf(-1)
	case f.IsNRes():
		return math.NaN()
	default:
		return float64(f.Mantissa) * math.Pow10(int(f.Exponent))
	}
}

func (f *FLOATType) IsSpecial() bool {
	return f.IsNaN() || f.IsNRes() || f.IsInfPos() || f.IsInfNeg()
}

func (f *FLOATType) IsNaN() bool {
	return f.Mantissa == mantissaNaN
}

func (f *FLOATType) IsInfPos() bool {
	return f.Mantissa == mantissaInfPos
}

func (f *FLOATType) IsInfNeg() bool {
	return f.Mantissa == mantissaInfNeg
}

func (f *FLOATType) IsNRes() bool {
	return f.Mantissa == mantissaNRes
}

func (f *FLOATType) String() string {
	if f == nil {
		return "<nil>"
	}
	switch {
	case f.IsNaN():
		return "NaN"
	case f.IsInfPos():
		return "+INF"
	case f.IsInfNeg():
		return "-INF"
	case f.IsNRes():
		return "NRes"
	default:
		return fmt.Sprintf("%f", f.Float64())
	}
}
