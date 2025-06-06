package base

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"strings"
)

type Handle struct {
	Value uint16
}

func (h Handle) Size() uint16 {
	return 2
}

func (h Handle) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, h.Value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *Handle) UnmarshalBinary(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &h.Value)
}

func (h *Handle) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%sHandle: 0x%04X", indent, h.Value)
}
