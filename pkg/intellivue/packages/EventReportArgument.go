package packages

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"sync"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

type EventReportArgument struct {
	SPpdu         structures.SPpdu     `json:"sp_pdu"`
	ROapdus       structures.ROapdus   `json:"ro_apdus"`
	ROIVapdu      structures.ROIVapdu  `json:"roiv_apdu"`
	ManagedObject base.ManagedObjectId `json:"managed_object"`
	EventTime     base.RelativeTime    `json:"event_time"`
	EventType     base.OIDType         `json:"event_type"`
	Length        uint16               `json:"length"`
	EventData     []byte               `json:"event_data"`
}

func (e *EventReportArgument) Size() uint16 {
	return 6 + e.ManagedObject.Size() + e.Length
}

func (e *EventReportArgument) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, e.Size()))

	spduData, err := e.SPpdu.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(spduData)

	roapdusData, err := e.ROapdus.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(roapdusData)

	roivapduData, err := e.ROIVapdu.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(roivapduData)

	managedObjectData, err := e.ManagedObject.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(managedObjectData)

	if err := binary.Write(buf, binary.BigEndian, e.EventTime); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, e.EventType); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, e.Length); err != nil {
		return nil, err
	}

	buf.Write(e.EventData)

	return buf.Bytes(), nil
}

func (e *EventReportArgument) UnmarshalBinary(r io.Reader) error {
	if err := e.SPpdu.UnmarshalBinary(r); err != nil {
		return err
	}

	if err := e.ROapdus.UnmarshalBinary(r); err != nil {
		return err
	}

	if err := e.ROIVapdu.UnmarshalBinary(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &e.ManagedObject); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &e.EventTime); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &e.EventType); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &e.Length); err != nil {
		return err
	}

	// e.EventData = make([]byte, e.Length)
	// if err := binary.Read(r, binary.BigEndian, e.EventData); err != nil {
	// 	return err
	// }

	return nil
}

func (e *EventReportArgument) ShowInfo(mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	e.SPpdu.ShowInfo(1)
	e.ROapdus.ShowInfo(1)
	e.ROIVapdu.ShowInfo(1)
	e.ManagedObject.ShowInfo(1)
	log.Printf("  EventTime: 0x%08X", e.EventTime)
	log.Printf("  EventType: 0x%04X", e.EventType)
	log.Printf("  Length: %d", e.Length)
	log.Printf("  EventData: %v", e.EventData)
}
