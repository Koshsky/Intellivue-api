package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type ObservationPoll struct {
	Handle     Handle
	Attributes *AttributeList
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
		return nil, fmt.Errorf("ошибка записи Handle: %w", err)
	}

	if o.Attributes != nil {
		attributeData, err := o.Attributes.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга AttributeList: %w", err)
		}
		buf.Write(attributeData)
	} else {
		buf.Write([]byte{0x00, 0x00, 0x00, 0x00})
	}

	return buf.Bytes(), nil
}

func (op *ObservationPoll) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &op.Handle); err != nil {
		return fmt.Errorf("ошибка чтения Handle в ObservationPoll: %w", err)
	}

	op.Attributes = &AttributeList{}
	if err := op.Attributes.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("ошибка парсинга Attributes в ObservationPoll: %w", err)
	}

	return nil
}
