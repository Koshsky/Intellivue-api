package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	// . "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

// handleDataExportPacket обрабатывает пакеты типа Data Export Protocol (0xE1).
// Здесь должна быть дальнейшая диспетчеризация по подтипам (MDS Poll, MDS Result, etc.).
// ctx используется для возможности отмены обработки длительных пакетов.
func (c *ComputerClient) handleDataExportPacket(ctx context.Context, data []byte) error {
	log.Println("Обработка пакета Data Export Protocol (0xE1)...")

	// ro_type находится со смещением 4-6 байт (uint16)
	if len(data) < 6 {
		return fmt.Errorf("недостаточно данных для определения ro_type в пакете 0xE1")
	}

	roType := binary.BigEndian.Uint16(data[4:6])

	log.Printf("Пакет 0xE1, ro_type: 0x%04X", roType)

	switch roType {
	case constants.ROIV_APDU:
		return c.handleROIVAPDU(ctx, data)

	case constants.RORS_APDU:
		log.Println("Обнаружен пакет: Single Poll Data Result (RORS_APDU)")

	case constants.ROLRS_APDU:
		log.Println("Обнаружен пакет: Single Poll Data Result (Linked) (ROLRS_APDU)")

	case constants.ROER_APDU:
		log.Println("Обнаружен пакет: RO Error (ROER_APDU)")
		utils.PrintHexDump("ROER_APDU", data)

	default:
		log.Printf("Обнаружен пакет 0xE1 с неизвестным ro_type: 0x%04X. Игнорируем.", roType)
	}

	return nil
}

func (c *ComputerClient) handleROIVAPDU(ctx context.Context, data []byte) error {

	// command_type находится со смещением 10-12 байт (uint16)
	if len(data) < 12 {
		return fmt.Errorf("недостаточно данных для определения command_type в пакете ROIV_APDU")
	}

	commandType := binary.BigEndian.Uint16(data[10:12])

	log.Printf("ROIV APDU, command_type: 0x%04X", commandType)

	switch commandType {
	case constants.CMD_CONFIRMED_EVENT_REPORT:
		log.Println("Обнаружен пакет: MDS CREATE EVENT (CMD_CONFIRMED_EVENT_REPORT).")

		log.Println("Отправка MDS CREATE RESULT...")
		msg := packages.NewMDSCreateResult()
		dataToSend, err := msg.MarshalBinary()
		if err != nil {
			return fmt.Errorf("ошибка при создании сообщения MDS CREATE RESULT: %v", err)
		}
		if err := c.sendData(dataToSend); err != nil {
			return fmt.Errorf("ошибка отправки MDS CREATE RESULT: %v", err)
		}

	case constants.CMD_EVENT_REPORT:
		log.Println("Обнаружен пакет: CONNECT INDICATION EVENT (CMD_EVENT_REPORT).")
	case constants.CMD_CONFIRMED_ACTION:
		log.Println("Обнаружен пакет: SINGLE POLL DATA REQUEST / EXTENDED POLL DATA REQUEST (CMD_CONFIRMED_ACTION).")
	case constants.CMD_GET:
		log.Println("Обнаружен пакет: GET PRIORITY LIST REQUEST (CMD_GET).")
	case constants.CMD_CONFIRMED_SET:
		log.Println("Обнаружен пакет: CMD_CONFIRMED_SET.")
	default:
		log.Printf("Обнаружен пакет ROIV_APDU с неизвестным command_type: 0x%04X.", commandType)
	}

	return nil
}
