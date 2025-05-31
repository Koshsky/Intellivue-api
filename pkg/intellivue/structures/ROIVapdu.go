package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Remote Operation Invoke
type ROIVapdu struct {
	InvokeID    uint16  // identifies the transaction
	CommandType CMDType // identifies type of command
	Length      uint16  // no. of bytes in rest of message
}

func (r *ROIVapdu) Size() uint16 {
	return 6
}

func (r *ROIVapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("ошибка записи InvokeID: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, r.CommandType); err != nil {
		return nil, fmt.Errorf("ошибка записи CommandType: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *ROIVapdu) UnmarshalBinary(reader io.Reader) error {
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
