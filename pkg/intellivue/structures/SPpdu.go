package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

// Session/Presentation Header
type SPpdu struct {
	SessionID  uint16 // TODO: contains a fixed value 0xE100
	PContextID uint16 // negotiated in association phase
}

func (s *SPpdu) Size() uint16 {
	return 4
}

func (s *SPpdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, s.SessionID); err != nil {
		return nil, fmt.Errorf("ошибка записи SessionID: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, s.PContextID); err != nil {
		return nil, fmt.Errorf("ошибка записи PContextID: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *SPpdu) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &s.SessionID); err != nil {
		return fmt.Errorf("ошибка чтения SessionID: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &s.PContextID); err != nil {
		return fmt.Errorf("ошибка чтения PContextID: %w", err)
	}

	return nil
}

func (s *SPpdu) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	mu.Lock()
	log.Printf("%s<SPpdu>", indent)
	log.Printf("%s  SessionID: 0x%X", indent, s.SessionID)
	log.Printf("%s  PContextID: %d", indent, s.PContextID)
	mu.Unlock()
}
