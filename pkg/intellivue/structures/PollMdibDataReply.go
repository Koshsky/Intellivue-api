package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type PollMdibDataReply struct {
	PollNumber    uint16
	RelTimeStamp  RelativeTime
	AbsTimeStamp  AbsoluteTime
	PolledObjType TYPE
	PolledAttrGrp OIDType
	PollInfoList  *PollInfoList
}

func (p *PollMdibDataReply) Size() uint16 {
	return 20 + p.PollInfoList.Size()
}

func (p *PollMdibDataReply) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.PollNumber); err != nil {
		return nil, fmt.Errorf("failed to marshal PollNumber: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.RelTimeStamp); err != nil {
		return nil, fmt.Errorf("failed to marshal RelTimeStamp: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.AbsTimeStamp); err != nil {
		return nil, fmt.Errorf("failed to marshal AbsTimeStamp: %w", err)
	}

	PolledObjTypeData, err := p.PolledObjType.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PolledObjType: %w", err)
	}
	buf.Write(PolledObjTypeData)

	if err := binary.Write(buf, binary.BigEndian, p.PolledAttrGrp); err != nil {
		return nil, fmt.Errorf("failed to marshal PolledAttrGrp: %w", err)
	}

	pollInfoListData, err := p.PollInfoList.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PollInfoList: %w", err)
	}
	buf.Write(pollInfoListData)

	return buf.Bytes(), nil
}

func (p *PollMdibDataReply) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.PollNumber); err != nil {
		return fmt.Errorf("failed to unmarshal PollNumber: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.RelTimeStamp); err != nil {
		return fmt.Errorf("failed to unmarshal RelTimeStamp: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.AbsTimeStamp); err != nil {
		return fmt.Errorf("failed to unmarshal AbsTimeStamp: %w", err)
	}
	if err := p.PolledObjType.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PolledObjType: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.PolledAttrGrp); err != nil {
		return fmt.Errorf("failed to unmarshal PolledAttrGrp: %w", err)
	}

	p.PollInfoList = &PollInfoList{}
	if err := p.PollInfoList.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PollInfoList: %w", err)
	}

	return nil
}

func (p *PollMdibDataReply) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<PollMdibDataReply>", indent)
	log.Printf("%s  PollNumber: 0x%X", indent, p.PollNumber)
	log.Printf("%s  RelTimeStamp: 0x%X", indent, p.RelTimeStamp)
	p.AbsTimeStamp.ShowInfo(indentationLevel + 1)
	p.PolledObjType.ShowInfo(indentationLevel + 1)
	log.Printf("%s  PolledAttrGrp: %d", indent, p.PolledAttrGrp)
	p.PollInfoList.ShowInfo(indentationLevel + 1)

}
