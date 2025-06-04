package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type EventReportResult struct {
	ManagedObject base.ManagedObjectId `json:"managed_object"`
	CurrentTime   base.RelativeTime    `json:"current_time"`
	EventType     base.OIDType         `json:"event_type"`
	Length        uint16               `json:"length"`
}

func (e *EventReportResult) Size() uint16 {
	return e.ManagedObject.Size() + 8
}

func (e *EventReportResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	moData, err := e.ManagedObject.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ManagedObject: %w", err)
	}
	buf.Write(moData)

	if err := binary.Write(buf, binary.BigEndian, e.CurrentTime); err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentTime: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, e.EventType); err != nil {
		return nil, fmt.Errorf("failed to marshal EventType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, e.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (e *EventReportResult) UnmarshalBinary(r io.Reader) error {
	if e == nil {
		return fmt.Errorf("nil EventReportResult receiver")
	}
	if err := e.ManagedObject.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ManagedObject: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &e.CurrentTime); err != nil {
		return fmt.Errorf("failed to unmarshal CurrentTime: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &e.EventType); err != nil {
		return fmt.Errorf("failed to unmarshal EventType: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &e.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}
