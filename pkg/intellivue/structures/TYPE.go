package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type TYPE struct {
	Partition uint16
	Code      uint16
}

func (t *TYPE) Size() uint16 {
	return 4
}

func (t *TYPE) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, t.Partition); err != nil {
		return nil, fmt.Errorf("failed to write Partition: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, t.Code); err != nil {
		return nil, fmt.Errorf("failed to write Code: %w", err)
	}

	return buf.Bytes(), nil
}

func (t *TYPE) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &t.Partition); err != nil {
		return fmt.Errorf("failed to read Partition: %w", err)
	}

	if err := binary.Read(r, binary.BigEndian, &t.Code); err != nil {
		return fmt.Errorf("failed to read Code: %w", err)
	}

	return nil
}
