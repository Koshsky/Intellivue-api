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

// Remote Operation Result
type RORSapdu struct {
	InvokeID    uint16       `json:"invoke_id"`
	CommandType base.CMDType `json:"command_type"`
	Length      uint16       `json:"length"`
}

func (r RORSapdu) Size() uint16 {
	return 6
}

func (r *RORSapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("failed to marshal InvokID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.CommandType); err != nil {
		return nil, fmt.Errorf("failed to marshal CommandType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *RORSapdu) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &r.InvokeID); err != nil {
		return fmt.Errorf("failed to unmarshal InvokeID: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &r.CommandType); err != nil {
		return fmt.Errorf("failed to unmarshal CommandType: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &r.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}

func (r *RORSapdu) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<RORSapdu>", indent)
	log.Printf("%s  InvokeID: %d", indent, r.InvokeID)
	log.Printf("%s  CommandType: %d", indent, r.CommandType)
	log.Printf("%s  Length: %d", indent, r.Length)
}
