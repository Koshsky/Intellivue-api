package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	CMD_EVENT_REPORT           CMDType = 0x0000
	CMD_CONFIRMED_EVENT_REPORT CMDType = 0x0001
	CMD_GET                    CMDType = 0x0003
	CMD_SET                    CMDType = 0x0004
	CMD_CONFIRMED_SET          CMDType = 0x0005
	CMD_CONFIRMED_ACTION       CMDType = 0x0007
)

type ActionResult struct {
	ManagedObject ManagedObjectId
	ActionType    OIDType
	Length        uint16
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
	log.Printf("%s  ActionType: 0x%X", indent, a.ActionType)
	log.Printf("%s  Length: %d", indent, a.Length)
}
