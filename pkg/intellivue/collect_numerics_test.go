package intellivue

import (
	"context"
	"testing"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func TestCollectNumerics(t *testing.T) {
	client := client.NewComputerClient("192.168.247.101", "24105")

	ctxTest, cancelTest := context.WithTimeout(context.Background(), 10*time.Second)

	if err := client.Connect(ctxTest); err != nil {
		cancelTest()
		t.Fatalf("error when establishing connectionы: %v", err)
	}

	doneCollecting := make(chan struct{})
	defer close(doneCollecting)
	client.CollectNumerics(ctxTest)

	<-ctxTest.Done()

	<-doneCollecting

	time.Sleep(100 * time.Millisecond)
	cancelTest()
	client.Close()
}
