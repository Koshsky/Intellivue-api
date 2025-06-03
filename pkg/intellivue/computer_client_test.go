package intellivue

import (
	"context"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestRunClientAndPoll(t *testing.T) {
	client := client.NewComputerClient("192.168.247.101", "24105")
	// defer client.Close() // Закрытие клиента будет выполнено после ожидания горутины

	// Увеличиваем таймаут до 10 секунд для тестового опроса
	ctxTest, cancelTest := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancelTest() // Отмена контекста будет выполнена после ожидания горутины

	client.SafeLog("Установка соединения...")
	if err := client.Connect(ctxTest); err != nil {
		cancelTest() // Отменяем контекст при ошибке соединения
		t.Fatalf("Ошибка при установке соединения: %v", err)
	}
	client.SafeLog("Соединение установлено.")

	// Канал для сигнализации о завершении горутины CollectNumerics
	doneCollecting := make(chan struct{})

	go func() {
		defer close(doneCollecting) // Сигнализируем о завершении горутины
		client.CollectNumerics(ctxTest)
	}()

	client.SafeLog("Ожидание завершения теста (10 секунд)...")
	<-ctxTest.Done() // Ждем отмены контекста теста (через 10 секунд или при ошибке)
	client.SafeLog("Контекст теста отменен. Ожидание завершения CollectNumerics...")

	<-doneCollecting // Ждем завершения горутины CollectNumerics
	client.SafeLog("CollectNumerics завершена.")

	time.Sleep(100 * time.Millisecond) // Небольшая задержка для чистоты вывода
	client.SafeLog("Закрытие клиента...")
	cancelTest()   // Отменяем контекст теста
	client.Close() // Теперь безопасно закрыть клиент и его логгер
	client.SafeLog("Клиент закрыт.")
}
