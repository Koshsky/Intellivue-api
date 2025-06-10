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

type NuObsValue struct {
	PhysioID base.OIDType     `json:"physio_id"`
	State    base.MeasureMode `json:"state"`
	UnitCode base.OIDType     `json:"-"`
	Value    base.FLOATType   `json:"value"`
}

func (n *NuObsValue) MarshalJSON() ([]byte, error) {
	unit := base.UnitCodes[n.UnitCode]
	return json.Marshal(map[string]interface{}{
		"physio_id": n.PhysioID,
		"state":     n.State,
		"value":     n.Value,
		"unit":      unit,
	})
}

func (n *NuObsValue) Size() uint16 {
	return 10
}

func (n *NuObsValue) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, n.PhysioID); err != nil {
		return nil, fmt.Errorf("failed to marshal PhysioID: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, n.State); err != nil {
		return nil, fmt.Errorf("failed to marshal State: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, n.UnitCode); err != nil {
		return nil, fmt.Errorf("failed to marshal UnitCode: %w", err)
	}

	valueBytes, err := n.Value.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Value: %w", err)
	}
	buf.Write(valueBytes)

	return buf.Bytes(), nil
}

func (n *NuObsValue) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &n.PhysioID); err != nil {
		return fmt.Errorf("failed to unmarshal PhysioID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &n.State); err != nil {
		return fmt.Errorf("failed to unmarshal State: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &n.UnitCode); err != nil {
		return fmt.Errorf("failed to unmarshal UnitCode: %w", err)
	}

	if err := n.Value.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal Value: %w", err)
	}

	return nil
}

func (h *NuObsValue) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<NuObsValue>\n", indent)
	log.Printf("%s  PhysioID: %d\n", indent, h.PhysioID)
	log.Printf("%s  State: %d\n", indent, h.State)
	log.Printf("%s  UnitCode: %d\n", indent, h.UnitCode)
	log.Printf("%s  Value: %s\n", indent, h.Value.String())
}
