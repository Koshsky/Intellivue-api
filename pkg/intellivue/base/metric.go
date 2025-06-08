package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
)

type MetricCategory uint16

type MetricSpec uint16

type MetricAccess uint16

type MetricRelevance uint16

type MetricModality uint16

type MetricState uint16

type MetricStructure struct {
	MSStruct uint8 `json:"ms_struct"`  // describes the structure of the object, 0 means simple, 1 means compound object.
	MSCompNo uint8 `json:"ms_comp_no"` // contains the maximum number of components in the compound
}

func (mc MetricStructure) Size() uint16 {
	return 2
}

func (mc *MetricStructure) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, mc.MSStruct); err != nil {
		return nil, fmt.Errorf("failed to marshal MSStruct: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, mc.MSCompNo); err != nil {
		return nil, fmt.Errorf("failed to marshal MSCompNo: %w", err)
	}

	return buf.Bytes(), nil
}

func (mc *MetricStructure) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &mc.MSStruct); err != nil {
		return fmt.Errorf("failed to unmarshal MSStruct: %w", err)
	}
	if err := binary.Read(reader, binary.BigEndian, &mc.MSCompNo); err != nil {
		return fmt.Errorf("failed to unmarshal MSCompNo: %w", err)
	}

	return nil
}

func (mc *MetricStructure) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<GlbHandle>", indent)
	log.Printf("%s MSStruct: 0x%X", indent, mc.MSStruct)
	log.Printf("%s MSCompNo: %d", indent, mc.MSCompNo)
}
