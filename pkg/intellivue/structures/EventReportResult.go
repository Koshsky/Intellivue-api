package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type EventReportResult struct {
	ManagedObject ManagedObjectId
	CurrentTime   RelativeTime
	EventType     uint16
	Length        uint16
}

func (e *EventReportResult) Size() uint16 {
	return e.ManagedObject.Size() + 4 + 2 + 2
}

func (e *EventReportResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	moData, err := e.ManagedObject.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ManagedObject: %w", err)
	}
	if _, err := buf.Write(moData); err != nil {
		return nil, fmt.Errorf("failed to write ManagedObject: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, e.CurrentTime); err != nil {
		return nil, fmt.Errorf("failed to write CurrentTime: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, e.EventType); err != nil {
		return nil, fmt.Errorf("failed to write EventType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, e.Length); err != nil {
		return nil, fmt.Errorf("failed to write Length: %w", err)
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
		return fmt.Errorf("failed to read CurrentTime: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &e.EventType); err != nil {
		return fmt.Errorf("failed to read EventType: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &e.Length); err != nil {
		return fmt.Errorf("failed to read Length: %w", err)
	}

	return nil
}
