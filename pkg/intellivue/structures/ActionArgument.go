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

type ActionArgument struct {
	ManagedObject base.ManagedObjectId `json:"managed_object"`
	Scope         uint32               `json:"scope"`
	ActionType    base.OIDType         `json:"action_type"`
	Length        uint16               `json:"length"`
}

func (a *ActionArgument) Size() uint16 {
	return a.ManagedObject.Size() + 4 + 2 + 2
}

func (a *ActionArgument) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	objData, err := a.ManagedObject.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ManagedObject: %w", err)
	}
	buf.Write(objData)

	if err := binary.Write(buf, binary.BigEndian, a.Scope); err != nil {
		return nil, fmt.Errorf("failed to marshal Scope: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, a.ActionType); err != nil {
		return nil, fmt.Errorf("failed to marshal ActionType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, a.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (a *ActionArgument) UnmarshalBinary(reader io.Reader) error {
	if err := a.ManagedObject.UnmarshalBinary(reader); err != nil {
		return fmt.Errorf("failed to unmarshal ManagedObject: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Scope); err != nil {
		return fmt.Errorf("failed to unmarshal Scope: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.ActionType); err != nil {
		return fmt.Errorf("failed to unmarshal ActionType: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	return nil
}

func (a *ActionArgument) ShowInfo(indentationLevel int) {
	indent := strings.Repeat(" ", indentationLevel*2)
	log.Printf("%s<ActionArgument>\n", indent)
	log.Printf("%s  <ManagedObject>\n", indent)
	log.Printf("%s    MObjClass: 0x%04X\n", indent, a.ManagedObject.MObjClass)
	log.Printf("%s  </ManagedObject>\n", indent)
	log.Printf("%s  Scope: %d\n", indent, a.Scope)
	log.Printf("%s  ActionType: 0x%04X\n", indent, a.ActionType)
	log.Printf("%s  Length: %d\n", indent, a.Length)
}
