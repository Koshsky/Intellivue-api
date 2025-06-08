package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures/attributes"
)

type MDSEUserInfoStd struct {
	ProtocolVersion     uint32
	NomenclatureVersion uint32
	FunctionalUnits     uint32
	SystemType          uint32
	StartupMode         uint32
	OptionList          *attributes.AttributeList
	SupportedAprofiles  *attributes.AttributeList
}

// TODO: write UnmarshalBinary, ShowInfo

func (m *MDSEUserInfoStd) Size() uint16 {
	return 20 + m.OptionList.Size() + m.SupportedAprofiles.Size()
}

func (u *MDSEUserInfoStd) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<MDSEUserInfoStd>", indent)
	log.Printf("%s  ProtocolVersion: 0x%X", indent, u.ProtocolVersion)
	log.Printf("%s  NomenclatureVersion: 0x%X", indent, u.NomenclatureVersion)
	log.Printf("%s  FunctionalUnits: 0x%X", indent, u.FunctionalUnits)
	log.Printf("%s  SystemType: 0x%X", indent, u.SystemType)
	log.Printf("%s  StartupMode: 0x%X", indent, u.StartupMode)
	u.OptionList.ShowInfo(indentationLevel + 1)
	u.SupportedAprofiles.ShowInfo(indentationLevel + 1)
}

func (m *MDSEUserInfoStd) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, m.ProtocolVersion); err != nil {
		return nil, fmt.Errorf("failed to marshal ProtocolVersion: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.NomenclatureVersion); err != nil {
		return nil, fmt.Errorf("failed to marshal NomenclatureVersion: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.FunctionalUnits); err != nil {
		return nil, fmt.Errorf("failed to marshal FunctionalUnits: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.SystemType); err != nil {
		return nil, fmt.Errorf("failed to marshal SystemType: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.StartupMode); err != nil {
		return nil, fmt.Errorf("failed to marshal StartupMode: %v", err)
	}

	optionListData, err := m.OptionList.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OptionList: %v", err)
	}
	buf.Write(optionListData)

	profileData, err := m.SupportedAprofiles.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SupportedAprofiles: %v", err)
	}
	buf.Write(profileData)

	return buf.Bytes(), nil
}

func (m *MDSEUserInfoStd) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &m.ProtocolVersion); err != nil {
		return fmt.Errorf("failed to unmarshal ProtocolVersion: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &m.NomenclatureVersion); err != nil {
		return fmt.Errorf("failed to unmarshal NomenclatureVersion: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &m.FunctionalUnits); err != nil {
		return fmt.Errorf("failed to unmarshal FunctionalUnits: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &m.SystemType); err != nil {
		return fmt.Errorf("failed to unmarshal SystemType: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &m.StartupMode); err != nil {
		return fmt.Errorf("failed to unmarshal StartupMode: %v", err)
	}
	m.OptionList = &attributes.AttributeList{}
	if err := m.OptionList.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal OptionList: %w", err)
	}
	m.SupportedAprofiles = &attributes.AttributeList{}
	if err := m.SupportedAprofiles.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal SupportedAprofiles: %w", err)
	}

	return nil
}
