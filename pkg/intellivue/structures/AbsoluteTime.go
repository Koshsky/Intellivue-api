package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type AbsoluteTime struct {
	Century      uint8
	Year         uint8
	Month        uint8
	Day          uint8
	Hour         uint8
	Minute       uint8
	Second       uint8
	SecFractions uint8
}

func (a *AbsoluteTime) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, a.Century); err != nil {
		return nil, fmt.Errorf("failed to marshal Century: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Year); err != nil {
		return nil, fmt.Errorf("failed to marshal Year: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Month); err != nil {
		return nil, fmt.Errorf("failed to marshal Month: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Day); err != nil {
		return nil, fmt.Errorf("failed to marshal Day: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Hour); err != nil {
		return nil, fmt.Errorf("failed to marshal Hour: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Minute); err != nil {
		return nil, fmt.Errorf("failed to marshal Minute: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Second); err != nil {
		return nil, fmt.Errorf("failed to marshal Second: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.SecFractions); err != nil {
		return nil, fmt.Errorf("failed to marshal SecFractions: %w", err)
	}

	return buf.Bytes(), nil
}

func (a *AbsoluteTime) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &a.Century); err != nil {
		return fmt.Errorf("failed to unmarshal Century: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Year); err != nil {
		return fmt.Errorf("failed to unmarshal Year: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Month); err != nil {
		return fmt.Errorf("failed to unmarshal Month: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Day); err != nil {
		return fmt.Errorf("failed to unmarshal Day: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Hour); err != nil {
		return fmt.Errorf("failed to unmarshal Hour: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Minute); err != nil {
		return fmt.Errorf("failed to unmarshal Minute: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Second); err != nil {
		return fmt.Errorf("failed to unmarshal Second: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.SecFractions); err != nil {
		return fmt.Errorf("failed to unmarshal SecFractions: %w", err)
	}

	return nil
}

func (a *AbsoluteTime) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<AbsoluteTime>", indent)
	log.Printf("%s  Century: %d", indent, a.Century)
	log.Printf("%s  Year: %d", indent, a.Year)
	log.Printf("%s  Month: %d", indent, a.Month)
	log.Printf("%s  Day: %d", indent, a.Day)
	log.Printf("%s  Hour: %d", indent, a.Hour)
	log.Printf("%s  Minute: %d", indent, a.Minute)
	log.Printf("%s  Second: %d", indent, a.Second)
	log.Printf("%s  SecFractions: %d", indent, a.SecFractions)
}
