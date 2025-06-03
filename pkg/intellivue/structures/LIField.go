package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type LIField uint16

func (l LIField) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if l < 255 {
		buf.WriteByte(byte(l))
	} else {
		buf.WriteByte(0xFF)
		binary.Write(&buf, binary.BigEndian, uint16(l))
	}

	return buf.Bytes(), nil
}

func writeLIField(buf *bytes.Buffer, val LIField) {
	data, _ := val.MarshalBinary()
	buf.Write(data)
}

func (l LIField) Size() uint16 {
	if l < 255 {
		return 1
	} else {
		return 3
	}
}

func (l *LIField) UnmarshalBinary(r io.Reader) error {
	var firstByte [1]byte
	if _, err := r.Read(firstByte[:]); err != nil {
		return fmt.Errorf("failed to unmarshal LIField first byte: %w", err)
	}
	if firstByte[0] == 0xFF {
		var val uint16
		if err := binary.Read(r, binary.BigEndian, &val); err != nil {
			return fmt.Errorf("failed to unmarshal LIField uint16: %w", err)
		}
		*l = LIField(val)
	} else {
		*l = LIField(firstByte[0])
	}
	return nil
}
