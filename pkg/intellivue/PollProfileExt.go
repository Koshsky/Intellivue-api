package intellivue

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	POLL_EXT_PERIOD_NU_1SEC          uint32 = 0x80000000 // 1 sec Real-time Numerics
	POLL_EXT_PERIOD_NU_AVG_12SEC     uint32 = 0x40000000 // 12 sec averaged Numerics
	POLL_EXT_PERIOD_NU_AVG_60SEC     uint32 = 0x20000000 // 60 sec averaged Numerics
	POLL_EXT_PERIOD_NU_AVG_300SEC    uint32 = 0x10000000 // 300 sec averaged Numerics
	POLL_EXT_PERIOD_RTSA             uint32 = 0x08000000 // Real-time Status and Alarms
	POLL_EXT_ENUM                    uint32 = 0x04000000 // allow enumeration objects
	POLL_EXT_NU_PRIO_LIST            uint32 = 0x02000000 // allow numeric priority list to be set
	POLL_EXT_DYN_MODALITIES          uint32 = 0x01000000 // send timestamps for numerics with dynamic modalities
)

type PollProfileExtension struct {
	options uint32
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

func (o PollProfileExtension) Length() uint16 {
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
