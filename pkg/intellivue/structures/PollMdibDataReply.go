package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PollMdibDataReply представляет структуру данных ответа на опрос MDS.
type PollMdibDataReply struct {
	PollNumber    uint16
	RelTimeStamp  uint32
	AbsTimeStamp  uint64
	Partition     uint16
	Code          uint16
	PolledAttrGrp uint16
	PollInfoList  PollInfoList // Используем уже созданную структуру
}

// Size возвращает общую длину PollMdibDataReply в байтах.
func (p *PollMdibDataReply) Size() uint16 {
	// Суммируем длины всех полей.
	length := uint16(0)
	length += 2 // PollNumber
	length += 4 // RelTimeStamp
	length += 8 // AbsTimeStamp
	length += 2 // Partition
	length += 2 // Code
	length += 2 // PolledAttrGrp
	length += p.PollInfoList.Size()

	return length
}

// MarshalBinary кодирует структуру PollMdibDataReply в бинарный формат.
func (p *PollMdibDataReply) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.PollNumber); err != nil {
		return nil, fmt.Errorf("ошибка записи PollNumber: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.RelTimeStamp); err != nil {
		return nil, fmt.Errorf("ошибка записи RelTimeStamp: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.AbsTimeStamp); err != nil {
		return nil, fmt.Errorf("ошибка записи AbsTimeStamp: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.Partition); err != nil {
		return nil, fmt.Errorf("ошибка записи Partition: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.Code); err != nil {
		return nil, fmt.Errorf("ошибка записи Code: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, p.PolledAttrGrp); err != nil {
		return nil, fmt.Errorf("ошибка записи PolledAttrGrp: %w", err)
	}

	pollInfoListData, err := p.PollInfoList.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга PollInfoList: %w", err)
	}
	buf.Write(pollInfoListData)

	return buf.Bytes(), nil
}
