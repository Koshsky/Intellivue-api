package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type ActionResult struct {
	ManagedObject base.ManagedObjectId `json:"managed_object"`
	ActionType    base.OIDType         `json:"action_type"`
	Length        uint16               `json:"length"`
}

func (a *ActionResult) Size() uint16 {
	return a.ManagedObject.Size() + 4
}

func (a *ActionResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	managedObjectBytes, err := a.ManagedObject.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ManagedObject: %w", err)
	}
	if _, err := buf.Write(managedObjectBytes); err != nil {
		return nil, fmt.Errorf("failed to marshal ManagedObject в буфер: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, a.ActionType); err != nil {
		return nil, fmt.Errorf("failed to marshal ActionType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, a.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (a *ActionResult) UnmarshalBinary(reader io.Reader) error {
	if err := a.ManagedObject.UnmarshalBinary(reader); err != nil {
		return fmt.Errorf("failed to unmarshal ManagedObject: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.ActionType); err != nil {
		return fmt.Errorf("failed to unmarshal ActionType: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}

func (a *ActionResult) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<ActionResult>", indent)
	a.ManagedObject.ShowInfo(indentationLevel + 1)
	log.Printf("%s  ActionType: 0x%04X", indent, a.ActionType)
	log.Printf("%s  Length: %d", indent, a.Length)
}
