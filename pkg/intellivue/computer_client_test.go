package intellivue

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestRunClientAndPoll(t *testing.T) {
	// Создаем клиента с тестовыми параметрами
	client := client.NewComputerClient("192.168.247.101", "24105")
	defer client.Close() // Убедимся, что соединение будет закрыто

	// Создаем контекст с таймаутом для всего теста
	ctxTest, cancelTest := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelTest() // Гарантируем отмену контекста

	// Устанавливаем соединение и запускаем обработчик входящих пакетов
	log.Println("Установка соединения...")
	// Передаем контекст теста в Connect
	if err := client.Connect(ctxTest); err != nil {
		t.Fatalf("Ошибка при установке соединения: %v", err)
	}
	log.Println("Соединение установлено.")

	log.Println("Тест RunClientAndPoll завершен.")

	// Небольшая задержка для уверенности, что все горутины успели завершить логирование
	time.Sleep(100 * time.Millisecond)
}
