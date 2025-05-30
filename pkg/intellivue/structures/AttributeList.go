package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type AttributeList struct {
	Value []AVAType
}

func (a *AttributeList) Size() uint16 {
	return 4 + a.Length() // count + length + list
}

func (a *AttributeList) Count() uint16 {
	if a.Value == nil {
		return 0
	}
	return uint16(len(a.Value))
}

func (a *AttributeList) Length() uint16 {
	if a.Value == nil || len(a.Value) == 0 {
		return 0
	}

	var total uint16
	for _, ava := range a.Value {
		total += ava.Size()
	}
	return total
}

func (a *AttributeList) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.Count()); err != nil {
		return nil, fmt.Errorf("ошибка записи Count: %w", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, a.Length()); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	for _, ava := range a.Value {
		avaData, err := ava.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга AVAType: %w", err)
		}
		buf.Write(avaData)
	}

	return buf.Bytes(), nil
}

func (a *AttributeList) UnmarshalBinary(r io.Reader) error {
	var attrCount uint16
	if err := binary.Read(r, binary.BigEndian, &attrCount); err != nil {
		return fmt.Errorf("failed to read attribute count: %w", err)
	}

	var attrDataLength uint16
	if err := binary.Read(r, binary.BigEndian, &attrDataLength); err != nil {
		return fmt.Errorf("failed to read attributes data length: %w", err)
	}

	listReader := io.LimitReader(r, int64(attrDataLength))

	a.Value = make([]AVAType, attrCount)
	for i := uint16(0); i < attrCount; i++ {
		if err := a.Value[i].UnmarshalBinary(listReader); err != nil {
			return fmt.Errorf("failed to unmarshal AVAType at index %d: %w", i, err)
		}
	}

	return nil
}
