package structures

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type UserData struct {
	Length base.ASNLength
	Data   MDSEUserInfoStd
}

func (u *UserData) Size() uint16 {
	u.Length = base.ASNLength(u.Data.Size())
	return u.Length.Size() + u.Data.Size()
}

func (u *UserData) ShowInfo(indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)
	log.Printf("%s<UserData>", indent)
	log.Printf("%s  Length: 0x%X", indent, u.Length)
	u.Data.ShowInfo(indentationLevel + 1)
}

func (u *UserData) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	u.Length = base.ASNLength(u.Data.Size())

	lengthBytes, err := u.Length.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %v", err)
	}
	buf.Write(lengthBytes)

	userData, err := u.Data.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UserData: %v", err)
	}
	buf.Write(userData)

	return buf.Bytes(), nil
}

func (u *UserData) UnmarshalBinary(r io.Reader) error {
	if u == nil {
		return fmt.Errorf("nil UserData receiver")
	}
	if err := u.Length.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal Length: %w", err)
	}
	if err := u.Data.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal UserData: %w", err)
	}
	return nil
}
