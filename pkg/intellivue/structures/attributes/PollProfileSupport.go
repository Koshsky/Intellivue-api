package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

// Attribute: Poll Profile Support
type PollProfileSupport struct {
	PollProfileRevision base.PollProfileRevision
	MinPollPeriod       base.RelativeTime
	MaxMtuRx            uint32
	MaxMtuTx            uint32
	MaxBwTx             uint32
	Options             base.PollProfileOptions
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
