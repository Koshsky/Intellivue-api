package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

// Command Types
const (
	CMD_EVENT_REPORT           CMDType = 0x0000 // Event Report
	CMD_CONFIRMED_EVENT_REPORT CMDType = 0x0001 // Confirmed Event Report
	CMD_GET                    CMDType = 0x0003 // Get
	CMD_SET                    CMDType = 0x0004 // Set
	CMD_CONFIRMED_SET          CMDType = 0x0005 // Confirmed Set
	CMD_CONFIRMED_ACTION       CMDType = 0x0007 // Confirmed Action
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

func (a *ActionResult) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	mu.Lock()
	log.Printf("%s<ActionResult>", indent)
	mu.Unlock()
	a.ManagedObject.ShowInfo(mu, indentationLevel+1)
	mu.Lock()
	log.Printf("%s  ActionType: 0x%X", indent, a.ActionType)
	log.Printf("%s  Length: %d", indent, a.Length)
	mu.Unlock()
}
