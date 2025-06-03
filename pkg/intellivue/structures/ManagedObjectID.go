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

type ManagedObjectId struct {
	MObjClass OIDType
	MObjInst  GlbHandle
}

func (m *ManagedObjectId) Size() uint16 {
	return 2 + m.MObjInst.Size()
}

func (m *ManagedObjectId) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, m.MObjClass); err != nil {
		return nil, fmt.Errorf("ошибка записи MObjClass: %w", err)
	}

	mobjInstBytes, err := m.MObjInst.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка MarshalBinary для MObjInst: %w", err)
	}
	if _, err := buf.Write(mobjInstBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи MObjInst в буфер: %w", err)
	}

	return buf.Bytes(), nil
}

func (m *ManagedObjectId) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &m.MObjClass); err != nil {
		return fmt.Errorf("ошибка чтения MObjClass: %w", err)
	}

	if err := m.MObjInst.UnmarshalBinary(reader); err != nil {
		return fmt.Errorf("ошибка UnmarshalBinary для MObjInst: %w", err)
	}

	return nil
}

// ShowInfo выводит информацию о структуре ManagedObjectId через предоставленный канал логов.
func (m *ManagedObjectId) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	mu.Lock()
	log.Printf("%s<ManagedObjectId>", indent)
	log.Printf("%s  MObjClass: 0x%X", indent, m.MObjClass)
	mu.Unlock()
}
