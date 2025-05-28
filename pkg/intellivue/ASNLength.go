package intellivue

import (
	"fmt"
)

type ASNLength uint32

func (l ASNLength) MarshalBinary() ([]byte, error) {
    // Handle short form (length <= 127)
    if l <= 127 {
        return []byte{byte(l)}, nil
    }

    value := uint32(l)
    numBytes := l.Length()

    // ASN.1 requires at least 2 bytes for long form lengths
    if numBytes < 2 {
        numBytes = 2
    }

    // Check for maximum supported length (we're using uint32)
    if numBytes > 4 {
        return nil, fmt.Errorf("length %d exceeds maximum supported value", l)
    }

    result := make([]byte, numBytes+1)
    result[0] = byte(0x80 | numBytes)
    
    // Write bytes in big-endian order
    for i := numBytes; i > 0; i-- {
        result[i] = byte(value & 0xff)
        value >>= 8
    }

    return result, nil
}

func (l ASNLength) Length() uint16 {
    if l == 0 {
        return 1  // Special case for zero length
    }

    value := uint32(l)
    numBytes := uint16(0)
    for value > 0 {
        numBytes++
        value >>= 8
    }
    return numBytes
}

