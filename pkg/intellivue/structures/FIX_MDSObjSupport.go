package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// TODO: что это ебать? определи где используется и что это за структура!
type MDIBObjSupport struct {
	parameter0  uint32
	parameter1  uint32
	parameter2  uint32
	parameter3  uint32
	parameter4  uint32
	parameter5  uint32
	parameter6  uint32
	parameter7  uint32
	parameter8  uint32
	parameter9  uint32
	parameter10 uint32
	parameter11 uint32
	parameter12 uint32
}

func NewMDIBObjSupport() *MDIBObjSupport {
	return &MDIBObjSupport{
		parameter0:  0x00060000,
		parameter1:  0x00010021,
		parameter2:  0x00000001,
		parameter3:  0x00010006,
		parameter4:  0x000000c9,
		parameter5:  0x000100c9,
		parameter6:  0x0000003c,
		parameter7:  0x00010005,
		parameter8:  0x00000010,
		parameter9:  0x0001002a,
		parameter10: 0x00000001,
		parameter11: 0x00010036,
		parameter12: 0x00000001,
	}
}

func (m MDIBObjSupport) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, m.parameter0); err != nil {
		return nil, fmt.Errorf("failed to write parameter0: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter1); err != nil {
		return nil, fmt.Errorf("failed to write parameter1: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter2); err != nil {
		return nil, fmt.Errorf("failed to write parameter2: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter3); err != nil {
		return nil, fmt.Errorf("failed to write parameter3: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter4); err != nil {
		return nil, fmt.Errorf("failed to write parameter4: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter5); err != nil {
		return nil, fmt.Errorf("failed to write parameter5: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter6); err != nil {
		return nil, fmt.Errorf("failed to write parameter6: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter7); err != nil {
		return nil, fmt.Errorf("failed to write parameter7: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter8); err != nil {
		return nil, fmt.Errorf("failed to write parameter8: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter9); err != nil {
		return nil, fmt.Errorf("failed to write parameter9: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter10); err != nil {
		return nil, fmt.Errorf("failed to write parameter10: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter11); err != nil {
		return nil, fmt.Errorf("failed to write parameter11: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter12); err != nil {
		return nil, fmt.Errorf("failed to write parameter12: %v", err)
	}

	return buf.Bytes(), nil
}

func (m MDIBObjSupport) Size() uint16 {
	return 52
}
