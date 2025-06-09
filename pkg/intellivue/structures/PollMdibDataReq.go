package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures/attributes"
)

type PollMdibDataReq struct {
	PollNumber    uint16          // recommended to use this field as a counter
	PolledObjType attributes.TYPE // Numerics / Alarms / MDS / Patient Demographics
	PolledAttrGrp base.OIDType
}

func (p *PollMdibDataReq) Size() uint16 {
	return 2 + p.PolledObjType.Size() + 2
}

func (p *PollMdibDataReq) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.PollNumber); err != nil {
		return nil, fmt.Errorf("failed to marshal PollNumber: %w", err)
	}

	typeData, err := p.PolledObjType.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PolledObjType: %w", err) // Fixed error message
	}
	buf.Write(typeData)

	if err := binary.Write(buf, binary.BigEndian, p.PolledAttrGrp); err != nil {
		return nil, fmt.Errorf("failed to marshal PolledAttrGrp: %w", err)
	}

	return buf.Bytes(), nil
}

func (p *PollMdibDataReq) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.PollNumber); err != nil {
		return fmt.Errorf("failed to unmarshal PollNumber: %v", err)
	}
	if err := p.PolledObjType.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PolledObjType: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.PolledAttrGrp); err != nil {
		return fmt.Errorf("failed to unmarshal PolledAttrGrp: %v", err)
	}

	return nil
}

func (p *PollMdibDataReq) ShowInfo(indentationLevel int) {
	indent := strings.Repeat(" ", indentationLevel*2)
	log.Printf("%s<PollMdibDataReq>\n", indent)
	log.Printf("%s  PollNumber: %d\n", indent, p.PollNumber)
	log.Printf("%s  PolledObjType: Partition=0x%04X, Code=%d\n", indent, p.PolledObjType.Partition, p.PolledObjType.Code)
	log.Printf("%s  PolledAttrGrp: 0x%04X\n", indent, p.PolledAttrGrp)
}
