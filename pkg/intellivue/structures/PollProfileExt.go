package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
)

type PollProfileExtension struct {
	options  uint32
	ext_attr *AttributeList
}

func (o PollProfileExtension) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, o.options); err != nil {
		return nil, fmt.Errorf("failed to write options: %v", err)
	}
	ext_attr, err := o.ext_attr.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ext_attr: %v", err)
	}
	buf.Write(ext_attr)
	return buf.Bytes(), nil
}

func (o PollProfileExtension) Size() uint16 {
	return 8
}

func NewPollProfileExtension() *PollProfileExtension {
	return &PollProfileExtension{
		options: POLL_EXT_PERIOD_NU_1SEC |
			POLL_EXT_PERIOD_RTSA |
			POLL_EXT_ENUM |
			POLL_EXT_NU_PRIO_LIST |
			POLL_EXT_DYN_MODALITIES,
		ext_attr: NewAttributeList(),
	}
}
