package client

import (
	"encoding/binary"
	"fmt"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

// runDataExportHandler прослушивает канал dataExportChan и обрабатывает пакеты Data Export Protocol.
func (c *ComputerClient) runDataExportHandler() {
	defer c.wg.Done() // Уменьшаем счетчик WaitGroup при завершении горутины
	c.SafeLog("Запуск Data Export Handler рутины...")

	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("Data Export Handler: Контекст отменен, завершение рутины.")
			return
		case packetData, ok := <-c.dataExportChan:
			if !ok {
				c.SafeLog("Data Export Handler: Канал dataExportChan закрыт, завершение рутины.")
				return
			}
			if len(packetData) > 0 {
				c.SafeLog("Data Export Handler: Пакет E1 (0x%02X) извлечен из dataExportChan.", packetData[0])
			}

			c.handleDataExportPacket(packetData)
		}
	}
}

func (c *ComputerClient) handleROIVAPDU(data []byte) error {
	if len(data) < 8 {
		return fmt.Errorf("недостаточно данных для определения command_type в пакете ROIV_APDU")
	}

	commandType := structures.CMDType(binary.BigEndian.Uint16(data[6:8]))

	c.SafeLog(fmt.Sprintf("Получен пакет ROIV_APDU, command_type: 0x%04X", commandType))

	switch commandType {
	case structures.CMD_CONFIRMED_EVENT_REPORT:
		c.SafeLog("\n\nОбнаружен CommandType: CMD_CONFIRMED_EVENT_REPORT")
		// TODO: Парсинг и обработка CMD_CONFIRMED_EVENT_REPORT

	case structures.CMD_EVENT_REPORT:
		c.SafeLog("\n\nОбнаружен CommandType: CMD_EVENT_REPORT")
		// TODO: Парсинг и обработка CMD_EVENT_REPORT

	case structures.CMD_CONFIRMED_ACTION:
		c.SafeLog("\n\nОбнаружен CommandType: CMD_CONFIRMED_ACTION")
		// TODO: Парсинг и обработка CMD_CONFIRMED_ACTION

	case structures.CMD_GET:
		c.SafeLog("\n\nОбнаружен CommandType: CMD_GET")
		// TODO: парсинг и обработка CMD_GET

	case structures.CMD_CONFIRMED_SET:
		c.SafeLog("\n\nОбнаружен CommandType: CMD_CONFIRMED_SET")
		// TODO: Парсинг и обработка CMD_CONFIRMED_SET

	default:
		c.SafeLog(fmt.Sprintf("  Неизвестный command_type в ROIV_APDU: 0x%04X. Игнорируем.", commandType))
	}

	utils.PrintHexDump(&c.printMu, "ROIV_APDU", data)

	return nil
}
