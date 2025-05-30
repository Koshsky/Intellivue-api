package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PollMdibDataReply struct {
	PollNumber    uint16
	RelTimeStamp  uint32
	AbsTimeStamp  uint64
	Partition     uint16
	Code          uint16
	PolledAttrGrp uint16
	PollInfoList  *PollInfoList
}

func (p *PollMdibDataReply) Size() uint16 {
	return 2 + 4 + 8 + 2 + 2 + 2 + p.PollInfoList.Size()
}

func (p *PollMdibDataReply) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.PollNumber); err != nil {
		return nil, fmt.Errorf("ошибка записи PollNumber: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.RelTimeStamp); err != nil {
		return nil, fmt.Errorf("ошибка записи RelTimeStamp: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.AbsTimeStamp); err != nil {
		return nil, fmt.Errorf("ошибка записи AbsTimeStamp: %w", err)
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

	pollInfoListData, err := p.PollInfoList.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга PollInfoList: %w", err)
	}
	buf.Write(pollInfoListData)

	return buf.Bytes(), nil
}

func (p *PollMdibDataReply) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.PollNumber); err != nil {
		return fmt.Errorf("failed to read PollMdibDataReply.PollNumber: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.RelTimeStamp); err != nil {
		return fmt.Errorf("failed to read PollMdibDataReply.RelTimeStamp: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.AbsTimeStamp); err != nil {
		return fmt.Errorf("failed to read PollMdibDataReply.AbsTimeStamp: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.Partition); err != nil {
		return fmt.Errorf("failed to read PollMdibDataReply.Partition: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.Code); err != nil {
		return fmt.Errorf("failed to read PollMdibDataReply.Code: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.PolledAttrGrp); err != nil {
		return fmt.Errorf("failed to read PollMdibDataReply.PolledAttrGrp: %w", err)
	}

	if p.PollInfoList == nil {
		p.PollInfoList = &PollInfoList{}
	}
	if err := p.PollInfoList.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PollMdibDataReply.PollInfoList: %w", err)
	}

	return nil
}
