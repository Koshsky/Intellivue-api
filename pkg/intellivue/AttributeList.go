package intellivue

import (
	"bytes"
	"encoding/binary"
)

type AttributeList struct {
	Count           uint16
	Value           []AVAType
	OptionalPackages *AttributeList
}

func NewAttributeList() *AttributeList {
	return &AttributeList{
		Count:           0,
		Value:           make([]AVAType, 0),
		OptionalPackages: nil,
	}
}

func (a *AttributeList) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.Count); err != nil {
		return nil, err
	}
	
	length := a.Length()
	if err := binary.Write(&buf, binary.BigEndian, length); err != nil {
		return nil, err
	}

	for _, ava := range a.Value {
		avaData, err := ava.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(avaData)
	}

	return buf.Bytes(), nil
}

func (a *AttributeList) Length() uint16 {
	totalLength := uint16(0)

	for _, ava := range a.Value {
		totalLength += ava.Length() + 4
	}

	return totalLength
}
