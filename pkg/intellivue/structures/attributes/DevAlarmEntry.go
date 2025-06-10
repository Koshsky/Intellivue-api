package attributes

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

const (
	BEDSIDE_AUDIBLE        = 0x4000
	CENTRAL_AUDIBLE        = 0x2000
	VISUAL_LATCHING        = 0x1000
	AUDIBLE_LATCHING       = 0x0800
	SHORT_YELLOW_EXTENSION = 0x0400
	DERIVED                = 0x0200

	AL_INHIBITED      base.AlertState = 0x8000
	AL_SUSPENDED      base.AlertState = 0x4000
	AL_LATCHED        base.AlertState = 0x2000
	AL_SILENCED_RESET base.AlertState = 0x1000

	AL_DEV_IN_TEST_MODE base.AlertState = 0x0400
	AL_DEV_IN_STANDBY   base.AlertState = 0x0200
	AL_DEV_IN_DEMO_MODE base.AlertState = 0x0100

	NO_ALERT     base.AlertType = 0x0000
	LOW_PRI_T_AL base.AlertType = 0x0001
	MED_PRI_T_AL base.AlertType = 0x0002
	HI_PRI_T_AL  base.AlertType = 0x0004
	LOW_PRI_P_AL base.AlertType = 0x0100
	MED_PRI_P_AL base.AlertType = 0x0200
	HI_PRI_P_AL  base.AlertType = 0x0400
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

func alertTypeToString(alertType base.AlertType) string {
	switch alertType {
	case NO_ALERT:
		return "no_alert"
	case LOW_PRI_T_AL:
		return "low_pri_t_al"
	case MED_PRI_T_AL:
		return "med_pri_t_al"
	case HI_PRI_T_AL:
		return "hi_pri_t_al"
	case LOW_PRI_P_AL:
		return "low_pri_p_al"
	case MED_PRI_P_AL:
		return "med_pri_p_al"
	case HI_PRI_P_AL:
		return "hi_pri_p_al"
	default:
		return fmt.Sprintf("unknown_alert_type_0x%04X", uint16(alertType))
	}
}
func alertStateToString(alertState base.AlertState) string {
	switch alertState & 0xF000 {
	case AL_INHIBITED:
		return "inhibited"
	case AL_SUSPENDED:
		return "suspended"
	case AL_LATCHED:
		return "latched"
	case AL_SILENCED_RESET:
		return "silenced_reset"
	default:
		return fmt.Sprintf("unknown_alert_state_0x%04X", uint16(alertState&0xF000))
	}
}
func devModeToString(alertState base.AlertState) string {
	switch alertState & 0x0F00 {
	case AL_DEV_IN_TEST_MODE:
		return "test"
	case AL_DEV_IN_STANDBY:
		return "standby"
	case AL_DEV_IN_DEMO_MODE:
		return "demo"
	default:
		return fmt.Sprintf("unknown_dev_mode_0x%04X", uint16(alertState&0x0F00))
	}
}

func (d *DevAlarmEntry) MarshalJSON() ([]byte, error) {
	jsonData := map[string]interface{}{
		// "alert_source"
		// "alert_code"
		"alert_type":  d.AlertType,
		"alert_state": d.AlertState,
		// "object"
		// "alert_info_id"
		// "priority"
		"flags": d.Info.Flags,
		"value": d.Info.String.Value,
	}
	return json.Marshal(jsonData)
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
