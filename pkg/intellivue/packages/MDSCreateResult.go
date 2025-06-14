package packages

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type MdsCreateEventResult struct {
	structures.SPpdu             `json:"sp_pdu"`
	structures.ROapdus           `json:"ro_apdus"`
	structures.RORSapdu          `json:"rors_apdu"`
	structures.EventReportResult `json:"event_result_report"`
}

func (m *MdsCreateEventResult) Size() uint16 {
	return m.SPpdu.Size() + m.ROapdus.Size() + m.RORSapdu.Size() + m.EventReportResult.Size()
}

func (m *MdsCreateEventResult) MarshalBinary() ([]byte, error) {
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

func (m *MdsCreateEventResult) UnmarshalBinary(r io.Reader) error {
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

func NewMdsCreateEventResult() *MdsCreateEventResult {
	return &MdsCreateEventResult{
		SPpdu: structures.SPpdu{
			SessionID:  0xE100,
			PContextID: 2,
		},
		ROapdus: structures.ROapdus{
			ROType: base.RORS_APDU,
			Length: 20,
		},
		RORSapdu: structures.RORSapdu{
			InvokeID:    1,
			CommandType: base.CMD_CONFIRMED_EVENT_REPORT,
			Length:      14,
		},
		EventReportResult: structures.EventReportResult{
			ManagedObject: base.ManagedObjectId{
				MObjClass: base.NOM_MOC_VMO_METRIC_NU,
				MObjInst: base.GlbHandle{
					ContextID: 0x0001,
					Handle:    base.Handle{0x00000000},
				},
			},
			CurrentTime: 4736768,
			EventType:   base.NOM_NOTI_MDS_CREAT,
			Length:      0,
		},
	}
}
