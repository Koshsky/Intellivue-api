package client

import (
	"fmt"
	"log"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
)

// SendPollNumericAction отправляет запрос MDS Poll Action для сбора числовых данных.
// invokeID: Invoke ID для запроса (должен быть уникальным для каждой отправки запроса).
func (c *ComputerClient) SendPollNumericAction(invokeID uint16) error {
	log.Printf("Отправка MDS Poll Action запроса с InvokeID: %d", invokeID)

	// Создаем сообщение MDSPollAction с InvokeID и кодом для числовых данных
	// Используем constants.NOM_MOC_VMO_METRIC_NU в качестве кода, как вы указали.
	msg := packages.NewMDSPollAction(invokeID, constants.NOM_MOC_VMO_METRIC_NU)

	// Маршалинг сообщения в байты
	dataToSend, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("ошибка маршалинга MDSPollAction: %w", err)
	}

	// Отправка данных
	if err := c.sendData(dataToSend); err != nil {
		return fmt.Errorf("ошибка отправки MDSPollAction: %w", err)
	}

	log.Println("MDSPollAction запрос отправлен.")

	return nil
}
