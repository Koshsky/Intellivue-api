package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SPpdu представляет структуру SPpdu из MDSPollAction.
type SPpdu struct {
	SessionID  uint16
	PContextID uint16
}

// Size возвращает длину SPpdu в байтах.
func (s *SPpdu) Size() uint16 {
	return 2 + 2 // SessionID (2) + PContextID (2)
}

// MarshalBinary кодирует структуру SPpdu в бинарный формат.
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
