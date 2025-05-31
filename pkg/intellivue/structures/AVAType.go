package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type AVAType struct {
	AttributeID OIDType
	Length      uint16
	Value       interface{}
}

func (a *AVAType) Size() uint16 { // TODO: это работает корректно???
	valueSize := uint16(0)
	if sizedValue, ok := a.Value.(interface{ Size() uint16 }); ok {
		valueSize = sizedValue.Size()
	}
	return 4 + valueSize
}

func (a *AVAType) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.AttributeID); err != nil {
		return nil, fmt.Errorf("ошибка записи AttributeID: %w", err)
	}

	valueSize := uint16(0)
	if sizedValue, ok := a.Value.(interface{ Size() uint16 }); ok {
		valueSize = sizedValue.Size()
	}
	if err := binary.Write(&buf, binary.BigEndian, valueSize); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	if marshalerValue, ok := a.Value.(interface{ MarshalBinary() ([]byte, error) }); ok {
		data, err := marshalerValue.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга Value: %w", err)
		}
		buf.Write(data)
	}

	return buf.Bytes(), nil
}

func (a *AVAType) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.AttributeID); err != nil {
		return fmt.Errorf("ошибка чтения AttributeID: %w", err)
	}

	if err := binary.Read(r, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("ошибка чтения Value Length: %w", err)
	}

	switch a.AttributeID {
	case NOM_MOC_VMO_METRIC_NU:
		return fmt.Errorf("ошибка парсинга Numeric Value")

	default:
		log.Printf("Неизвестный AttributeID: 0x%04X. Читаем %d байт значения как []byte.", a.AttributeID, a.Length)
		valueBytes := make([]byte, a.Length)
		if a.Length > 0 {
			if _, err := io.ReadFull(r, valueBytes); err != nil {
				return fmt.Errorf("ошибка чтения байт неизвестного Value: %w", err)
			}
		}
		a.Value = valueBytes
		log.Printf("  Прочитанные байты значения: %x", valueBytes)
	}

	return nil
}
