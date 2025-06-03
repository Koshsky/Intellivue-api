package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	NOM_ACT_POLL_MDIB_DATA     OIDType = 0x0c16
	NOM_ACT_POLL_MDIB_DATA_EXT OIDType = 0xf13b
)

type ActionArgument struct {
	ManagedObject ManagedObjectId
	Scope         uint32  //fixed value 0
	ActionType    OIDType // identification of method
	Length        uint16
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
