package packages

import (
	"bytes"
	"fmt"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type MDSPollAction struct {
	SPpdu
	ROapdus
	ROIVapdu
	ActionArgument
	PollMdibDataReq
}

func (m *MDSPollAction) Size() uint16 {
	return m.SPpdu.Size() + m.ROapdus.Size() + m.ROIVapdu.Size() + m.ActionArgument.Size() + m.PollMdibDataReq.Size()
}

func NewMDSPollAction(invoke_id, code uint16) *MDSPollAction {
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
	pollReq := PollMdibDataReq{
		PollNumber: 1,
		PolledObjType: TYPE{
			Partition: 0,
			Code:      code,
		},
		PolledAttrGrp: 0,
	}

	actionArg := ActionArgument{
		ManagedObject: ManagedObjectId{
			MObjClass: 0,
			MObjInst: GlbHandle{
				ContextID: 0,
				Handle:    0,
			},
		},
		Scope:      0,
		ActionType: NOM_ACT_POLL_MDIB_DATA,
	}

	roiv.Length = pollReq.Size() + actionArg.Size()
	roap.Length = roiv.Size() + roiv.Length

	return &MDSPollAction{
		SPpdu:           sp,
		ROapdus:         roap,
		ROIVapdu:        roiv,
		ActionArgument:  actionArg,
		PollMdibDataReq: pollReq,
	}
}

func (m *MDSPollAction) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	SPpduBytes, err := m.SPpdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка MarshalBinary для SPpdu: %w", err)
	}
	ROapdusBytes, err := m.ROapdus.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка MarshalBinary для ROapdus: %w", err)
	}
	ROIVapduBytes, err := m.ROIVapdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка MarshalBinary для ROIVapdu: %w", err)
	}
	ActionArgumentBytes, err := m.ActionArgument.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка MarshalBinary для ActionArgument: %w", err)
	}
	PollMdibDataReqBytes, err := m.PollMdibDataReq.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка MarshalBinary для PollMdibDataReq: %w", err)
	}

	if _, err := buf.Write(SPpduBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи SPpdu в буфер: %w", err)
	}
	if _, err := buf.Write(ROapdusBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи ROapdus в буфер: %w", err)
	}
	if _, err := buf.Write(ROIVapduBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи ROIVapdu в буфер: %w", err)
	}
	if _, err := buf.Write(ActionArgumentBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи ActionArgument в буфер: %w", err)
	}
	if _, err := buf.Write(PollMdibDataReqBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи PollMdibDataReq в буфер: %w", err)
	}

	return buf.Bytes(), nil
}
