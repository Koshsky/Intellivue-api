package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

// Remote Operation Header
type ROapdus struct {
	ROType uint16 // ID for operation
	Length uint16 // bytes to follow
}

func (r *ROapdus) Size() uint16 {
	return 4
}

func (r *ROapdus) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.ROType); err != nil {
		return nil, fmt.Errorf("failed to marshal ROType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *ROapdus) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &r.ROType); err != nil {
		return fmt.Errorf("failed to unmarshal ROType: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &r.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}

func (r *ROapdus) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<ROapdus>", indent)
	log.Printf("%s  ROType: 0x%X", indent, r.ROType)
	log.Printf("%s  Length: %d", indent, r.Length)
}
