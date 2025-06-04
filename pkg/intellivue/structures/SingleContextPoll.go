package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type SingleContextPoll struct {
	ContextID base.MdsContext   `json:"context_id"`
	Count     uint16            `json:"count"`
	Length    uint16            `json:"length"`
	PollInfo  []ObservationPoll `json:"poll_info"`
}

func (sp *SingleContextPoll) Size() uint16 {
	return 6 + sp.Length
}

func (s *SingleContextPoll) MarshalBinary() ([]byte, error) {
	s.Count = uint16(len(s.PollInfo))
	total := uint16(0)
	for _, op := range s.PollInfo {
		total += op.Size()
	}
	s.Length = total

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, s.ContextID); err != nil {
		return nil, fmt.Errorf("failed to marshal ContextID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, s.Count); err != nil {
		return nil, fmt.Errorf("failed to marshal Count: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, s.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	for _, op := range s.PollInfo {
		opData, err := op.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ObservationPoll: %w", err)
		}
		buf.Write(opData)
	}

	return buf.Bytes(), nil
}

func (sp *SingleContextPoll) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &sp.ContextID); err != nil {
		return fmt.Errorf("failed to unmarshal ContextID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &sp.Count); err != nil {
		return fmt.Errorf("failed to unmarshal Count: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &sp.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	sp.PollInfo = make([]ObservationPoll, sp.Count)
	for i := uint16(0); i < sp.Count; i++ {
		if err := sp.PollInfo[i].UnmarshalBinary(r); err != nil {
			return fmt.Errorf("failed to unmarshal ObservationPoll[%d]: %w", i, err)
		}
	}

	return nil
}

func (s *SingleContextPoll) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<SingleContextPoll>", indent)
	log.Printf("%s  ContextID: %d", indent, s.ContextID)

	for _, op := range s.PollInfo {
		op.ShowInfo(indentationLevel + 1)
	}
}
