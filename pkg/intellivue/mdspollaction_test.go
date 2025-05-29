package intellivue

import (
	"testing"
	// Импортируем fmt и encoding/hex для PrintHexDump, если не используются в другом месте файла
	// Если PrintHexDump находится в отдельном пакете, нужно будет импортировать его

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

func TestMDSPollActionMarshal(t *testing.T) {
	// Создаем новый экземпляр MDSPollAction с предопределенными значениями
	mdsPollAction := NewMDSPollAction(2, 1, NOM_MOC_VMO_AL_MON)

	// Выполняем маршалинг структуры в бинарный формат
	data, err := mdsPollAction.MarshalBinary()
	if err != nil {
		t.Fatalf("Ошибка при маршалинге MDSPollAction: %v", err)
	}

	// Выводим бинарные данные в HEX-формате, используя общую функцию
	PrintHexDump("MDSPollAction HEX\n", data)

	// TODO: Добавить проверки на ожидаемые значения байт в будущем, когда спецификация будет более точной
}
