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

type AttributeValue interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(io.Reader) error
	Size() uint16
	ShowInfo(indentationLevel int)
}

type AVAType struct {
	AttributeID base.OIDType   `json:"attribute_id"`
	Length      uint16         `json:"length"`
	Value       AttributeValue `json:"value"`
}

func (a *AVAType) Size() uint16 {
	return 4 + a.Value.Size()
}

func (a *AVAType) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.AttributeID); err != nil {
		return nil, fmt.Errorf("failed to marshal AttributeID: %w", err)
	}

	valBytes, err := a.Value.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Value: %w", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, uint16(len(valBytes))); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}
	buf.Write(valBytes)

	return buf.Bytes(), nil
}

func (a *AVAType) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.AttributeID); err != nil {
		return fmt.Errorf("failed to unmarshal AttributeID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	var val AttributeValue
	switch a.AttributeID {
	case base.NOM_ATTR_AL_MON_P_AL_LIST, base.NOM_ATTR_AL_MON_T_AL_LIST:
		list := &DevAlarmList{}
		if err := list.UnmarshalBinary(io.LimitReader(r, int64(a.Length))); err != nil {
			return fmt.Errorf("failed to unmarshal DevAlarmList: %w", err)
		}
		val = list
	case base.NOM_ATTR_DEV_AL_COND:
		condition := &DeviceAlertCondition{}
		if err := condition.UnmarshalBinary(io.LimitReader(r, int64(a.Length))); err != nil {
			return fmt.Errorf("failed to unmarshal DevAlarmList: %w", err)
		}
		val = condition
	case base.NOM_ATTR_ID_TYPE:
		obj := &TYPE{}
		if err := obj.UnmarshalBinary(io.LimitReader(r, int64(a.Length))); err != nil {
			return fmt.Errorf("failed to unmarshal DevAlarmList: %w", err)
		}
		val = obj
	default:
		hb := make(HexBytes, a.Length)
		if a.Length > 0 {
			if _, err := io.ReadFull(r, hb); err != nil {
				return fmt.Errorf("failed to unmarshal HexBytes: %w", err)
			}
		}
		val = &hb
	}
	a.Value = val
	return nil
}

func (a *AVAType) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<AVAType>", indent)
	log.Printf("%s  AttributeID: %#04x", indent, a.AttributeID)
	log.Printf("%s  Length: %d", indent, a.Length)
	a.Value.ShowInfo(indentationLevel + 1)
}
