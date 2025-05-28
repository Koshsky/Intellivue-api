package intellivue

import (
	"testing"
)

func TestSendAssociationRequest(t *testing.T) {
	// Создаем клиента с тестовыми параметрами
	client := NewComputerClient("192.168.247.101", "24105")

	// Отправляем запрос на ассоциацию по UDP
	err := client.SendAssociationRequest()
	if err != nil {
		t.Fatalf("Ошибка при отправке UDP запроса на ассоциацию: %v", err)
	}
} 