package structures

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type PollInfoList struct {
	Count  uint16              `json:"count"`
	Length uint16              `json:"length"`
	Value  []SingleContextPoll `json:"value"`
}

func (p *PollInfoList) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *PollInfoList) Size() uint16 {
	return 4 + p.Length
}

func (p *PollInfoList) MarshalBinary() ([]byte, error) {
	p.Count = uint16(len(p.Value))
	total := uint16(0)
	for _, poll := range p.Value {
		total += poll.Size()
	}
	p.Length = total

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.Count); err != nil {
		return nil, fmt.Errorf("failed to marshal Count: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	for _, poll := range p.Value {
		pollData, err := poll.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal SingleContextPoll: %w", err)
		}
		buf.Write(pollData)
	}

	return buf.Bytes(), nil
}

func (p *PollInfoList) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.Count); err != nil {
		return fmt.Errorf("failed to unmarshal Count: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &p.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	p.Value = make([]SingleContextPoll, p.Count)
	for i := uint16(0); i < p.Count; i++ {
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
	log.Printf("%s  Count: %d", indent, p.Count)
	log.Printf("%s  Length: %d", indent, p.Length)

	for _, poll := range p.Value {
		poll.ShowInfo(indentationLevel + 1)
	}
}
