package client

import (
	"context"
	"time"
)

// CollectNumerics собирает и обрабатывает пакеты жизненно важных показателей.
// Реализована отправка SinglePollDataRequest запроса для сбора числовых данных.
func (c *ComputerClient) CollectNumerics(ctx context.Context) {
	var invokeID uint16 = 1

	pollInterval := 15 * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	if err := c.SendPollNumericAction(invokeID); err != nil {
		c.SafeLog("Failed to process SinglePollDataResultLinked: %v", err)
	}
	invokeID++

	for {
		select {
		case <-ctx.Done():
			c.SafeLog("CollectNumerics: Контекст отменен, завершение рутины.")
			return
		case <-ticker.C:
			c.SafeLog("Отправка очередного SinglePollDataRequest запроса по таймеру.")
			if err := c.SendPollNumericAction(invokeID); err != nil {
				c.SafeLog("Ошибка при отправке SinglePollDataRequest с InvokeID %d: %v", invokeID, err)
			}
			invokeID++
		}
	}
}

func (c *ComputerClient) CollectWaveforms(ctx context.Context) {
	c.SafeLog("Запуск рутины CollectWaveforms...")
	// Пример цикла, который будет остановлен при отмене контекста
	for {
		select {
		case <-ctx.Done():
			c.SafeLog("СollectWaveforms: Контекст отменен, завершение рутины.")
			return
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *ComputerClient) CollectAlarms(ctx context.Context) {
	var invokeID uint16 = 100

	pollInterval := 10 * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	if err := c.SendPollAlarmAction(invokeID); err != nil {
		c.SafeLog("Failed to process SinglePollDataResult: %v", err)
	}
	invokeID++

	for {
		select {
		case <-ctx.Done():
			c.SafeLog("CollectAlarms: Контекст отменен, завершение рутины.")
			return
		case <-ticker.C:
			c.SafeLog("Отправка очередного SinglePollDataRequest запроса по таймеру.")
			if err := c.SendPollAlarmAction(invokeID); err != nil {
				c.SafeLog("Ошибка при отправке SinglePollDataRequest с InvokeID %d: %v", invokeID, err)
			}
			invokeID++
		}
	}
}
