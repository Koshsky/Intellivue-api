package client

import (
	"context"
	"log"
	"time"
)

// CollectNumerics собирает и обрабатывает пакеты жизненно важных показателей.
// Реализована отправка MDS Poll Action запроса для сбора числовых данных.
func (c *ComputerClient) CollectNumerics(ctx context.Context) {
	log.Println("Запуск рутины CollectNumerics...")

	// Используем invokeID, начиная с 1. В реальном приложении, возможно, стоит сделать его потокобезопасным
	// или использовать другой механизм управления invokeID.
	var invokeID uint16 = 1

	// Интервал отправки запросов (например, каждые 5 секунд)
	pollInterval := 5 * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Отправляем первый запрос сразу при запуске
	if err := c.SendPollNumericAction(invokeID); err != nil {
		log.Printf("Ошибка при первой отправке MDS Poll Action: %v", err)
	}
	invokeID++

	// Цикл для периодической отправки запросов
	for {
		select {
		case <-ctx.Done():
			log.Println("CollectNumerics: Контекст отменен, завершение рутины.")
			return
		case <-ticker.C:
			// Время отправки очередного запроса
			log.Println("Отправка очередного MDS Poll Action запроса по таймеру.")
			if err := c.SendPollNumericAction(invokeID); err != nil {
				log.Printf("Ошибка при отправке MDS Poll Action с InvokeID %d: %v", invokeID, err)
			}
			invokeID++ // Инкрементируем invokeID для следующего запроса
		}
	}
}

// CollectWaveforms собирает и обрабатывает пакеты данных осциллограмм.
func (c *ComputerClient) CollectWaveforms(ctx context.Context) {
	log.Println("Запуск рутины CollectWaveforms...")
	// Пример цикла, который будет остановлен при отмене контекста
	for {
		select {
		case <-ctx.Done():
			log.Println("СollectWaveforms: Контекст отменен, завершение рутины.")
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
	log.Println("Запуск рутины CollectAlarms...")
	// Пример цикла, который будет остановлен при отмене контекста
	for {
		select {
		case <-ctx.Done():
			log.Println("CollectAlarms: Контекст отменен, завершение рутины.")
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
	log.Println("Отправка запроса Polled Event Report...")
	// Здесь будет логика создания и отправки пакета Polled Event Report
}
