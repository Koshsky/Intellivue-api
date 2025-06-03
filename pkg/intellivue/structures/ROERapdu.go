package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	NO_SUCH_OBJECT_CLASS    uint16 = 0x0000
	NO_SUCH_OBJECT_INSTANCE uint16 = 0x0001
	ACCESS_DENIED           uint16 = 0x0002
	GET_LIST_ERROR          uint16 = 0x0007
	SET_LIST_ERROR          uint16 = 0x0008
	NO_SUCH_ACTION          uint16 = 0x0009
	PROCESSING_FAILURE      uint16 = 0x000a
	INVALID_ARGUMENT_VALUE  uint16 = 0x000f
	INVALID_SCOPE           uint16 = 0x0010
	INVALID_OBJECT_INSTANCE uint16 = 0x0011
)

type ROERapdu struct {
	InvokeID uint16
	Length   uint16
}

func (r *ROERapdu) Size() uint16 {
	return 4
}

func (r *ROERapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("failed to marshal InvokeID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *ROERapdu) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &r.InvokeID); err != nil {
		return fmt.Errorf("failed to unmarshal InvokeID: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &r.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}
