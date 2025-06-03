package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

type PollInfoList struct {
	Value []SingleContextPoll
}

func (p *PollInfoList) Size() uint16 {
	return 4 + p.Length() // count + length + list
}

func (pil *PollInfoList) Count() uint16 {
	if pil == nil {
		return 0
	}
	return uint16(len(pil.Value))
}

func (pil *PollInfoList) Length() uint16 {
	if pil == nil || len(pil.Value) == 0 {
		return 0
	}

	var total uint16
	for _, op := range pil.Value {
		total += op.Size()
	}
	return total
}

func (pil *PollInfoList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, pil.Count()); err != nil {
		return nil, fmt.Errorf("ошибка записи Count: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, pil.Length()); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	for i, poll := range pil.Value {
		pollData, err := poll.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга SingleContextPoll %d: %w", i, err)
		}
		buf.Write(pollData)
	}

	return buf.Bytes(), nil
}

func (pil *PollInfoList) UnmarshalBinary(r io.Reader) error {
	var pollCount uint16
	if err := binary.Read(r, binary.BigEndian, &pollCount); err != nil {
		return fmt.Errorf("failed to read observation count: %w", err)
	}
	var pollDataLength uint16
	if err := binary.Read(r, binary.BigEndian, &pollDataLength); err != nil {
		return fmt.Errorf("failed to read observation data length: %w", err)
	}

	pil.Value = make([]SingleContextPoll, pollCount)
	for i := uint16(0); i < pollCount; i++ {
		if err := pil.Value[i].UnmarshalBinary(r); err != nil {
			return fmt.Errorf("failed to unmarshal SingleContextPoll at index %d: %w", i, err)
		}
	}

	return nil
}

func (p *PollInfoList) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	mu.Lock()
	log.Printf("%s<PollInfoList>", indent)
	log.Printf("%s  Value:", indent)
	log.Printf("%s  Count: %d", indent, p.Count())
	log.Printf("%s  Length: %d", indent, p.Length())
	mu.Unlock()

	for _, poll := range p.Value {
		poll.ShowInfo(mu, indentationLevel+1)
	}
}
