package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type TextID struct {
	Value uint32 `json:"label_id"`
}

func (t *TextID) Size() uint16 {
	return 4
}

func (t *TextID) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, t.Value); err != nil {
		return nil, fmt.Errorf("failed to marshal TextID: %w", err)
	}
	return buf.Bytes(), nil
}

func (t *TextID) UnmarshalBinary(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &t.Value)
}

func (t TextID) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<TextID> %d", indent, t.Value)
}
