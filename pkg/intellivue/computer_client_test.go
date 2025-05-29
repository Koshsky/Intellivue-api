package intellivue

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRunClientAndPoll(t *testing.T) {
	// Создаем клиента с тестовыми параметрами
	client := NewComputerClient("192.168.247.101", "24105")
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

	ctxCollect := ctxTest
	// Запускаем только рутину collectAlarms
	log.Println("Запуск рутины collectAlarms...")
	go client.collectAlarms(ctxCollect)

	// Ждем завершения рутины collectAlarms (она завершится по таймауту) ИЛИ отмены контекста теста
	select {
	case <-ctxCollect.Done():
		log.Println("Рутина collectAlarms завершена по таймауту или контекст отменен.")
	case <-ctxTest.Done():
		log.Println("Контекст теста отменен. Рутина collectAlarms также должна остановиться.")
	}

	log.Println("Тест RunClientAndPoll завершен.")

	// Небольшая задержка для уверенности, что все горутины успели завершить логирование
	time.Sleep(100 * time.Millisecond)
}

// TestEstablishUDPConnection проверяет установку только UDP соединения без отправки AssociationRequest.
func TestEstablishUDPConnection(t *testing.T) {
	// Создаем клиента с тестовыми параметрами
	client := NewComputerClient("192.168.247.101", "24105")
	defer client.Close() // Убедимся, что соединение будет закрыто

	// Создаем контекст с таймаутом для теста
	ctxTest, cancelTest := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTest() // Гарантируем отмену контекста

	log.Println("Установка только UDP соединения...")
	// Устанавливаем только UDP соединение и запускаем обработчик входящих пакетов
	if err := client.EstablishUDPConnection(ctxTest); err != nil {
		t.Fatalf("Ошибка при установке UDP соединения: %v", err)
	}
	log.Println("UDP соединение установлено. Ожидание пакетов в течение 5 секунд...")

	// Ожидаем отмены контекста. Входящие пакеты (если будут) будут логироваться в handleIncomingPackets.
	client.collectAlarms(ctxTest)
	<-ctxTest.Done()

	log.Println("Тест EstablishUDPConnection завершен.")

	// Небольшая задержка для уверенности, что все горутины успели завершить логирование
	time.Sleep(100 * time.Millisecond)
}
