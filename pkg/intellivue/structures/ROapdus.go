package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ROapdus представляет структуру ROapdus из MDSPollAction.
type ROapdus struct {
	ROType uint16
	Length uint16 // Длина оставшейся части сообщения.
}

// Size возвращает длину ROapdus в байтах.
func (r *ROapdus) Size() uint16 {
	return 2 + 2 // ROType (2) + Length (2)
}

// MarshalBinary кодирует структуру ROapdus в бинарный формат.
func (r *ROapdus) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.ROType); err != nil {
		return nil, fmt.Errorf("ошибка записи ROType: %w", err)
	}

	// Length будет установлена при маршалинге родительской структуры (MDSPollAction)
	// или структуры, которая ее включает (ROIVapdu), чтобы отразить длину последующих данных.
	// Здесь мы просто записываем текущее значение.
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}
