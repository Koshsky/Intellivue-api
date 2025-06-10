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

type AttributeList struct {
	Count  uint16    `json:"count"`
	Length uint16    `json:"length"`
	Value  []AVAType `json:"value"`
}

func (a *AttributeList) Size() uint16 {
	return 4 + a.Length
}

func (a *AttributeList) Append(attrbiuteID base.OIDType, ava AttributeValue) {
	a.Value = append(a.Value, AVAType{
		AttributeID: attrbiuteID,
		Length:      ava.Size(),
		Value:       ava,
	})
	a.Count += 1
	a.Length += ava.Size() + 4
}

func (a *AttributeList) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.Count); err != nil {
		return nil, fmt.Errorf("failed to marshal Count: %w", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, a.Length); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	for _, ava := range a.Value {
		avaData, err := ava.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal AVAType: %w", err)
		}
		buf.Write(avaData)
	}

	return buf.Bytes(), nil
}

func (a *AttributeList) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.Count); err != nil {
		return fmt.Errorf("failed to unmarshal attribute count: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("failed to unmarshal attributes data length: %w", err)
	}

	a.Value = make([]AVAType, a.Count)
	for i := uint16(0); i < a.Count; i++ {
		if err := a.Value[i].UnmarshalBinary(r); err != nil {
			return fmt.Errorf("unmarshal erorr AVAType[%d]: %w", i, err)
		}
	}

	return nil
}

func (a *AttributeList) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<AttributeList>", indent)
	log.Printf("%s  Count: %d", indent, a.Count)
	log.Printf("%s  Length: %d", indent, a.Length)

	for _, ava := range a.Value {
		ava.ShowInfo(indentationLevel + 1)
	}
}

func (a *AttributeList) MarshalJSON() ([]byte, error) {
	attrsMap := make(map[string]interface{})

	for _, attr := range a.Value {
		if _, ok := attr.Value.(*HexBytes); ok {
			continue
		}

		valueBytes, err := json.Marshal(attr.Value)
		if err != nil {
			continue
		}

		var valueMap map[string]interface{}
		if err := json.Unmarshal(valueBytes, &valueMap); err != nil {
			continue
		}

		for k, v := range valueMap {
			attrsMap[k] = v
		}
	}

	return json.Marshal(attrsMap)
}
