package base

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"strings"
)

const (
	COL_BLACK       uint16 = 0
	COL_RED         uint16 = 1
	COL_GREEN       uint16 = 2
	COL_YELLOW      uint16 = 3
	COL_BLUE        uint16 = 4
	COL_MAGENTA     uint16 = 5
	COL_CYAN        uint16 = 6
	COL_WHITE       uint16 = 7
	COL_PINK        uint16 = 20
	COL_ORANGE      uint16 = 35
	COL_LIGHT_GREEN uint16 = 50
	COL_LIGHT_RED   uint16 = 65
)

type SimpleColour struct {
	Value uint16
}

func (h SimpleColour) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h SimpleColour) Size() uint16 {
	return 2
}

func (h SimpleColour) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, h.Value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *SimpleColour) UnmarshalBinary(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &h.Value)
}

func (h *SimpleColour) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%sColour: %s", indent, h.String())
}

func (f *SimpleColour) String() string {
	switch f.Value {
	case COL_BLACK:
		return "Black"
	case COL_RED:
		return "Red"
	case COL_GREEN:
		return "Green"
	case COL_YELLOW:
		return "Yellow"
	case COL_BLUE:
		return "Blue"
	case COL_MAGENTA:
		return "Magenta"
	case COL_CYAN:
		return "Cyan"
	case COL_WHITE:
		return "White"
	case COL_PINK:
		return "Pink"
	case COL_ORANGE:
		return "Orange"
	case COL_LIGHT_GREEN:
		return "Light Green"
	case COL_LIGHT_RED:
		return "Light Red"
	default:
		return "Unknown Color"
	}
}
