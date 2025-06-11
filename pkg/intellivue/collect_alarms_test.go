package intellivue

import (
	"context"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestCollectAlarms(t *testing.T) {
	client := client.NewComputerClient("192.168.247.101:24105", "192.168.247.100:80")

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer connectCancel()

	if err := client.Connect(connectCtx); err != nil {
		t.Fatalf("error when establishing connection: %v", err)
	}

	client.CollectAlarms(1500 * time.Millisecond)

	select {}
}
