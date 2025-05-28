package intellivue

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PollProfileOptions uint32
type PollProfileRevision uint32
type RelativeTime uint32

const (
	P_OPT_DYN_CREATE_OBJECTS         PollProfileOptions = 0x40000000 // Dynamic Create Objects
	P_OPT_DYN_DELETE_OBJECTS         PollProfileOptions = 0x20000000 // Dynamic Delete Objects
)

type PollProfileSupport struct {
	PollProfileRevision PollProfileRevision  // If no matching version is found, the profile is not supported.
	MinPollPeriod       RelativeTime
	MaxMtuRx            uint32
	MaxMtuTx            uint32
	MaxBwTx             uint32
	Options             PollProfileOptions
	OptionalPackages    *AttributeList
}

func NewPollProfileSupport() *PollProfileSupport {
	optionalPackages := NewAttributeList()
	optionalPackages.Value = append(optionalPackages.Value, AVAType{
		AttributeID: NOM_ATTR_POLL_PROFILE_EXT,
		Value:       NewPollProfileExtension(),
	})
	optionalPackages.Count = 1
	
	return &PollProfileSupport{
		PollProfileRevision: POLL_PROFILE_REV_0,
		MinPollPeriod:       0x00000fa0,
		MaxMtuRx:            0x000005b0,
		MaxMtuTx:            0x000005b0,
		MaxBwTx:             0xFFFFFFFF,
		Options:             P_OPT_DYN_CREATE_OBJECTS | P_OPT_DYN_DELETE_OBJECTS,
		OptionalPackages:    optionalPackages,
	}
}

func (p PollProfileSupport) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Записываем все поля в порядке их объявления
	if err := binary.Write(&buf, binary.BigEndian, p.PollProfileRevision); err != nil {
		return nil, fmt.Errorf("failed to write PollProfileRevision: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MinPollPeriod); err != nil {
		return nil, fmt.Errorf("failed to write MinPollPeriod: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MaxMtuRx); err != nil {
		return nil, fmt.Errorf("failed to write MaxMtuRx: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MaxMtuTx); err != nil {
		return nil, fmt.Errorf("failed to write MaxMtuTx: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MaxBwTx); err != nil {
		return nil, fmt.Errorf("failed to write MaxBwTx: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.Options); err != nil {
		return nil, fmt.Errorf("failed to write Options: %v", err)
	}

	// Сериализуем OptionalPackages
	optPackagesData, err := p.OptionalPackages.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OptionalPackages: %v", err)
	}
	buf.Write(optPackagesData)

	return buf.Bytes(), nil
}

func (p PollProfileSupport) Length() uint16 {
	baseLength := uint16(24)
	
	if p.OptionalPackages != nil {
		baseLength += 4 // Count(2) + Length(2)
		baseLength += p.OptionalPackages.Length()
	}
	
	return baseLength
}