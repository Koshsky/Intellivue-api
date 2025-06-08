package base

import (
	"bytes"
	"math"
	"testing"
)

func TestFLOATType_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		exponent int8
		mantissa int32
		wantErr  bool
	}{
		{"Normal value", 2, 123456, false},
		{"Negative exponent", -3, -789012, false},
		{"Max mantissa", 0, mantissaMax, false},
		{"Min mantissa", 0, mantissaMin, false},
		{"NaN", 0, mantissaNaN, false},
		{"NRes", 0, mantissaNRes, false},
		{"+INF", 0, mantissaInfPos, false},
		{"-INF", 0, mantissaInfNeg, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FLOATType{
				Exponent: tt.exponent,
				Mantissa: tt.mantissa,
			}

			// Marshal
			data, err := f.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Unmarshal
			f2 := &FLOATType{}
			if err := f2.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				t.Errorf("UnmarshalBinary() error = %v", err)
				return
			}

			// Compare
			if f.Exponent != f2.Exponent || f.Mantissa != f2.Mantissa {
				t.Errorf("Values don't match: original %+v, unmarshaled %+v", f, f2)
			}
		})
	}
}

func TestFLOATType_Float64(t *testing.T) {
	tests := []struct {
		name     string
		exponent int8
		mantissa int32
		want     float64
	}{
		{"book 1", -3, 32000, 32},
		{"book 2", -1, 320, 32},
		{"book 3", 1, 320, 3200},
		{"book 4", 2, 32, 3200},
		{"Positive value", 2, 1234, 123400},
		{"Negative value", -1, -5000, -500.0},
		{"Zero", 0, 0, 0.0},
		{"NaN", 0, mantissaNaN, math.NaN()},
		{"+INF", 0, mantissaInfPos, math.Inf(1)},
		{"-INF", 0, mantissaInfNeg, math.Inf(-1)},
		{"NRes", 0, mantissaNRes, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FLOATType{
				Exponent: tt.exponent,
				Mantissa: tt.mantissa,
			}

			got := f.ToFloat64()

			if math.IsNaN(tt.want) {
				if !math.IsNaN(got) {
					t.Errorf("Float64() = %v, want NaN", got)
				}
			} else if got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFLOATType_String(t *testing.T) {
	tests := []struct {
		name     string
		exponent int8
		mantissa int32
		want     string
	}{
		{"Positive value", 2, 1234, "123400.000000"},
		{"Negative value", -1, -5000, "-500.0"},
		{"Zero", 0, 0, "0.000000"},
		{"NaN", 0, mantissaNaN, "NaN"},
		{"+INF", 0, mantissaInfPos, "+INF"},
		{"-INF", 0, mantissaInfNeg, "-INF"},
		{"NRes", 0, mantissaNRes, "NRes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FLOATType{
				Exponent: tt.exponent,
				Mantissa: tt.mantissa,
			}

			if got := f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFLOATType_IsMethods(t *testing.T) {
	tests := []struct {
		name      string
		mantissa  int32
		isNaN     bool
		isNRes    bool
		isInfPos  bool
		isInfNeg  bool
		isSpecial bool
	}{
		{"Normal value", 1234, false, false, false, false, false},
		{"NaN", mantissaNaN, true, false, false, false, true},
		{"NRes", mantissaNRes, false, true, false, false, true},
		{"+INF", mantissaInfPos, false, false, true, false, true},
		{"-INF", mantissaInfNeg, false, false, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FLOATType{Mantissa: tt.mantissa}

			if got := f.IsNaN(); got != tt.isNaN {
				t.Errorf("IsNaN() = %v, want %v", got, tt.isNaN)
			}
			if got := f.IsNRes(); got != tt.isNRes {
				t.Errorf("IsNRes() = %v, want %v", got, tt.isNRes)
			}
			if got := f.IsInfPos(); got != tt.isInfPos {
				t.Errorf("IsInfPos() = %v, want %v", got, tt.isInfPos)
			}
			if got := f.IsInfNeg(); got != tt.isInfNeg {
				t.Errorf("IsInfNeg() = %v, want %v", got, tt.isInfNeg)
			}
			if got := f.IsSpecial(); got != tt.isSpecial {
				t.Errorf("IsSpecial() = %v, want %v", got, tt.isSpecial)
			}
		})
	}
}

func TestFLOATType_Size(t *testing.T) {
	f := &FLOATType{}
	if got := f.Size(); got != 4 {
		t.Errorf("Size() = %v, want 4", got)
	}
}
