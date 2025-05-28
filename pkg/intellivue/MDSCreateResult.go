package intellivue

import (
	"bytes"
	"encoding/binary"
)

const (
	RORS_APDU = 0x0002
	CMD_CONFIRMED_EVENT_REPORT = 0x0001
	NOM_MOC_VMS_MDS = 0x0021
	NOM_NOTI_MDS_CREAT = 0x0d06
)

type MDSCreateResult struct {
	SessionID   uint16
	PContextID  uint16
	ROType      uint16
	Length1     uint16
	InvokeID    uint16
	CommandType uint16
	Length2     uint16
	MObjClass   uint16
	ContextID   uint16
	Handle      uint16
	EventTime   uint32
	EventType   uint16
	Length3     uint16
}

func (m *MDSCreateResult) Length() uint16 {
	return 28
}

func (m *MDSCreateResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, m.SessionID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.PContextID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.ROType); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.Length1); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.InvokeID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.CommandType); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.Length2); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.MObjClass); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.ContextID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.Handle); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.EventTime); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.EventType); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, m.Length3); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}



func NewMDSCreateResult() *MDSCreateResult {
	// Define the necessary fields based on the provided data
	return &MDSCreateResult{
		SessionID:   0xE100,
		PContextID:  2,
		ROType:      RORS_APDU,
		Length1:     20,
		InvokeID:    1,
		CommandType: CMD_CONFIRMED_EVENT_REPORT,
		Length2:     14,
		MObjClass:   NOM_MOC_VMS_MDS,
		ContextID:   0,
		Handle:      0,
		EventTime:   0x00afbd00,  // TODO: определить число
		EventType:   NOM_NOTI_MDS_CREAT,
		Length3:     0,
	}
}
