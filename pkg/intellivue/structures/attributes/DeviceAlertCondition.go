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

type DeviceAlertCondition struct {
	DeviceAlertState base.AlertState `json:"device_alert_state"`
	AlertStateChgCnt uint16          `json:"al_stat_chg_cnt"`
	MaxPAlarm        base.AlertType  `json:"max_p_alarm"`
	MaxTAlarm        base.AlertType  `json:"max_t_alarm"`
	MaxAudAlarm      base.AlertType  `json:"max_aud_alarm"`
}

func (d *DeviceAlertCondition) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, d.DeviceAlertState); err != nil {
		return nil, fmt.Errorf("failed to marshal DeviceAlertState: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.AlertStateChgCnt); err != nil {
		return nil, fmt.Errorf("failed to marshal AlertStateChgCnt: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.MaxPAlarm); err != nil {
		return nil, fmt.Errorf("failed to marshal MaxPAlarm: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.MaxTAlarm); err != nil {
		return nil, fmt.Errorf("failed to marshal MaxTAlarm: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, d.MaxAudAlarm); err != nil {
		return nil, fmt.Errorf("failed to marshal MaxAudAlarm: %w", err)
	}

	return buf.Bytes(), nil
}

func (d *DeviceAlertCondition) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &d.DeviceAlertState); err != nil {
		return fmt.Errorf("failed to unmarshal DeviceAlertState: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.AlertStateChgCnt); err != nil {
		return fmt.Errorf("failed to unmarshal AlertStateChgCnt: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.MaxPAlarm); err != nil {
		return fmt.Errorf("failed to unmarshal MaxPAlarm: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.MaxTAlarm); err != nil {
		return fmt.Errorf("failed to unmarshal MaxTAlarm: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &d.MaxAudAlarm); err != nil {
		return fmt.Errorf("failed to unmarshal MaxAudAlarm: %w", err)
	}

	return nil
}

func (d *DeviceAlertCondition) Size() uint16 {
	return 10
}

func (d *DeviceAlertCondition) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%sDeviceAlertCondition:", indent)
	log.Printf("%s  DeviceAlertState: %#04x", indent, d.DeviceAlertState)
	log.Printf("%s  AlertStateChgCnt: %#04x", indent, d.AlertStateChgCnt)
	log.Printf("%s  MaxPAlarm: %#04x", indent, d.MaxPAlarm)
	log.Printf("%s  MaxTAlarm: %#04x", indent, d.MaxTAlarm)
	log.Printf("%s  MaxAudAlarm: %#04x", indent, d.MaxAudAlarm)
}
