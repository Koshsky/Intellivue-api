package intellivue

import (
	"context"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestCollectAlarms(t *testing.T) {
	client := client.NewComputerClient("192.168.247.101:24105", "192.168.247.100:80")

	// Создаем контекст с таймаутом только для установки соединения
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer connectCancel()

	if err := client.Connect(connectCtx); err != nil {
		t.Fatalf("error when establishing connection: %v", err)
	}

	// После успешного соединения используем бесконечный контекст
	client.CollectAlarms()

	// Ждем бесконечно, так как контекст не должен отменяться
	select {}
}
