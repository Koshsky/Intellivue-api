package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PollProfileExtOptions uint32

const (
	POLL_EXT_PERIOD_NU_1SEC       PollProfileExtOptions = 0x80000000 // 1 sec Real-time
	POLL_EXT_PERIOD_NU_AVG_12SEC  PollProfileExtOptions = 0x40000000 // 12 sec averaged
	POLL_EXT_PERIOD_NU_AVG_60SEC  PollProfileExtOptions = 0x20000000 // 60 sec averaged
	POLL_EXT_PERIOD_NU_AVG_300SEC PollProfileExtOptions = 0x10000000 // 300 sec averaged
	POLL_EXT_PERIOD_RTSA          PollProfileExtOptions = 0x08000000 // Real-time Status and Alarms
	POLL_EXT_ENUM                 PollProfileExtOptions = 0x04000000 // allow enumeration objects
	POLL_EXT_NU_PRIO_LIST         PollProfileExtOptions = 0x02000000 // allow numeric priority list to be set
	POLL_EXT_DYN_MODALITIES       PollProfileExtOptions = 0x01000000 // send timestamps for numerics with dynamic modalities
)

// The Poll Profile Extensions attribute specifies
// some extensions for the standard polling profile.
type PollProfileExt struct {
	Options PollProfileExtOptions
	ExtAttr *AttributeList // reserved for future extensions
}

func (o PollProfileExt) Size() uint16 {
	size := uint16(8)

	if o.ExtAttr != nil {
		size += o.ExtAttr.Size()
	}

	return size
}

func (o PollProfileExt) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, o.Options); err != nil {
		return nil, fmt.Errorf("failed to marshal options: %v", err)
	}

	ext_attr, err := o.ExtAttr.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ext_attr: %v", err)
	}
	buf.Write(ext_attr)

	return buf.Bytes(), nil
}

func (o PollProfileExt) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &o.Options); err != nil {
		return fmt.Errorf("failed to unmarshal Options: %v", err)
	}

	o.ExtAttr = &AttributeList{}
	if err := o.ExtAttr.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ExtAttr: %v", err)
	}

	return nil
}

// Helper function for creating the PollProfileExt AVAType
// TODO: ????? должна ли эта фунцкция располагаться ЗДЕСЬ?
func NewPollProfileExtensionAVAType() AVAType {
	return AVAType{
		AttributeID: NOM_ATTR_POLL_PROFILE_EXT,
		Value: &PollProfileExt{
			Options: POLL_EXT_PERIOD_NU_1SEC |
				POLL_EXT_PERIOD_RTSA |
				POLL_EXT_ENUM |
				POLL_EXT_NU_PRIO_LIST |
				POLL_EXT_DYN_MODALITIES,
			ExtAttr: &AttributeList{},
		},
	}
}
