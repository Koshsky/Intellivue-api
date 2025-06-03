package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

const (
	RORLS_FIRST              uint8 = 0x0001 // set in the first message
	RORLS_NOT_FIRST_NOT_LAST uint8 = 0x0002
	RORLS_LAST               uint8 = 0x0003 // last RORLSapdu, one RORSapdu to follow
)

type RorlsId struct {
	State uint8
	Count uint8 // counter starts with 1
}

func (r *RorlsId) Size() uint16 {
	return 2
}

func (r *RorlsId) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, r.State); err != nil {
		return nil, fmt.Errorf("ошибка записи State: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, r.Count); err != nil {
		return nil, fmt.Errorf("ошибка записи Count: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *RorlsId) UnmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.BigEndian, &r.State); err != nil {
		return fmt.Errorf("ошибка чтения State: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &r.Count); err != nil {
		return fmt.Errorf("ошибка чтения Count: %w", err)
	}

	return nil
}

func (r *RorlsId) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	mu.Lock()
	log.Printf("%s<RorlsId>", indent)
	log.Printf("%s  State: 0x%X", indent, r.State)
	log.Printf("%s  Count: %d", indent, r.Count)
	mu.Unlock()
}
