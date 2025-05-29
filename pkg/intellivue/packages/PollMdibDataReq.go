package packages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PollMdibDataReq представляет структуру PollMdibDataReq из MDSPollAction.
type PollMdibDataReq struct {
	PollNumber    uint16
	Partition     uint16
	Code          uint16
	PolledAttrGrp uint16
}

// Size возвращает длину PollMdibDataReq в байтах.
func (p *PollMdibDataReq) Size() uint16 {
	return 2 + 2 + 2 + 2 // PollNumber (2) + Partition (2) + Code (2) + PolledAttrGrp (2)
}

// MarshalBinary кодирует структуру PollMdibDataReq в бинарный формат.
func (p *PollMdibDataReq) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.PollNumber); err != nil {
		return nil, fmt.Errorf("ошибка записи PollNumber: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.Partition); err != nil {
		return nil, fmt.Errorf("ошибка записи Partition: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.Code); err != nil {
		return nil, fmt.Errorf("ошибка записи Code: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.PolledAttrGrp); err != nil {
		return nil, fmt.Errorf("ошибка записи PolledAttrGrp: %w", err)
	}

	return buf.Bytes(), nil
}
