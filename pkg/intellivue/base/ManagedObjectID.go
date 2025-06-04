package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type ManagedObjectId struct {
	MObjClass OIDType   `json:"mobj_class"`
	MObjInst  GlbHandle `json:"mobj_inst"`
}

func (m *ManagedObjectId) Size() uint16 {
	return 2 + m.MObjInst.Size()
}

func (m *ManagedObjectId) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, m.MObjClass); err != nil {
		return nil, fmt.Errorf("failed to marshal MObjClass: %w", err)
	}

	mobjInstBytes, err := m.MObjInst.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MObjInst: %w", err)
	}
	buf.Write(mobjInstBytes)

	return buf.Bytes(), nil
}

func (m *ManagedObjectId) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &m.MObjClass); err != nil {
		return fmt.Errorf("failed to unmarshal MObjClass: %w", err)
	}
	if err := m.MObjInst.UnmarshalBinary(reader); err != nil {
		return fmt.Errorf("failed to unmarshal MObjInst: %w", err)
	}

	return nil
}

func (m *ManagedObjectId) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<ManagedObjectId>", indent)
	log.Printf("%s  MObjClass: %#04x", indent, m.MObjClass)
}
