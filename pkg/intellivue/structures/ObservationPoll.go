package structures

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures/attributes"
)

type ObservationPoll struct {
	Handle     base.Handle               `json:"handle"`
	Attributes *attributes.AttributeList `json:"attributes"`
}

func (o *ObservationPoll) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Attributes)
}

func (o *ObservationPoll) Size() uint16 {
	length := uint16(2)
	if o.Attributes != nil {
		length += o.Attributes.Size()
	} else {
		length += 4
	}
	return length
}

func (o *ObservationPoll) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, o.Handle); err != nil {
		return nil, fmt.Errorf("failed to marshal Handle: %w", err)
	}

	if o.Attributes != nil {
		attributeData, err := o.Attributes.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal AttributeList: %w", err)
		}
		buf.Write(attributeData)
	} else {
		buf.Write([]byte{0x00, 0x00, 0x00, 0x00})
	}

	return buf.Bytes(), nil
}

func (op *ObservationPoll) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &op.Handle); err != nil {
		return fmt.Errorf("failed to unmarshal Handle: %w", err)
	}

	op.Attributes = &attributes.AttributeList{}
	if err := op.Attributes.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal Attributes: %w", err)
	}

	return nil
}

func (op *ObservationPoll) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<ObservationPoll>", indent)
	log.Printf("%s  Handle: %#04x", indent, op.Handle)

	op.Attributes.ShowInfo(indentationLevel + 1)
}
