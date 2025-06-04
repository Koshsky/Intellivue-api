package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type GlbHandle struct {
	ContextID MdsContext `json:"context_id"`
	Handle    Handle     `json:"handle"`
}

func (g *GlbHandle) Size() uint16 {
	return 2 * 2
}

func (g *GlbHandle) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, g.ContextID); err != nil {
		return nil, fmt.Errorf("failed to marshal ContextID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, g.Handle); err != nil {
		return nil, fmt.Errorf("failed to marshal Handle: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *GlbHandle) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &g.ContextID); err != nil {
		return fmt.Errorf("failed to unmarshal ContextID: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &g.Handle); err != nil {
		return fmt.Errorf("failed to unmarshal Handle: %w", err)
	}

	return nil
}

func (g *GlbHandle) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s--- GlbHandle ---", indent)
	log.Printf("%s ContextID: 0x%X", indent, g.ContextID)
	log.Printf("%s Handle: %d", indent, g.Handle)
}
