package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ManagedObjectId представляет структуру Managed Object ID.
type ManagedObjectId struct {
	MObjClass uint16
	ContextId uint16 // Часто 0x0000
	Handle    uint16 // Часто 0x0000, указывает на сам MDS объект
}

// Size возвращает длину ManagedObjectId в байтах.
func (m *ManagedObjectId) Size() uint16 {
	return 3 * 2 // 3 поля типа uint16
}

// MarshalBinary кодирует структуру ManagedObjectId в бинарный формат.
func (m *ManagedObjectId) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, m.MObjClass); err != nil {
		return nil, fmt.Errorf("ошибка записи MObjClass: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, m.ContextId); err != nil {
		return nil, fmt.Errorf("ошибка записи ContextId: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, m.Handle); err != nil {
		return nil, fmt.Errorf("ошибка записи Handle: %w", err)
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary декодирует бинарные данные в структуру ManagedObjectId.
func (m *ManagedObjectId) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &m.MObjClass); err != nil {
		return fmt.Errorf("ошибка чтения MObjClass: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &m.ContextId); err != nil {
		return fmt.Errorf("ошибка чтения ContextId: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &m.Handle); err != nil {
		return fmt.Errorf("ошибка чтения Handle: %w", err)
	}

	return nil
}
