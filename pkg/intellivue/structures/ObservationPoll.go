package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ObservationPoll представляет структуру данных для одного опроса наблюдения.
type ObservationPoll struct {
	ObjHandle  uint16
	Attributes *AttributeList // Предполагается стандартная структура AttributeList с Count и Length
}

// Size возвращает длину ObservationPoll в байтах.
// Длина включает ObjHandle (2 байта) + длина AttributeList (если не nil).
func (o *ObservationPoll) Size() uint16 {
	length := uint16(2) // ObjHandle
	if o.Attributes != nil {
		length += o.Attributes.Size()
	}
	return length
}

// MarshalBinary кодирует структуру ObservationPoll в бинарный формат.
func (o *ObservationPoll) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, o.ObjHandle); err != nil {
		return nil, fmt.Errorf("ошибка записи ObjHandle: %w", err)
	}

	if o.Attributes != nil {
		attributeData, err := o.Attributes.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга AttributeList: %w", err)
		}
		buf.Write(attributeData)
	}

	return buf.Bytes(), nil
}
