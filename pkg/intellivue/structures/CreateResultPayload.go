package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// CreateResultPayload представляет полезную нагрузку результата создания.
type CreateResultPayload struct {
	EventTime uint32
	EventType uint16
	Length    uint16 // Длина того, что следует за Length (если есть)
}

// Size возвращает длину CreateResultPayload в байтах.
func (c *CreateResultPayload) Size() uint16 {
	return 4 + 2 + 2 // EventTime (4) + EventType (2) + Length (2)
}

// MarshalBinary кодирует структуру CreateResultPayload в бинарный формат.
func (c *CreateResultPayload) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, c.EventTime); err != nil {
		return nil, fmt.Errorf("ошибка записи EventTime: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, c.EventType); err != nil {
		return nil, fmt.Errorf("ошибка записи EventType: %w", err)
	}

	// Length будет установлена при маршалинге родительской структуры (MDSCreateResult)
	// чтобы отразить длину последующих данных.
	// Здесь мы просто записываем текущее значение.
	if err := binary.Write(buf, binary.BigEndian, c.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}
