package packages

import (
	"bytes"
	"fmt"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

// MDSPollResult представляет структуру данных для результата опроса,
// состоящую из нескольких логических частей.
type MDSPollResult struct {
	SPpdu
	ROapdus
	ROIVapdu
	// ActionResult - В MDS Poll Result это не Action Argument, а скорее Managed Object ID и связанные поля.
	ManagedObjectId ManagedObjectId
	// Поля, специфичные для ответа на опрос MDS
	PollMdibDataReply PollMdibDataReply
}

// Size возвращает общую длину MDSPollResult в байтах.
func (m *MDSPollResult) Size() uint16 {
	return m.SPpdu.Size() + m.ROapdus.Size() + m.ROIVapdu.Size() + m.ManagedObjectId.Size() + m.PollMdibDataReply.Size()
}

// MarshalBinary кодирует структуру MDSPollResult в бинарный формат.
func (m *MDSPollResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Обновляем поля длины во вложенных структурах перед маршалингом.

	// Длина PollMdibDataReply фиксирована.
	// В MDS Poll Result нет поля Length3, но PollMdibDataReply содержит PollInfoList, у которого есть Length.
	// Это поле Length в PollInfoList должно отражать длину данных в PollInfoList (элементы SingleContextPoll).
	// Значение этого поля устанавливается при маршалинге самой PollInfoList.

	// Длина ROIVapdu: ManagedObjectId + PollMdibDataReply
	m.ROIVapdu.Length = m.ManagedObjectId.Size() + m.PollMdibDataReply.Size()

	// Длина ROapdus: ROIVapdu (без своего поля Length) + ее поле Length
	m.ROapdus.Length = (m.ROIVapdu.Size() - 2) + m.ROIVapdu.Length // Subtract ROIVapdu.Length size, add its value

	// Маршалинг каждой части в правильном порядке.
	spduData, err := m.SPpdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга SPpdu: %w", err)
	}
	buf.Write(spduData)

	roapdusData, err := m.ROapdus.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ROapdus: %w", err)
	}
	buf.Write(roapdusData)

	roivapduData, err := m.ROIVapdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ROIVapdu: %w", err)
	}
	buf.Write(roivapduData)

	managedObjectData, err := m.ManagedObjectId.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ManagedObjectId: %w", err)
	}
	buf.Write(managedObjectData)

	pollMdibDataReplyData, err := m.PollMdibDataReply.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга PollMdibDataReply: %w", err)
	}
	buf.Write(pollMdibDataReplyData)

	return buf.Bytes(), nil
}
