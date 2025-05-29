package packages

import (
	"bytes"
	"fmt"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

// MDSPollAction представляет структуру сообщения MDS Poll Action Request.
// Она включает стандартные заголовки и специфическую для действия полезную нагрузку.
type MDSPollAction struct {
	SPpdu           SPpdu
	ROapdus         ROapdus
	ROIVapdu        ROIVapdu
	ActionArgument  ActionArgument  // Полезная нагрузка действия
	PollMdibDataReq PollMdibDataReq // Специфические данные для запроса Poll MDIB Data
}

// NewMDSPollAction создает новый экземпляр MDSPollAction.
// pcontext_id: Presentation Context ID
// invoke_id: Invoke ID для запроса
// code: Код для поля Code в PollMdibDataReq (например, NOM_MOC_VMO_AL_MON)
func NewMDSPollAction(pcontext_id, invoke_id, code uint16) *MDSPollAction {
	// Инициализируем вложенные структуры с базовыми значениями
	sp := SPpdu{
		SessionID:  0xE100, // Пример значения, возможно, потребуется динамическое определение
		PContextID: pcontext_id,
	}

	roap := ROapdus{
		ROType: ROIV_APDU, // Тип APDU для запроса (Invoke)
		// Length будет рассчитана перед маршалингом
	}

	roiv := ROIVapdu{
		InvokeID:    invoke_id,
		CommandType: CMD_CONFIRMED_ACTION, // Тип команды (Confirmed Action)
		// Length будет рассчитана перед маршалингом
	}

	pollReq := PollMdibDataReq{
		PollNumber:    1, // Пример значения, возможно, потребуется инкремент
		Partition:     0, // Пример значения
		Code:          code,
		PolledAttrGrp: 0, // Пример значения
	}

	// ActionArgument содержит PollMdibDataReq как часть полезной нагрузки.
	// ManagedObjectId в ActionArgument может быть нулем или указывать на конкретный объект.
	// Scope и ActionType также должны быть установлены.
	actionArg := ActionArgument{
		ActionArgumentManagedObject: ManagedObjectId{MObjClass: 0, ContextId: 0, Handle: 0}, // TODO: Определить ManagedObjectId, если требуется конкретный объект
		Scope:                       0,                                                      // TODO: Определить Scope
		ActionType:                  NOM_ACT_POLL_MDIB_DATA,                                 // Тип действия
		// Length будет рассчитана перед маршалингом
	}

	return &MDSPollAction{
		SPpdu:           sp,
		ROapdus:         roap,
		ROIVapdu:        roiv,
		ActionArgument:  actionArg,
		PollMdibDataReq: pollReq,
	}
}

// MarshalBinary кодирует структуру MDSPollAction в бинарный формат.
func (m *MDSPollAction) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Маршалинг PollMdibDataReq
	pollReqData, err := m.PollMdibDataReq.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга PollMdibDataReq: %w", err)
	}

	// Обновляем поле Length в ActionArgument. Оно включает только PollMdibDataReq.
	m.ActionArgument.Length = uint16(len(pollReqData))

	// Маршалинг ActionArgument (теперь с обновленной длиной)
	actionArgData, err := m.ActionArgument.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ActionArgument: %w", err)
	}

	// Обновляем поле Length в ROIVapdu. Оно включает ActionArgument.
	m.ROIVapdu.Length = uint16(len(actionArgData))

	// Маршалинг ROIVapdu (теперь с обновленной длиной)
	roivData, err := m.ROIVapdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ROIVapdu: %w", err)
	}

	// Обновляем поле Length в ROapdus. Оно включает ROIVapdu.
	m.ROapdus.Length = uint16(len(roivData))

	// Маршалинг ROapdus (теперь с обновленной длиной)
	roapdusData, err := m.ROapdus.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ROapdus: %w", err)
	}

	// Маршалинг SPpdu
	spduData, err := m.SPpdu.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга SPpdu: %w", err)
	}

	// Записываем все части в буфер в правильном порядке
	buf.Write(spduData)
	buf.Write(roapdusData)
	buf.Write(roivData)
	buf.Write(actionArgData)
	buf.Write(pollReqData)

	return buf.Bytes(), nil
}

// Size возвращает общую длину MDSPollAction в байтах.
func (m *MDSPollAction) Size() uint16 {
	// Суммируем размеры всех вложенных структур.
	// Размеры полей Length внутри ROapdus, ROIVapdu и ActionArgument
	// учитываются в их собственных методах Size().
	return m.SPpdu.Size() + m.ROapdus.Size() + m.ROIVapdu.Size() + m.ActionArgument.Size() + m.PollMdibDataReq.Size()
}
