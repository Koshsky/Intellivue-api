package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type PollInfoList struct {
	Value []SingleContextPoll
}

func (p *PollInfoList) Size() uint16 {
	return 4 + p.Length()
}

func (p *PollInfoList) Count() uint16 {
	if p == nil {
		return 0
	}
	return uint16(len(p.Value))
}

func (p *PollInfoList) Length() uint16 {
	if p == nil || len(p.Value) == 0 {
		return 0
	}

	var total uint16
	for _, op := range p.Value {
		total += op.Size()
	}
	return total
}

func (p *PollInfoList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.Count()); err != nil {
		return nil, fmt.Errorf("failed to marshal Count: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Length()); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	for i, poll := range p.Value {
		pollData, err := poll.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal SingleContextPoll %d: %w", i, err)
		}
		buf.Write(pollData)
	}

	return buf.Bytes(), nil
}

func (p *PollInfoList) UnmarshalBinary(r io.Reader) error {
	var pollCount uint16
	if err := binary.Read(r, binary.BigEndian, &pollCount); err != nil {
		return fmt.Errorf("failed to unmarshal Count: %w", err)
	}
	var pollDataLength uint16
	if err := binary.Read(r, binary.BigEndian, &pollDataLength); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	p.Value = make([]SingleContextPoll, pollCount)
	for i := uint16(0); i < pollCount; i++ {
		if err := p.Value[i].UnmarshalBinary(r); err != nil {
			return fmt.Errorf("failed to unmarshal SingleContextPoll[%d]: %w", i, err)
		}
	}

	return nil
}

func (p *PollInfoList) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<PollInfoList>", indent)
	log.Printf("%s  Value:", indent)
	log.Printf("%s  Count: %d", indent, p.Count())
	log.Printf("%s  Length: %d", indent, p.Length())

	for _, poll := range p.Value {
		poll.ShowInfo(indentationLevel + 1)
	}
}
