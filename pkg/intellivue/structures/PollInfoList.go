package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PollInfoList представляет структуру данных для списка информации об опросах.
type PollInfoList struct {
	Count             uint16              // Количество элементов SingleContextPoll в списке
	Length            uint16              // Длина оставшейся части структуры (все SingleContextPoll)
	SingleContextPoll []SingleContextPoll // Список опросов в рамках контекста
}

// Size возвращает общую длину PollInfoList в байтах.
// Длина включает Count, Length (каждое 2 байта) + суммарная длина всех SingleContextPoll.
func (p *PollInfoList) Size() uint16 {
	length := uint16(4) // Count (2) + Length (2)
	for _, poll := range p.SingleContextPoll {
		length += poll.Size() // Используем Size() из SingleContextPoll
	}
	return length
}

// MarshalBinary кодирует структуру PollInfoList в бинарный формат.
func (p *PollInfoList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Рассчитываем и записываем Count и Length.
	// Count - это количество элементов в срезе SingleContextPoll.
	// Length - это суммарная длина всех элементов SingleContextPoll в байтах.
	count := uint16(len(p.SingleContextPoll))
	contentLength := uint16(0)
	for _, poll := range p.SingleContextPoll {
		contentLength += poll.Size()
	}

	if err := binary.Write(buf, binary.BigEndian, count); err != nil {
		return nil, fmt.Errorf("ошибка записи Count: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, contentLength); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	// Маршалинг каждого SingleContextPoll в списке
	for i, poll := range p.SingleContextPoll {
		pollData, err := poll.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга SingleContextPoll %d: %w", i, err)
		}
		buf.Write(pollData)
	}

	return buf.Bytes(), nil
}
