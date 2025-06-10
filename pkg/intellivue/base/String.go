package base

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf16"
)

const (
	SUBSCRIPT_CAPITAL_E_CHAR       = 0xE145
	SUBSCRIPT_CAPITAL_L_CHAR       = 0xE14C
	LITER_PER_CHAR                 = 0xE400
	HYDROGEN_CHAR                  = 0xE401
	ALARM_STAR_CHAR                = 0xE40D
	CAPITAL_V_WITH_DOT_ABOVE_CHAR  = 0xE425
	ZERO_WIDTH_NO_BREAK_SPACE_CHAR = 0xFEFF
)

var privateSymbols = map[rune]string{
	SUBSCRIPT_CAPITAL_E_CHAR:       "E",
	SUBSCRIPT_CAPITAL_L_CHAR:       "L",
	LITER_PER_CHAR:                 "l/min",
	HYDROGEN_CHAR:                  "cmH20",
	ALARM_STAR_CHAR:                "*",
	CAPITAL_V_WITH_DOT_ABOVE_CHAR:  "V",
	ZERO_WIDTH_NO_BREAK_SPACE_CHAR: " ",
}

type String struct {
	Length uint16
	Value  string
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"label": s.Value,
	})
}

func (s *String) Size() uint16 {
	utf16Runes := utf16.Encode([]rune(s.Value))
	return 2 + uint16(len(utf16Runes)*2) + 2
}

func (s *String) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	utf16Data := utf16.Encode([]rune(s.Value))

	var result []uint16

	for _, v := range utf16Data {
		if privateSymbol, ok := privateSymbols[rune(v)]; ok {
			result = append(result, utf16.Encode([]rune(privateSymbol))...)
		} else {
			result = append(result, v)
		}
	}

	utf16Data = result

	dataLength := len(utf16Data)*2 + 2

	if s.Length != 0 && uint16(dataLength) != s.Length {
		return nil, fmt.Errorf("value length (%d) does not match Length field (%d)", dataLength, s.Length)
	}

	if s.Length == 0 {
		s.Length = uint16(dataLength)
	}

	if err := binary.Write(buf, binary.BigEndian, s.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	for _, v := range utf16Data {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return nil, fmt.Errorf("failed to marshal UTF-16 char: %w", err)
		}
	}

	if err := binary.Write(buf, binary.BigEndian, uint16(0)); err != nil {
		return nil, fmt.Errorf("failed to marshal terminator: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *String) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	if s.Length == 0 {
		s.Value = ""
		return nil
	}

	data := make([]byte, s.Length)
	if _, err := io.ReadFull(r, data); err != nil {
		return fmt.Errorf("failed to read string data: %w", err)
	}

	if len(data) < 2 || binary.BigEndian.Uint16(data[len(data)-2:]) != 0 {
		return fmt.Errorf("missing UTF-16 terminator (0x0000)")
	}

	utf16Runes := make([]uint16, (len(data)-2)/2)
	for i := 0; i < len(utf16Runes); i++ {
		utf16Runes[i] = binary.BigEndian.Uint16(data[i*2 : (i+1)*2])
	}

	s.Value = string(utf16.Decode(utf16Runes))

	for _, v := range s.Value {
		if privateSymbol, ok := privateSymbols[v]; ok {
			s.Value = strings.ReplaceAll(s.Value, string(v), privateSymbol)
		}
	}
	return nil
}

func (s *String) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<String> Length: %d, Value: '%s'", indent, s.Length, s.Value)
}
