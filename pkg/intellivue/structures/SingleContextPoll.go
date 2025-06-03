package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type SingleContextPoll struct {
	ContextID MdsContext
	PollInfo  []ObservationPoll
}

func (sp *SingleContextPoll) Count() uint16 {
	if sp == nil {
		return 0
	}
	return uint16(len(sp.PollInfo))
}

func (sp *SingleContextPoll) Length() uint16 {
	if sp == nil || len(sp.PollInfo) == 0 {
		return 0
	}

	var total uint16
	for _, op := range sp.PollInfo {
		total += op.Size()
	}
	return total
}

func (s *SingleContextPoll) Size() uint16 {
	return 6 + s.Length() // ContexId + count + length + list
}

func (s *SingleContextPoll) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, s.ContextID); err != nil {
		return nil, fmt.Errorf("failed to marshal ContextID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, s.Count()); err != nil {
		return nil, fmt.Errorf("failed to marshal Count: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, s.Length()); err != nil {
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

func (s *SingleContextPoll) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.ContextID); err != nil {
		return fmt.Errorf("failed to unmarshal ContextID: %w", err)
	}
	var obervationCount uint16
	if err := binary.Read(r, binary.BigEndian, &obervationCount); err != nil {
		return fmt.Errorf("failed to unmarshal Count: %w", err)
	}

	var observationDataLength uint16
	if err := binary.Read(r, binary.BigEndian, &observationDataLength); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	s.PollInfo = make([]ObservationPoll, obervationCount)
	for i := uint16(0); i < obervationCount; i++ {
		if err := s.PollInfo[i].UnmarshalBinary(r); err != nil {
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
