package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

// supported components ( Attribute: System Specification)
type MdibObjectSupport struct {
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

func NewMDIBObjSupport() *MdibObjectSupport {
	return &MdibObjectSupport{
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

func (m *MdibObjectSupport) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, m.parameter0); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter0: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter1); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter1: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter2); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter2: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter3); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter3: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter4); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter4: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter5); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter5: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter6); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter6: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter7); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter7: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter8); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter8: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter9); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter9: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter10); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter10: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter11); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter11: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.parameter12); err != nil {
		return nil, fmt.Errorf("failed to marshal parameter12: %v", err)
	}

	return buf.Bytes(), nil
}

func (m *MdibObjectSupport) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &m.parameter0); err != nil {
		return fmt.Errorf("failed to unmarshal parameter0: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter1); err != nil {
		return fmt.Errorf("failed to unmarshal parameter1: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter2); err != nil {
		return fmt.Errorf("failed to unmarshal parameter2: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter3); err != nil {
		return fmt.Errorf("failed to unmarshal parameter3: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter4); err != nil {
		return fmt.Errorf("failed to unmarshal parameter4: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter5); err != nil {
		return fmt.Errorf("failed to unmarshal parameter5: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter6); err != nil {
		return fmt.Errorf("failed to unmarshal parameter6: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter7); err != nil {
		return fmt.Errorf("failed to unmarshal parameter7: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter8); err != nil {
		return fmt.Errorf("failed to unmarshal parameter8: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter9); err != nil {
		return fmt.Errorf("failed to unmarshal parameter9: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter10); err != nil {
		return fmt.Errorf("failed to unmarshal parameter10: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter11); err != nil {
		return fmt.Errorf("failed to unmarshal parameter11: %v", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &m.parameter12); err != nil {
		return fmt.Errorf("failed to unmarshal parameter12: %v", err)
	}

	return nil
}

func (m *MdibObjectSupport) Size() uint16 {
	return 52
}

func (m *MdibObjectSupport) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<MdibObjectSupport>", indent)
	log.Printf("%s  parameter0: 0x%08X", indent, m.parameter0)
	log.Printf("%s  parameter1: 0x%08X", indent, m.parameter1)
	log.Printf("%s  parameter2: 0x%08X", indent, m.parameter2)
	log.Printf("%s  parameter3: 0x%08X", indent, m.parameter3)
	log.Printf("%s  parameter4: 0x%08X", indent, m.parameter4)
	log.Printf("%s  parameter5: 0x%08X", indent, m.parameter5)
	log.Printf("%s  parameter6: 0x%08X", indent, m.parameter6)
	log.Printf("%s  parameter7: 0x%08X", indent, m.parameter7)
	log.Printf("%s  parameter8: 0x%08X", indent, m.parameter8)
	log.Printf("%s  parameter9: 0x%08X", indent, m.parameter9)
	log.Printf("%s  parameter10: 0x%08X", indent, m.parameter10)
	log.Printf("%s  parameter11: 0x%08X", indent, m.parameter11)
	log.Printf("%s  parameter12: 0x%08X", indent, m.parameter12)
}
