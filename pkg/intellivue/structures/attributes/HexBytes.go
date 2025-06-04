package attributes

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type HexBytes []byte

func (h HexBytes) MarshalJSON() ([]byte, error) {
	if len(h) == 0 {
		return []byte("\"\""), nil
	}
	result := make([]string, len(h))
	for i, b := range h {
		result[i] = fmt.Sprintf("%02X", b)
	}
	hexStr := "\"" + strings.Join(result, " ") + "\""
	return []byte(hexStr), nil
}

func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ReplaceAll(s, " ", "")
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*h = b
	return nil
}

func (h HexBytes) Size() uint16 {
	return uint16(len(h))
}

func (h *HexBytes) MarshalBinary() ([]byte, error) {
	return *h, nil
}

func (h *HexBytes) UnmarshalBinary(r io.Reader) error {
	buf := make([]byte, h.Size())
	if len(buf) > 0 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
	}
	*h = buf
	return nil
}

func (h *HexBytes) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<HexBytes> % X", indent, *h)
}
