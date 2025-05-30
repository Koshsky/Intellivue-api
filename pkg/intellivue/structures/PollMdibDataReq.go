package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PollMdibDataReq struct {
	PollNumber    uint16
	PolledObjType TYPE
	PolledAttrGrp uint16
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
