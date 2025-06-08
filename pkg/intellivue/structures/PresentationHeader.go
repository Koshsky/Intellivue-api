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

type PresentationHeader struct {
	Type   byte
	Length base.LIField
	Data   []byte
}

func (p *PresentationHeader) Size() uint16 {
	return 1 + p.Length.Size() + uint16(len(p.Data))
}

func (s *PresentationHeader) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<PresentationHeader>", indent)
	log.Printf("%s  Type: 0x%X", indent, s.Type)
	log.Printf("%s  Length: 0x%X", indent, s.Length)
	log.Printf("%s  Data: 0x%X", indent, s.Data)
}

func (p *PresentationHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.Type); err != nil {
		return nil, fmt.Errorf("failed to marshal Type: %w", err)
	}
	lengthData, err := p.Length.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ManagedObject: %w", err)
	}
	buf.Write(lengthData)
	buf.Write(p.Data)

	return buf.Bytes(), nil
}

func (p *PresentationHeader) UnmarshalBinary(r io.Reader, length uint16) error {
	if p == nil {
		return fmt.Errorf("nil PresentationHeader receiver")
	}
	if err := binary.Read(r, binary.BigEndian, &p.Type); err != nil {
		return fmt.Errorf("failed to unmarshal Type: %w", err)
	}
	if err := p.Length.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ManagedObject: %w", err)
	}
	p.Data = make([]byte, length)
	read, err := io.ReadFull(r, p.Data)
	if err != nil {
		return fmt.Errorf("failed to read session data: %w", err)
	} else if read != int(length) {
		return fmt.Errorf("failed to read session data: read %d bytes, expected %d", read, length)
	}
	return nil
}
