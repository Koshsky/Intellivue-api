package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ROIVapdu представляет структуру ROIVapdu из MDSPollAction.
type ROIVapdu struct {
	InvokeID    uint16
	CommandType uint16 // command_type := CMD_CONFIRMED_ACTION
	Length      uint16 // Длина оставшейся части сообщения.
}

// Size возвращает длину ROIVapdu в байтах.
func (r *ROIVapdu) Size() uint16 {
	return 2 + 2 + 2 // InvokeID (2) + CommandType (2) + Length (2)
}

// MarshalBinary кодирует структуру ROIVapdu в бинарный формат.
func (r *ROIVapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("ошибка записи InvokeID: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, r.CommandType); err != nil {
		return nil, fmt.Errorf("ошибка записи CommandType: %w", err)
	}

	// Length будет установлена при маршалинге родительской структуры (MDSPollAction)
	// или структуры, которая ее включает (ActionArgument), чтобы отразить длину последующих данных.
	// Здесь мы просто записываем текущее значение.
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}
