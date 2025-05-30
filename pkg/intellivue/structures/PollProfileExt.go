package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
)

type PollProfileExtension struct {
	Options uint32
	ExtAttr *AttributeList
}

func (o PollProfileExtension) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, o.Options); err != nil {
		return nil, fmt.Errorf("failed to write options: %v", err)
	}
	ext_attr, err := o.ExtAttr.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ext_attr: %v", err)
	}
	buf.Write(ext_attr)
	return buf.Bytes(), nil
}

func (o PollProfileExtension) Size() uint16 {
	return 8
}

// Helper function for creating the PollProfileExtension AVAType
// TODO: убрать в какой-нибудь другой файл?
func NewPollProfileExtensionAVAType() AVAType {
	return AVAType{
		AttributeID: NOM_ATTR_POLL_PROFILE_EXT,
		Value: &PollProfileExtension{
			Options: POLL_EXT_PERIOD_NU_1SEC |
				POLL_EXT_PERIOD_RTSA |
				POLL_EXT_ENUM |
				POLL_EXT_NU_PRIO_LIST |
				POLL_EXT_DYN_MODALITIES,
			ExtAttr: &AttributeList{},
		},
	}
}
