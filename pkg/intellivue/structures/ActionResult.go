package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("ошибка MarshalBinary для ManagedObject: %w", err)
	}
	if _, err := buf.Write(managedObjectBytes); err != nil {
		return nil, fmt.Errorf("ошибка записи ManagedObject в буфер: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.ActionType); err != nil {
		return nil, fmt.Errorf("ошибка записи ActionType: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (a *ActionResult) UnmarshalBinary(reader io.Reader) error {
	if err := a.ManagedObject.UnmarshalBinary(reader); err != nil {
		return fmt.Errorf("ошибка UnmarshalBinary для ManagedObject: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &a.ActionType); err != nil {
		return fmt.Errorf("ошибка чтения ActionType: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("ошибка чтения Length: %w", err)
	}

	return nil
}
