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

type TYPE struct {
	Partition base.NomPartition `json:"partition"`
	Code      base.OIDType      `json:"code"`
}

func (t *TYPE) Size() uint16 {
	return 4
}

func (t *TYPE) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, t.Partition); err != nil {
		return nil, fmt.Errorf("failed to marshal Partition: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, t.Code); err != nil {
		return nil, fmt.Errorf("failed to marshal Code: %w", err)
	}

	return buf.Bytes(), nil
}

func (t *TYPE) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &t.Partition); err != nil {
		return fmt.Errorf("failed to unmarshal Partition: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &t.Code); err != nil {
		return fmt.Errorf("failed to unmarshal Code: %w", err)
	}

	return nil
}

func (t *TYPE) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<TYPE>", indent)
	log.Printf("%s  Partition: 0x%X", indent, t.Partition)
	log.Printf("%s  Code: %d", indent, t.Code)
}
