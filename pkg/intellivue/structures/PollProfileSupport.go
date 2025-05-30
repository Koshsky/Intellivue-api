package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type RelativeTime uint32        // TODO: move all defines to one file!
type PollProfileRevision uint32 // TODO: move all defines to one file!
type PollProfileOptions uint32  // TODO: move all defines to one file!

type PollProfileSupport struct {
	PollProfileRevision uint32
	MinPollPeriod       uint32
	MaxMtuRx            uint32
	MaxMtuTx            uint32
	MaxBwTx             uint32
	Options             uint32
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

	optPackagesData, err := p.OptionalPackages.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OptionalPackages: %v", err)
	}
	buf.Write(optPackagesData)

	return buf.Bytes(), nil
}

func (p *PollProfileSupport) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.PollProfileRevision); err != nil {
		return fmt.Errorf("failed to read PollProfileRevision: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MinPollPeriod); err != nil {
		return fmt.Errorf("failed to read MinPollPeriod: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MaxMtuRx); err != nil {
		return fmt.Errorf("failed to read MaxMtuRx: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MaxMtuTx); err != nil {
		return fmt.Errorf("failed to read MaxMtuTx: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.MaxBwTx); err != nil {
		return fmt.Errorf("failed to read MaxBwTx: %v", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.Options); err != nil {
		return fmt.Errorf("failed to read Options: %v", err)
	}

	p.OptionalPackages = &AttributeList{}
	if err := p.OptionalPackages.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal OptionalPackages: %v", err)
	}

	return nil
}
