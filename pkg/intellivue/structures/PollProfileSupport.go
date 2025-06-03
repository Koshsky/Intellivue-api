package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PollProfileRevision uint32
type PollProfileOptions uint32

const (
	POLL_PROFILE_REV_0 PollProfileRevision = 0x80000000

	P_OPT_DYN_CREATE_OBJECTS PollProfileOptions = 0x40000000
	P_OPT_DYN_DELETE_OBJECTS PollProfileOptions = 0x20000000
)

// The Poll Profile Support attribute contains the specification
// of the polling profile supported by the system.
type PollProfileSupport struct {
	PollProfileRevision PollProfileRevision
	MinPollPeriod       RelativeTime
	MaxMtuRx            uint32
	MaxMtuTx            uint32
	MaxBwTx             uint32
	Options             PollProfileOptions
	OptionalPackages    *AttributeList
}

func (p PollProfileSupport) Size() uint16 {
	size := uint16(24)

	if p.OptionalPackages != nil {
		size += p.OptionalPackages.Size()
	}

	return size
}

func (p PollProfileSupport) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, p.PollProfileRevision); err != nil {
		return nil, fmt.Errorf("failed to marshal PollProfileRevision: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MinPollPeriod); err != nil {
		return nil, fmt.Errorf("failed to marshal MinPollPeriod: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MaxMtuRx); err != nil {
		return nil, fmt.Errorf("failed to marshal MaxMtuRx: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MaxMtuTx); err != nil {
		return nil, fmt.Errorf("failed to marshal MaxMtuTx: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.MaxBwTx); err != nil {
		return nil, fmt.Errorf("failed to marshal MaxBwTx: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, p.Options); err != nil {
		return nil, fmt.Errorf("failed to marshal Options: %v", err)
	}

	optPackagesData, err := p.OptionalPackages.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OptionalPackages: %v", err)
	}
	buf.Write(optPackagesData)

	return buf.Bytes(), nil
}

func (p *PollProfileSupport) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.PollProfileRevision); err != nil {
		return fmt.Errorf("failed to unmarshal PollProfileRevision: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MinPollPeriod); err != nil {
		return fmt.Errorf("failed to unmarshal MinPollPeriod: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MaxMtuRx); err != nil {
		return fmt.Errorf("failed to unmarshal MaxMtuRx: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MaxMtuTx); err != nil {
		return fmt.Errorf("failed to unmarshal MaxMtuTx: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MaxBwTx); err != nil {
		return fmt.Errorf("failed to unmarshal MaxBwTx: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.Options); err != nil {
		return fmt.Errorf("failed to unmarshal Options: %v", err)
	}

	p.OptionalPackages = &AttributeList{}
	if err := p.OptionalPackages.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal OptionalPackages: %v", err)
	}

	return nil
}
