package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

// Attribute: Device P-Alarm List
// Attribute: Device T-Alarm List
type DevAlarmList struct {
	Count  uint16          `json:"count"`
	Length uint16          `json:"length"`
	Entry  []DevAlarmEntry `json:"entries"`
}

func (d *DevAlarmList) Size() uint16 {
	size := uint16(4)
	for i := range d.Entry {
		size += d.Entry[i].Size()
	}
	return size
}

func (d *DevAlarmList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, d.Count); err != nil {
		return nil, fmt.Errorf("failed to marshal Count: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}
	for i := range d.Entry {
		entryBytes, err := d.Entry[i].MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Entry[%d]: %w", i, err)
		}
		buf.Write(entryBytes)
	}
	return buf.Bytes(), nil
}

func (d *DevAlarmList) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &d.Count); err != nil {
		return fmt.Errorf("failed to unmarshal Count: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}
	d.Entry = make([]DevAlarmEntry, d.Count)
	for i := 0; i < int(d.Count); i++ {
		if err := d.Entry[i].UnmarshalBinary(r); err != nil {
			return fmt.Errorf("failed to unmarshal Entry[%d]: %w", i, err)
		}
	}
	return nil
}

func (d *DevAlarmList) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<DevAlarmList>", indent)
	log.Printf("%s  Count: %d", indent, d.Count)
	log.Printf("%s  Length: %d", indent, d.Length)
	for i := range d.Entry {
		log.Printf("%s  Entry[%d]:", indent, i)
		d.Entry[i].ShowInfo(indentationLevel + 1)
	}
}
