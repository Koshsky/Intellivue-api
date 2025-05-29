package packages

import (
	"bytes"
	"fmt"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

// MDSCreateResult представляет структуру данных для результата операции создания,
// состоящую из нескольких логических частей.
type MDSCreateResult struct {
	SPpdu
	ROapdus
	ROIVapdu
	ManagedObjectId // Используем существующую структуру
	CreateResultPayload
}

// Size возвращает общую длину MDSCreateResult в байтах.
func (m *MDSCreateResult) Size() uint16 {
	return m.SPpdu.Size() + m.ROapdus.Size() + m.ROIVapdu.Size() + m.ManagedObjectId.Size() + m.CreateResultPayload.Size()
}

// MarshalBinary кодирует структуру MDSCreateResult в бинарный формат.
func (m *MDSCreateResult) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Обновляем поля длины во вложенных структурах перед маршалингом.
	// Длина CreateResultPayload фиксирована (EventTime, EventType, Length).
	m.CreateResultPayload.Length = m.CreateResultPayload.Size() - 2 // Exclude its own Length field

	// Длина ROIVapdu: ManagedObjectId + CreateResultPayload
	m.ROIVapdu.Length = m.ManagedObjectId.Size() + m.CreateResultPayload.Size()

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

	createResultPayloadData, err := m.CreateResultPayload.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга CreateResultPayload: %w", err)
	}
	buf.Write(createResultPayloadData)

	return buf.Bytes(), nil
}

// NewMDSCreateResult создает новый экземпляр MDSCreateResult с предопределенными значениями.
func NewMDSCreateResult() *MDSCreateResult {
	sp := SPpdu{
		SessionID:  0xE100,
		PContextID: 2,
	}

	roap := ROapdus{
		ROType: RORS_APDU,
		// Length будет рассчитана перед маршалингом
	}

	roiv := ROIVapdu{
		InvokeID:    1,
		CommandType: CMD_CONFIRMED_EVENT_REPORT,
		// Length будет рассчитана перед маршалингом
	}

	managedObj := ManagedObjectId{
		MObjClass: NOM_MOC_VMS_MDS,
		ContextId: 0,
		Handle:    0,
	}

	payload := CreateResultPayload{
		EventTime: 0x00afbd00, // TODO: определить число
		EventType: NOM_NOTI_MDS_CREAT,
		// Length будет рассчитана перед маршалингом
	}

	return &MDSCreateResult{
		SPpdu:               sp,
		ROapdus:             roap,
		ROIVapdu:            roiv,
		ManagedObjectId:     managedObj,
		CreateResultPayload: payload,
	}
}
