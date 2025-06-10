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

type MetricSpec struct {
	UpdatePeriod base.RelativeTime    `json:"update_period"`
	Category     base.MetricCategory  `json:"category"`
	Access       base.MetricAccess    `json:"access"`
	Structure    base.MetricStructure `json:"structure"`
	Relevance    base.MetricRelevance `json:"relevance"`
}

func (ms *MetricSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"update_period": ms.UpdatePeriod,
		"category":      ms.Category,
	})
}

func (ms *MetricSpec) Size() uint16 {
	return 10 + ms.Structure.Size()
}

func (ms *MetricSpec) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, ms.UpdatePeriod); err != nil {
		return nil, fmt.Errorf("failed to marshal UpdatePeriod: %w", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, ms.Category); err != nil {
		return nil, fmt.Errorf("failed to marshal Category: %w", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, ms.Access); err != nil {
		return nil, fmt.Errorf("failed to marshal Access: %w", err)
	}
	structureData, err := ms.Structure.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Structure: %w", err)
	}
	buf.Write(structureData)
	if err := binary.Write(&buf, binary.BigEndian, ms.Relevance); err != nil {
		return nil, fmt.Errorf("failed to marshal Relevance: %w", err)
	}

	return buf.Bytes(), nil
}

func (ms *MetricSpec) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &ms.UpdatePeriod); err != nil {
		return fmt.Errorf("failed to unmarshal attribute UpdatePeriod: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &ms.Category); err != nil {
		return fmt.Errorf("failed to unmarshal Category: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &ms.Access); err != nil {
		return fmt.Errorf("failed to unmarshal Access: %w", err)
	}
	if err := ms.Structure.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to umarshal Structure: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &ms.Relevance); err != nil {
		return fmt.Errorf("failed to unmarshal Relevance: %w", err)
	}

	return nil
}

func (ms *MetricSpec) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<MetricSpec>", indent)
	log.Printf("%s  UpdatePeriod: %d", indent, ms.UpdatePeriod)
	log.Printf("%s  Category: %d", indent, ms.Category)
	log.Printf("%s  Access: %d", indent, ms.Access)
	ms.Structure.ShowInfo(indentationLevel + 1)
	log.Printf("%s  Relevance: %d", indent, ms.Relevance)
}
