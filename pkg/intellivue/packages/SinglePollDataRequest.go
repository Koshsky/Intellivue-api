package packages

import (
	"bytes"
	"fmt"
	"io"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

// This message can be sent as soon as the logical connection is established and the MDS Create Event/
// Reply message sequence is finished. The message calls a method that returns monitor device data in a
// single response message.
type SinglePollDataRequest struct {
	SPpdu
	ROapdus
	ROIVapdu
	ActionArgument
	PollMdibDataReq
}

func (s *SinglePollDataRequest) Size() uint16 {
	return s.SPpdu.Size() + s.ROapdus.Size() + s.ROIVapdu.Size() + s.ActionArgument.Size() + s.PollMdibDataReq.Size()
}

func NewSinglePollDataRequest(invoke_id uint16, code OIDType) *SinglePollDataRequest {
	sp := SPpdu{
		SessionID:  0xE100,
		PContextID: 2,
	}
	roap := ROapdus{
		ROType: ROIV_APDU,
	}
	roiv := ROIVapdu{
		InvokeID:    invoke_id,
		CommandType: CMD_CONFIRMED_ACTION,
	}
	actionArg := ActionArgument{
		ManagedObject: ManagedObjectId{
			MObjClass: NOM_MOC_VMS_MDS,
			MObjInst: GlbHandle{
				ContextID: 0,
				Handle:    0,
			},
		},
		Scope:      0,
		ActionType: NOM_ACT_POLL_MDIB_DATA,
	}
	pollReq := PollMdibDataReq{
		PollNumber: 1,
		PolledObjType: TYPE{
			Partition: NOM_PART_OBJ,
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
		return fmt.Errorf("ошибка UnmarshalBinary для SPpdu: %w", err)
	}
	if err := s.ROapdus.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("ошибка UnmarshalBinary для ROapdus: %w", err)
	}
	if err := s.ROIVapdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("ошибка UnmarshalBinary для ROIVapdu: %w", err)
	}
	if err := s.ActionArgument.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("ошибка UnmarshalBinary для ActionArgument: %w", err)
	}
	if err := s.PollMdibDataReq.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("ошибка UnmarshalBinary для PollMdibDataReq: %w", err)
	}

	return nil
}
