package attributes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type NuObsValCmp struct {
	Count  uint16       `json:"count"`
	Length uint16       `json:"length"`
	Value  []NuObsValue `json:"value"`
}

func (a *NuObsValCmp) Size() uint16 {
	return 4 + a.Length // count + length + list
}

func (a *NuObsValCmp) MarshalBinary() ([]byte, error) {
	a.Count = uint16(len(a.Value))
	total := uint16(0)
	for _, ava := range a.Value {
		total += ava.Size()
	}
	a.Length = total

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
			return nil, fmt.Errorf("failed to marshal NuObsValue: %w", err)
		}
		buf.Write(avaData)
	}

	return buf.Bytes(), nil
}

func (a *NuObsValCmp) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.Count); err != nil {
		return fmt.Errorf("failed to unmarshal attribute count: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("failed to unmarshal attributes data length: %w", err)
	}

	a.Value = make([]NuObsValue, a.Count)
	for i := uint16(0); i < a.Count; i++ {
		if err := a.Value[i].UnmarshalBinary(r); err != nil {
			return fmt.Errorf("unmarshal erorr NuObsValue[%d]: %w", i, err)
		}
	}

	return nil
}

func (a *NuObsValCmp) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<NuObsValCmp>", indent)
	log.Printf("%s  Count: %d", indent, a.Count)
	log.Printf("%s  Length: %d", indent, a.Length)

	for _, ava := range a.Value {
		ava.ShowInfo(indentationLevel + 1)
	}
}
