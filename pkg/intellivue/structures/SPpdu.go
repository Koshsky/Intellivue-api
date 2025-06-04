package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

// Session/Presentation Header
type SPpdu struct {
	SessionID  uint16 `json:"session_id"`
	PContextID uint16 `json:"p_context_id"`
}

func (s *SPpdu) Size() uint16 {
	return 4
}

func (s *SPpdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, s.SessionID); err != nil {
		return nil, fmt.Errorf("failed to marshal SessionID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, s.PContextID); err != nil {
		return nil, fmt.Errorf("failed to marshal PContextID: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *SPpdu) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &s.SessionID); err != nil {
		return fmt.Errorf("failed to unmarshal SessionID: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &s.PContextID); err != nil {
		return fmt.Errorf("failed to unmarshal PContextID: %w", err)
	}

	return nil
}

func (s *SPpdu) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<SPpdu>", indent)
	log.Printf("%s  SessionID: 0x%X", indent, s.SessionID)
	log.Printf("%s  PContextID: %d", indent, s.PContextID)
}
