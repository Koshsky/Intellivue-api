package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

type AVAType struct {
	AttributeID OIDType
	Length      uint16
	Value       interface{}
}

func (a *AVAType) Size() uint16 { // TODO: это работает корректно???
	valueSize := uint16(0)
	if sizedValue, ok := a.Value.(interface{ Size() uint16 }); ok {
		valueSize = sizedValue.Size()
	}
	return 4 + valueSize
}

func (a *AVAType) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, a.AttributeID); err != nil {
		return nil, fmt.Errorf("failed to marshal AttributeID: %w", err)
	}

	// TODO: можно проще получить размер a.Value в байтах?
	valueSize := uint16(0)
	if sizedValue, ok := a.Value.(interface{ Size() uint16 }); ok {
		valueSize = sizedValue.Size()
	}
	if err := binary.Write(&buf, binary.BigEndian, valueSize); err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %w", err)
	}

	if marshalerValue, ok := a.Value.(interface{ MarshalBinary() ([]byte, error) }); ok {
		data, err := marshalerValue.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Value: %w", err)
		}
		buf.Write(data)
	}

	return buf.Bytes(), nil
}

func (a *AVAType) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &a.AttributeID); err != nil {
		return fmt.Errorf("failed to unmarshal AttributeID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &a.Length); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}

	valueBytes := make([]byte, a.Length)
	if a.Length > 0 {
		limitedReader := io.LimitReader(r, int64(a.Length))
		if _, err := io.ReadFull(limitedReader, valueBytes); err != nil {
			return fmt.Errorf("failed to unmarshal Value: %w", err)
		}
	}
	a.Value = valueBytes

	return nil
}

func (a *AVAType) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	log.Printf("%s<AttributeList>", indent)
	log.Printf("%s  AttributeID: %#04x", indent, a.AttributeID)
	log.Printf("%s  Length: %d", indent, a.Length)
	if valueBytes, ok := a.Value.([]byte); ok {
		hexStr := utils.PrintHex(valueBytes)
		indentedHexStr := strings.Replace(hexStr, "\n", "\n"+indent+"         ", -1)
		log.Printf("%s  value: %s", indent, indentedHexStr)
	} else {
		log.Printf("%s  value: <not []byte>", indent)
	}
}
