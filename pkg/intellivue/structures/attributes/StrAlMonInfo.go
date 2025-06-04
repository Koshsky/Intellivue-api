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

type StrAlMonInfo struct {
	AlertInstNo uint16             `json:"al_inst_no"`
	AlertText   base.TextID        `json:"al_text"`
	Priority    base.AlertPriority `json:"priority"`
	Flags       base.AlertFlags    `json:"flags"`
	String      base.String        `json:"string"`
}

func (s *StrAlMonInfo) Size() uint16 {
	return 2 + 4 + 2 + 2 + s.String.Size() // AlertInstNo:uint16, AlertText:uint32, Priority:uint16, Flags:uint16, String
}

func (s *StrAlMonInfo) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, s.AlertInstNo); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertInstNo: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, uint32(s.AlertText)); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertText: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(s.Priority)); err != nil {
		return nil, fmt.Errorf("failed to marshal Priority: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(s.Flags)); err != nil {
		return nil, fmt.Errorf("failed to marshal Flags: %w", err)
	}
	strData, err := s.String.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal String: %w", err)
	}
	if _, err := buf.Write(strData); err != nil {
		return nil, fmt.Errorf("failed to write String: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *StrAlMonInfo) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.AlertInstNo); err != nil {
		return fmt.Errorf("failed to unmarshal AlertInstNo: %w", err)
	}
	var alertText uint32
	if err := binary.Read(r, binary.BigEndian, &alertText); err != nil {
		return fmt.Errorf("failed to unmarshal AlertText: %w", err)
	}
	s.AlertText = base.TextID(alertText)
	var priority uint16
	if err := binary.Read(r, binary.BigEndian, &priority); err != nil {
		return fmt.Errorf("failed to unmarshal Priority: %w", err)
	}
	s.Priority = base.AlertPriority(priority)
	var flags uint16
	if err := binary.Read(r, binary.BigEndian, &flags); err != nil {
		return fmt.Errorf("failed to unmarshal Flags: %w", err)
	}
	s.Flags = base.AlertFlags(flags)
	if err := s.String.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal String: %w", err)
	}
	return nil
}

func (s *StrAlMonInfo) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<StrAlMonInfo>", indent)
	log.Printf("%s  AlertInstNo: %d", indent, s.AlertInstNo)
	log.Printf("%s  AlertText: %#04x", indent, uint32(s.AlertText))
	log.Printf("%s  Priority: %d", indent, uint16(s.Priority))
	log.Printf("%s  Flags: %#04x", indent, uint16(s.Flags))
	s.String.ShowInfo(indentationLevel + 1)
}
