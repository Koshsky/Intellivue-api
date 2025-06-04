package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type DevAlarmEntry struct {
	AlertSource base.OIDType         `json:"al_source"`
	AlertCode   base.OIDType         `json:"al_code"`
	AlertType   base.AlertType       `json:"al_type"`
	AlertState  base.AlertState      `json:"al_state"`
	Object      base.ManagedObjectId `json:"object"`
	AlertInfoID base.PrivateOID      `json:"alert_info_id"`
	Length      uint16               `json:"length"`
	Info        StrAlMonInfo         `json:"info"`
}

func (d *DevAlarmEntry) Size() uint16 {
	return 12 + d.Object.Size() + d.Info.Size()
}

func (d *DevAlarmEntry) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, d.AlertSource); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertSource: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.AlertCode); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertCode: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.AlertType); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertType: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.AlertState); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertState: %w", err)
	}
	objBytes, err := d.Object.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Object: %w", err)
	}
	buf.Write(objBytes)
	if err := binary.Write(buf, binary.BigEndian, d.AlertInfoID); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertInfoID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}
	infoBytes, err := d.Info.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Info: %w", err)
	}
	buf.Write(infoBytes)
	return buf.Bytes(), nil
}

func (d *DevAlarmEntry) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &d.AlertSource); err != nil {
		return fmt.Errorf("failed to unmarshal AlertSource: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.AlertCode); err != nil {
		return fmt.Errorf("failed to unmarshal AlertCode: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.AlertType); err != nil {
		return fmt.Errorf("failed to unmarshal AlertType: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.AlertState); err != nil {
		return fmt.Errorf("failed to unmarshal AlertState: %w", err)
	}
	if err := d.Object.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal Object: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.AlertInfoID); err != nil {
		return fmt.Errorf("failed to unmarshal AlertInfoID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}
	if err := d.Info.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal Info: %w", err)
	}
	return nil
}

func (d *DevAlarmEntry) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<DevAlarmEntry>", indent)
	log.Printf("%s  AlertSource: %#04x", indent, d.AlertSource)
	log.Printf("%s  AlertCode: %#04x", indent, d.AlertCode)
	log.Printf("%s  AlertType: %#04x", indent, d.AlertType)
	log.Printf("%s  AlertState: %#04x", indent, d.AlertState)
	d.Object.ShowInfo(indentationLevel + 1)
	log.Printf("%s  AlertInfoID: %#04x", indent, d.AlertInfoID)
	log.Printf("%s  Length: %d", indent, d.Length)
	d.Info.ShowInfo(indentationLevel + 1)
}
