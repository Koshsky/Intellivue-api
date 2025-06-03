package client

import (
	"fmt"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
)

func (c *ComputerClient) SendPollNumericAction(invokeID uint16) error {
	c.SafeLog(fmt.Sprintf("Отправка SinglePollDataRequest запроса с InvokeID: %d", invokeID))

	msg := packages.NewSinglePollDataRequest(invokeID, NOM_MOC_VMO_METRIC_NU)

	dataToSend, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("ошибка маршалинга SinglePollDataRequest: %w", err)
	}

	if err := c.sendData(dataToSend); err != nil {
		return fmt.Errorf("ошибка отправки SinglePollDataRequest: %w", err)
	}

	c.SafeLog("SinglePollDataRequest пакет отправлен.")

	return nil
}
