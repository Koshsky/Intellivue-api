package packages

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

// The Single Poll Data Result message contains a PollInfoList which is of variable length. The length
// fields in the message depend on the length of the PollInfoList.
type SinglePollDataResultLinked struct {
	SPpdu
	ROapdus
	ROLRSapdu
	ActionResult
	PollMdibDataReply
}

func (m *SinglePollDataResultLinked) Size() uint16 {
	return m.SPpdu.Size() + m.ROapdus.Size() + m.ROLRSapdu.Size() + m.ActionResult.Size() // + m.PollMdibDataReply.Size()
}

func (m *SinglePollDataResultLinked) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	SPDuData, err := m.SPpdu.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(SPDuData)

	ROapdusData, err := m.ROapdus.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(ROapdusData)

	RORSapduData, err := m.ROLRSapdu.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(RORSapduData)

	ActionResultData, err := m.ActionResult.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(ActionResultData)

	PollMdibDataReplyData, err := m.PollMdibDataReply.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(PollMdibDataReplyData)

	return buf.Bytes(), nil
}

func (m *SinglePollDataResultLinked) UnmarshalBinary(r io.Reader) error {
	if m == nil {
		return fmt.Errorf("nil SinglePollDataResult receiver")
	}

	if err := m.SPpdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal SPpdu: %w", err)
	}
	if err := m.ROapdus.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ROapdus: %w", err)
	}
	if err := m.ROLRSapdu.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal RORSapdu: %w", err)
	}
	if err := m.ActionResult.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal ActionResult: %w", err)
	}
	if err := m.PollMdibDataReply.UnmarshalBinary(r); err != nil {
		return fmt.Errorf("failed to unmarshal PollMdibDataReply: %w", err)
	}

	return nil
}

func (m *SinglePollDataResultLinked) ShowInfo(mu *sync.Mutex, indentationLevel int) {
	indent := strings.Repeat("  ", indentationLevel)

	mu.Lock()
	log.Printf("%s<SinglePollDataResultLinked>", indent)
	m.SPpdu.ShowInfo(indentationLevel + 1)
	m.ROapdus.ShowInfo(indentationLevel + 1)
	m.ROLRSapdu.ShowInfo(indentationLevel + 1)
	m.ActionResult.ShowInfo(indentationLevel + 1)
	m.PollMdibDataReply.ShowInfo(indentationLevel + 1)
	mu.Unlock()
}
