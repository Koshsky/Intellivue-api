package packages

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures/attributes"
)

// This message can be sent as soon as the logical connection is established and the MDS Create Event/
// Reply message sequence is finished. The message calls a method that returns monitor device data in a
// single response message.
type SinglePollDataRequest struct {
	structures.SPpdu           `json:"sp_pdu"`
	structures.ROapdus         `json:"ro_apdus"`
	structures.ROIVapdu        `json:"roiv_apdu"`
	structures.ActionArgument  `json:"action_argument"`
	structures.PollMdibDataReq `json:"poll_mdib_data_req"`
}

func (s *SinglePollDataRequest) Size() uint16 {
	return s.SPpdu.Size() + s.ROapdus.Size() + s.ROIVapdu.Size() + s.ActionArgument.Size() + s.PollMdibDataReq.Size()
}

func NewSinglePollDataRequest(invoke_id uint16, code base.OIDType) *SinglePollDataRequest {
	sp := structures.SPpdu{
		SessionID:  0xE100,
		PContextID: 2,
	}
	roap := structures.ROapdus{
		ROType: base.ROIV_APDU,
	}
	roiv := structures.ROIVapdu{
		InvokeID:    invoke_id,
		CommandType: base.CMD_CONFIRMED_ACTION,
	}
	actionArg := structures.ActionArgument{
		ManagedObject: base.ManagedObjectId{
			MObjClass: base.NOM_MOC_VMS_MDS,
			MObjInst: base.GlbHandle{
				ContextID: 0,
				Handle:    base.Handle{Value: 0},
			},
		},
		Scope:      0,
		ActionType: base.NOM_ACT_POLL_MDIB_DATA,
	}
	pollReq := structures.PollMdibDataReq{
		PollNumber: 1,
		PolledObjType: attributes.TYPE{
			Partition: base.NOM_PART_OBJ,
			Code:      code,
		},
		PolledAttrGrp: 0,
	}

	actionArg.Length = pollReq.Size()
	roiv.Length = actionArg.Length + actionArg.Size()
	roap.Length = roiv.Size() + roiv.Length

	return &SinglePollDataRequest{
		SPpdu:           sp,
		ROapdus:         roap,
		ROIVapdu:        roiv,
		ActionArgument:  actionArg,
		PollMdibDataReq: pollReq,
	}
}

func (s *SinglePollDataRequest) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	SPpduBytes, err := s.SPpdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SPpdu: %w", err)
	}
	buf.Write(SPpduBytes)

	ROapdusBytes, err := s.ROapdus.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ROapdus: %w", err)
	}
	buf.Write(ROapdusBytes)

	ROIVapduBytes, err := s.ROIVapdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ROIVapdu: %w", err)
	}
	buf.Write(ROIVapduBytes)

	ActionArgumentBytes, err := s.ActionArgument.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ActionArgument: %w", err)
	}
	buf.Write(ActionArgumentBytes)

	PollMdibDataReqBytes, err := s.PollMdibDataReq.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PollMdibDataReq: %w", err)
	}
	buf.Write(PollMdibDataReqBytes)

	return buf.Bytes(), nil
}

func (s *SinglePollDataRequest) UnmarshalBinary(r io.Reader) error {
	if s == nil {
		return fmt.Errorf("nil SinglePollDataRequest receiver")
	}

	if err := s.SPpdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal SPpdu: %w", err)
	}
	if err := s.ROapdus.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ROapdus: %w", err)
	}
	if err := s.ROIVapdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ROIVapdu: %w", err)
	}
	if err := s.ActionArgument.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ActionArgument: %w", err)
	}
	if err := s.PollMdibDataReq.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PollMdibDataReq: %w", err)
	}

	return nil
}
