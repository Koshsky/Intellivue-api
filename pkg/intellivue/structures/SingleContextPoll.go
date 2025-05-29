package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SingleContextPoll представляет структуру данных для опроса в рамках одного контекста.
// На основе предоставленного фрагмента кода.
type SingleContextPoll struct {
	ContextID       uint16
	Count           uint16
	Length          uint16
	ObservationPoll uint8 // Тип из предоставленного фрагмента, может быть некорректным
}

// Size возвращает общую длину SingleContextPoll в байтах.
// Длина включает ContextID, Count, Length (каждое 2 байта) + ObservationPoll (1 байт).
func (s *SingleContextPoll) Size() uint16 {
	return 2 + 2 + 2 + 1 // ContextID + Count + Length + ObservationPoll
}

// MarshalBinary кодирует структуру SingleContextPoll в бинарный формат.
func (s *SingleContextPoll) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, s.ContextID); err != nil {
		return nil, fmt.Errorf("ошибка записи ContextID: %w", err)
	}

	// Count и Length должны отражать последующие данные.
	// В данном случае, только поле ObservationPoll.
	count := uint16(1)         // Предполагаем, что ObservationPoll uint8 - это один элемент
	contentLength := uint16(1) // Длина ObservationPoll uint8

	if err := binary.Write(buf, binary.BigEndian, count); err != nil {
		return nil, fmt.Errorf("ошибка записи Count: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, contentLength); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, s.ObservationPoll); err != nil {
		return nil, fmt.Errorf("ошибка записи ObservationPoll: %w", err)
	}

	return buf.Bytes(), nil
}
