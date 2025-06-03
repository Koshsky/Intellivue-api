package structures

import (
	"fmt"
	"io"
)

type ASNLength uint32

func (l ASNLength) MarshalBinary() ([]byte, error) {
	if l <= 127 {
		return []byte{byte(l)}, nil
	}

	value := uint32(l)
	numBytes := l.Size()

	if numBytes < 2 {
		numBytes = 2
	}

	if numBytes > 4 {
		return nil, fmt.Errorf("length %d exceeds maximum supported value", l)
	}

	result := make([]byte, numBytes+1)
	result[0] = byte(0x80 | numBytes)

	for i := numBytes; i > 0; i-- {
		result[i] = byte(value & 0xff)
		value >>= 8
	}

	return result, nil
}

func (l ASNLength) Size() uint16 {
	if l == 0 {
		return 1
	}

	value := uint32(l)
	numBytes := uint16(0)
	for value > 0 {
		numBytes++
		value >>= 8
	}
	return numBytes
}

func (l *ASNLength) UnmarshalBinary(r io.Reader) error {
	var firstByte [1]byte
	if _, err := r.Read(firstByte[:]); err != nil {
		return fmt.Errorf("failed to unmarshal ASNLength first byte: %w", err)
	}
	if firstByte[0] <= 127 {
		*l = ASNLength(firstByte[0])
		return nil
	}
	lengthBytes := int(firstByte[0] & 0x7F)
	if lengthBytes == 0 || lengthBytes > 4 {
		return fmt.Errorf("invalid ASNLength length byte: %d", lengthBytes)
	}
	buf := make([]byte, lengthBytes)
	if _, err := io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("failed to unmarshal ASNLength value bytes: %w", err)
	}
	var value uint32
	for _, b := range buf {
		value = (value << 8) | uint32(b)
	}
	*l = ASNLength(value)
	return nil
}
