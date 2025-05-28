package intellivue

import (
	"bytes"
	"encoding/binary"
)

type AVAType struct {
	AttributeID uint16
	Value       interface{}
}

func (a *AVAType) Length() uint16 {
	return a.Value.(interface{ Length() uint16 }).Length()
}

func (a *AVAType) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.AttributeID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, a.Length()); err != nil {
		return nil, err
	}

	v := a.Value.(interface{ MarshalBinary() ([]byte, error) })
	data, err := v.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Bytes(), nil
}
