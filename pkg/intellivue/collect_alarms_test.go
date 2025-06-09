package intellivue

import (
	"context"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestCollectAlarms(t *testing.T) {
	client := client.NewComputerClient("192.168.247.101:24105", "192.168.247.100:80")

	ctxTest, cancelTest := context.WithTimeout(context.Background(), 20*time.Second)

	if err := client.Connect(ctxTest); err != nil {
		cancelTest()
		t.Fatalf("error when establishing connectionы: %v", err)
	}

	client.CollectAlarms()

	<-ctxTest.Done()

	time.Sleep(100 * time.Millisecond)
	cancelTest()
}
