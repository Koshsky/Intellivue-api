package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type PollProfileSupport struct {
	PollProfileRevision base.PollProfileRevision `json:"poll_profile_revision"`
	MinPollPeriod       base.RelativeTime        `json:"min_poll_period"`
	MaxMtuRx            uint32                   `json:"max_mtu_rx"`
	MaxMtuTx            uint32                   `json:"max_mtu_tx"`
	MaxBwTx             uint32                   `json:"max_bw_tx"`
	Options             base.PollProfileOptions  `json:"options"`
	OptionalPackages    *AttributeList           `json:"optional_packages"`
}

func (p *PollProfileSupport) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<PollProfileSupport>", indent)
	log.Printf("%s  PollProfileRevision: %d", indent, p.PollProfileRevision)
	log.Printf("%s  MinPollPeriod: %d", indent, p.MinPollPeriod)
	log.Printf("%s  MaxMtuRx: %d", indent, p.MaxMtuRx)
	log.Printf("%s  MaxMtuTx: %d", indent, p.MaxMtuTx)
	log.Printf("%s  MaxBwTx: %d", indent, p.MaxBwTx)
	log.Printf("%s  Options: 0x00000000", indent)
	p.OptionalPackages.ShowInfo(indentationLevel + 1)
}

func (p PollProfileSupport) Size() uint16 {
	size := uint16(24)

	size += p.OptionalPackages.Size()
	log.Printf("Size of PollProfileSupport: %d bytes\n", size)

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
