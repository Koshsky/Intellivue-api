package packages

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type SessionData struct {
	Data []byte
}

func (s *SessionData) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<SessionData>", indent)
	log.Printf("%s  Data: 0x%X", indent, s.Data)
}

func (s SessionData) Size() uint16 {
	return uint16(len(s.Data))
}

func (s *SessionData) MarshalBinary() ([]byte, error) {
	return s.Data, nil
}

func (s *SessionData) UnmarshalBinary(r io.Reader) error {
	s.Data = make([]byte, 1024)
	_, err := io.ReadFull(r, s.Data)
	if err != nil {
		return fmt.Errorf("failed to read session data: %w", err)
	}
	return nil
}
