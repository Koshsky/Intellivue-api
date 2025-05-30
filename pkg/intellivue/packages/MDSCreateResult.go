package packages

import (
	"bytes"
	"fmt"
	"io"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type MDSCreateResult struct { // TODO: rename to MdsCreateEventResult
	SPpdu
	ROapdus
	RORSapdu
	EventReportResult
}

func (m *MDSCreateResult) Size() uint16 {
	return m.SPpdu.Size() + m.ROapdus.Size() + m.RORSapdu.Size() + m.EventReportResult.Size()
}

func (m *MDSCreateResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	spData, err := m.SPpdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SPpdu: %w", err)
	}
	buf.Write(spData)
	roData, err := m.ROapdus.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ROapdus: %w", err)
	}
	buf.Write(roData)
	rorsData, err := m.RORSapdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RORSapdu: %w", err)
	}
	buf.Write(rorsData)
	eventData, err := m.EventReportResult.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EventReportResult: %w", err)
	}
	buf.Write(eventData)

	return buf.Bytes(), nil
}

func (m *MDSCreateResult) UnmarshalBinary(r io.Reader) error {
	if m == nil {
		return fmt.Errorf("nil MDSCreateResult receiver")
	}

	if err := m.SPpdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal SPpdu: %w", err)
	}
	if err := m.ROapdus.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ROapdus: %w", err)
	}
	if err := m.RORSapdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal RORSapdu: %w", err)
	}
	if err := m.EventReportResult.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal EventReportResult: %w", err)
	}

	return nil
}

func NewMDSCreateResult() *MDSCreateResult {
	sp := SPpdu{
		SessionID:  0xE100,
		PContextID: 2,
	}

	roap := ROapdus{
		ROType: RORS_APDU,
		Length: 0x0014,
	}

	roiv := RORSapdu{
		InvokeID:    1,
		CommandType: CMD_CONFIRMED_EVENT_REPORT,
		Length:      0x000e,
	}

	eventReport := EventReportResult{
		ManagedObject: ManagedObjectId{
			MObjClass: NOM_MOC_VMS_MDS,
			MObjInst: GlbHandle{
				ContextID: 0,
				Handle:    0,
			},
		},
		CurrentTime: 0x0000,
		EventType:   NOM_NOTI_MDS_CREAT,
		Length:      0x0000,
	}

	return &MDSCreateResult{
		SPpdu:             sp,
		ROapdus:           roap,
		RORSapdu:          roiv,
		EventReportResult: eventReport,
	}
}
