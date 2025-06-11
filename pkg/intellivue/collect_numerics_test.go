package intellivue

import (
	"context"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestCollectNumerics(t *testing.T) {
	client := client.NewComputerClient("192.168.247.101:24105", "192.168.247.100:80")

	ctxTest, cancelTest := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTest()

	if err := client.Connect(ctxTest); err != nil {
		t.Fatalf("error when establishing connection: %v", err)
	}

	go client.CollectNumerics(1500 * time.Millisecond)

	<-ctxTest.Done()
}
