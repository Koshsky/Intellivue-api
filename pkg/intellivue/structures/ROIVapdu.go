package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

// Remote Operation Invoke
type ROIVapdu struct {
	InvokeID    uint16       `json:"invoke_id"`
	CommandType base.CMDType `json:"command_type"`
	Length      uint16       `json:"length"`
}

func (r *ROIVapdu) Size() uint16 {
	return 6
}

func (r *ROIVapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("failed to marshal InvokeID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.CommandType); err != nil {
		return nil, fmt.Errorf("failed to marshal CommandType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *ROIVapdu) UnmarshalBinary(reader io.Reader) error {
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
