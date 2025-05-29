package structures

import (
	"bytes"
	"encoding/binary"
)

type AVAType struct {
	AttributeID uint16
	Value       interface{}
}

func (a *AVAType) Size() uint16 {
	return a.Value.(interface{ Size() uint16 }).Size()
}

func (a *AVAType) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.AttributeID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, a.Size()); err != nil {
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
