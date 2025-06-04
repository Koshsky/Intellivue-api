package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf16"
)

type String struct {
	Length uint16 `json:"length"` // Длина в байтах (включая терминатор 0x0000)
	Value  string `json:"value"`  // UTF-8 строка (для удобства работы в Go)
}

// Size возвращает длину в байтах (2 байта на Length + данные строки + терминатор 0x0000)
func (s *String) Size() uint16 {
	utf16Runes := utf16.Encode([]rune(s.Value))
	return 2 + uint16(len(utf16Runes)*2) + 2 // Length (2) + данные (n*2) + 0x0000 (2)
}

// MarshalBinary кодирует строку в UTF-16BE с терминатором 0x0000.
func (s *String) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Преобразуем строку в UTF-16 (каждый символ = 2 байта)
	utf16Data := utf16.Encode([]rune(s.Value))

	// Длина в байтах (данные + терминатор 0x0000)
	dataLength := len(utf16Data)*2 + 2

	// Проверяем, что Length совпадает (если задан)
	if s.Length != 0 && uint16(dataLength) != s.Length {
		return nil, fmt.Errorf("value length (%d) does not match Length field (%d)", dataLength, s.Length)
	}

	// Устанавливаем длину, если она не задана
	if s.Length == 0 {
		s.Length = uint16(dataLength)
	}

	// Записываем длину (BigEndian)
	if err := binary.Write(buf, binary.BigEndian, s.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	// Записываем UTF-16BE данные
	for _, v := range utf16Data {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return nil, fmt.Errorf("failed to marshal UTF-16 char: %w", err)
		}
	}

	// Добавляем терминатор 0x0000
	if err := binary.Write(buf, binary.BigEndian, uint16(0)); err != nil {
		return nil, fmt.Errorf("failed to marshal terminator: %w", err)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary декодирует UTF-16BE строку с терминатором 0x0000.
func (s *String) UnmarshalBinary(r io.Reader) error {
	// Читаем длину (2 байта)
	if err := binary.Read(r, binary.BigEndian, &s.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	if s.Length == 0 {
		s.Value = ""
		return nil
	}

	// Читаем данные (длина включает терминатор)
	data := make([]byte, s.Length)
	if _, err := io.ReadFull(r, data); err != nil {
		return fmt.Errorf("failed to read string data: %w", err)
	}

	// Проверяем терминатор (последние 2 байта должны быть 0x0000)
	if len(data) < 2 || binary.BigEndian.Uint16(data[len(data)-2:]) != 0 {
		return fmt.Errorf("missing UTF-16 terminator (0x0000)")
	}

	// Преобразуем UTF-16BE в UTF-8 (игнорируем терминатор)
	utf16Runes := make([]uint16, (len(data)-2)/2)
	for i := 0; i < len(utf16Runes); i++ {
		utf16Runes[i] = binary.BigEndian.Uint16(data[i*2 : (i+1)*2])
	}

	s.Value = string(utf16.Decode(utf16Runes))
	return nil
}

func (s *String) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<String> Length: %d, Value: '%s'", indent, s.Length, s.Value)
}
