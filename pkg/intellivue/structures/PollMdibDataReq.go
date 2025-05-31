package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PollMdibDataReq struct {
	PollNumber    uint16 // recommended to use this field as a counter
	PolledObjType TYPE   // Numerics / Alarms / MDS / Patient Demographics
	PolledAttrGrp OIDType
}

func (p *PollMdibDataReq) Size() uint16 {
	return 2 + p.PolledObjType.Size() + 2
}

func (p *PollMdibDataReq) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.PollNumber); err != nil {
		return nil, fmt.Errorf("ошибка записи PollNumber: %w", err)
	}

	typeData, err := p.PolledObjType.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга PolledObjType: %w", err) // Fixed error message
	}

	if _, err := buf.Write(typeData); err != nil {
		return nil, fmt.Errorf("ошибка записи PolledObjType: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.PolledAttrGrp); err != nil {
		return nil, fmt.Errorf("ошибка записи PolledAttrGrp: %w", err)
	}

	return buf.Bytes(), nil
}

func (p PollMdibDataReq) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.PollNumber); err != nil {
		return fmt.Errorf("failed to read PollNumber: %v", err)
	}

	if err := p.PolledObjType.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PolledObjType: %v", err)
	}

	if err := binary.Read(r, binary.BigEndian, &p.PolledAttrGrp); err != nil {
		return fmt.Errorf("failed to read PolledAttrGrp: %v", err)
	}

	return nil
}
