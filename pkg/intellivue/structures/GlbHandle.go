package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type GlbHandle struct {
	ContextID uint16
	Handle    uint16
}

func (g *GlbHandle) Size() uint16 {
	return 2 * 2
}

func (g *GlbHandle) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, g.ContextID); err != nil {
		return nil, fmt.Errorf("ошибка записи ContextID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, g.Handle); err != nil {
		return nil, fmt.Errorf("ошибка записи Handle: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *GlbHandle) UnmarshalBinary(reader io.Reader) error {
	// Читаем поля в порядке их следования
	if err := binary.Read(reader, binary.BigEndian, &g.ContextID); err != nil {
		return fmt.Errorf("ошибка чтения ContextID: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &g.Handle); err != nil {
		return fmt.Errorf("ошибка чтения Handle: %w", err)
	}

	return nil
}
