package packages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type SessionHeader struct {
	Type   byte         `json:"type"`
	Length base.LIField `json:"length"`
}

func (s *SessionHeader) Size() uint16 {
	return 1 + s.Length.Size()
}

func (s *SessionHeader) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<SessionHeader>", indent)
	log.Printf("%s  Type: 0x%X", indent, s.Type)
	log.Printf("%s  Length: 0x%X", indent, s.Length)
}

func (s *SessionHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, s.Type); err != nil {
		return nil, fmt.Errorf("failed to marshal Type: %w", err)
	}

	lengthData, err := s.Length.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}
	buf.Write(lengthData)

	return buf.Bytes(), nil
}

func (s *SessionHeader) UnmarshalBinary(r io.Reader) error {
	if s == nil {
		return fmt.Errorf("nil SessionHeader receiver")
	}
	if err := binary.Read(r, binary.BigEndian, &s.Type); err != nil {
		return fmt.Errorf("failed to unmarshal Type: %w", err)
	}
	if err := s.Length.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}
