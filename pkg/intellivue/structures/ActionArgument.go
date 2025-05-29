package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ActionArgument представляет структуру ActionArgument из MDSPollAction.
type ActionArgument struct {
	ActionArgumentManagedObject ManagedObjectId
	Scope                       uint32
	ActionType                  uint16
	Length                      uint16 // Длина того, что следует за Length
}

// Size возвращает длину ActionArgument в байтах.
// Длина включает ActionArgumentManagedObject, Scope, ActionType и Length.
// Предполагается, что ManagedObjectId имеет фиксированную длину (например, 4 байта).
func (a *ActionArgument) Size() uint16 {
	// Предполагая, что ManagedObjectId имеет размер 4 байта (uint32)
	return 4 + 4 + 2 + 2 // ManagedObjectId (4) + Scope (4) + ActionType (2) + Length (2)
}

// MarshalBinary кодирует структуру ActionArgument в бинарный формат.
func (a *ActionArgument) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Маршалинг ManagedObjectId. Требуется его MarshalBinary метод.
	objData, err := a.ActionArgumentManagedObject.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга ActionArgumentManagedObject: %w", err)
	}
	buf.Write(objData)

	if err := binary.Write(buf, binary.BigEndian, a.Scope); err != nil {
		return nil, fmt.Errorf("ошибка записи Scope: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, a.ActionType); err != nil {
		return nil, fmt.Errorf("ошибка записи ActionType: %w", err)
	}

	// Length будет установлена при маршалинге родительской структуры (MDSPollAction)
	// или структуры, которая ее включает (PollMdibDataReq), чтобы отразить длину последующих данных.
	// Здесь мы просто записываем текущее значение.
	if err := binary.Write(buf, binary.BigEndian, a.Length); err != nil {
		return nil, fmt.Errorf("ошибка записи Length: %w", err)
	}

	return buf.Bytes(), nil
}
