package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

// Remote Operation Result
type RORSapdu struct {
	InvokeID    uint16  // mirrored back from op. invoke
	CommandType CMDType // identifies type of command
	Length      uint16  // no of bytes in rest of message
}

func (r RORSapdu) Size() uint16 {
	return 6
}

func (r *RORSapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("ошибка записи InvokID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.CommandType); err != nil {
		return nil, fmt.Errorf("ошибка записи CommandType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *RORSapdu) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &r.InvokeID); err != nil {
		return fmt.Errorf("ошибка чтения InvokeID: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &r.CommandType); err != nil {
		return fmt.Errorf("ошибка чтения CommandType: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &r.Length); err != nil {
		return fmt.Errorf("ошибка чтения Length: %w", err)
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
