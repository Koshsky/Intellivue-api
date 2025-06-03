package client

import (
	"context"
	"time"
)

// CollectNumerics собирает и обрабатывает пакеты жизненно важных показателей.
// Реализована отправка SinglePollDataRequest запроса для сбора числовых данных.
func (c *ComputerClient) CollectNumerics(ctx context.Context) {
	var invokeID uint16 = 1

	pollInterval := 15 * time.Second // Изменил интервал на 2 секунды, как было запрошено в одном из предыдущих сообщений
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

// CollectWaveforms собирает и обрабатывает пакеты данных осциллограмм.
func (c *ComputerClient) CollectWaveforms(ctx context.Context) {
	c.SafeLog("Запуск рутины CollectWaveforms...")
	// Пример цикла, который будет остановлен при отмене контекста
	for {
		select {
		case <-ctx.Done():
			c.SafeLog("СollectWaveforms: Контекст отменен, завершение рутины.")
			return
		default:
			// Здесь будет логика чтения данных Waveforms
			// time.Sleep для предотвращения бесконечного цикла на данный момент
			time.Sleep(5 * time.Second)
		}
	}
}

// CollectAlarms собирает и обрабатывает пакеты тревог.
func (c *ComputerClient) CollectAlarms(ctx context.Context) {
	c.SafeLog("Запуск рутины CollectAlarms...")
	// Пример цикла, который будет остановлен при отмене контекста
	for {
		select {
		case <-ctx.Done():
			c.SafeLog("CollectAlarms: Контекст отменен, завершение рутины.")
			return
		default:
			// Здесь будет логика чтения данных Alarms
			// Например, отправка запросов Polled Event Report
			// time.Sleep для предотвращения бесконечного цикла на данный момент
			time.Sleep(5 * time.Second)
		}
	}
}

// sendPollActionAlarm отправляет запрос Polled Event Report для получения информации о тревогах.
func (c *ComputerClient) sendPollActionAlarm() {
	c.SafeLog("Отправка запроса Polled Event Report...")
	// Здесь будет логика создания и отправки пакета Polled Event Report
}
