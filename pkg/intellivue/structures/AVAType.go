package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
	// "sync" // Удаляем импорт sync, если он больше не нужен
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

// UnmarshalBinary читает данные из io.Reader и заполняет структуру AVAType.
// Возвращает сигнатуру к исходной без мьютекса и канала логов.
func (a *AVAType) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.AttributeID); err != nil {
		return fmt.Errorf("ошибка чтения AttributeID: %w", err)
	}

	if err := binary.Read(r, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("ошибка чтения Value Length: %w", err)
	}

	switch a.AttributeID {
	case NOM_MOC_VMO_METRIC_NU:
		// TODO: Парсинг Numeric Value. Логирование здесь не будет выполняться.
		return fmt.Errorf("TODO: Реализовать парсинг Numeric Value") // Пока возвращаем ошибку

	default:
		// Логирование неизвестного AttributeID и прочитанных байт будет выполняться в вызывающем коде.
		valueBytes := make([]byte, a.Length)
		if a.Length > 0 {
			limitedReader := io.LimitReader(r, int64(a.Length))
			if _, err := io.ReadFull(limitedReader, valueBytes); err != nil {
				return fmt.Errorf("ошибка чтения байт неизвестного Value: %w", err)
			}
		}
		a.Value = valueBytes
	}

	return nil
}

func (a *AVAType) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	mu.Lock()
	log.Printf("%s<AttributeList>", indent)
	log.Printf("%s  AttributeID: %d", indent, a.AttributeID)
	log.Printf("%s  Length: %d", indent, a.Length)
	if valueBytes, ok := a.Value.([]byte); ok {
		log.Printf("%s  value: %s", indent, utils.PrintHex(valueBytes))
	} else {
		log.Printf("%s  value: <not []byte>", indent)
	}
	mu.Unlock()
}
