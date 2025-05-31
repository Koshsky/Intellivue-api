package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// The Event Report Result command is used as a response message to the Event Report message. It is
// appended to the Operation Result message with the command_type
// CMD_CONFIRMED_EVENT_REPORT.
type EventReportResult struct {
	ManagedObject ManagedObjectId // mirrored from EvRep
	CurrentTime   RelativeTime    // result time stamp
	EventType     OIDType         // identification of event
	Length        uint16          // size of appended data
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
