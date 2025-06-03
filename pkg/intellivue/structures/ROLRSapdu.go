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

type ROLRSapdu struct {
	LinkedID    RorlsId
	InvokeID    uint16
	CommandType CMDType
	Length      uint16
}

func (r *ROLRSapdu) Size() uint16 {
	return 6 + r.LinkedID.Size()
}

func (r *ROLRSapdu) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	LinkedIdData, err := r.LinkedID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(LinkedIdData)

	if err := binary.Write(buf, binary.BigEndian, r.InvokeID); err != nil {
		return nil, fmt.Errorf("ошибка записи InvokeID: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, r.CommandType); err != nil {
		return nil, fmt.Errorf("ошибка записи CommandType: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}

func (r *ROLRSapdu) UnmarshalBinary(reader io.Reader) error {
	if err := r.LinkedID.UnmarshalBinary(reader); err != nil {
		return fmt.Errorf("failed to unmarshal LinkedID: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &r.InvokeID); err != nil {
		return fmt.Errorf("ошибка чтения InvokeID: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &r.CommandType); err != nil {
		return fmt.Errorf("ошибка чтения CommandType: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &r.Length); err != nil {
		return fmt.Errorf("ошибка чтения Length: %w", err)
	}

	return nil
}

func (r *ROLRSapdu) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	r.LinkedID.ShowInfo(mu, indentationLevel+1)
	mu.Lock()
	log.Printf("%s<ROLRSapdu>", indent)
	log.Printf("%s  InvokeID: %d", indent, r.InvokeID)
	log.Printf("%s  CommandType: %d", indent, r.CommandType)
	log.Printf("%s  Length: %d", indent, r.Length)
	mu.Unlock()
}
