package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type RORSapdu struct {
	InvokeID    uint16
	CommandType uint16
	Length      uint16
}

func (r RORSapdu) Size() uint16 {
	return 3 * 2
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
