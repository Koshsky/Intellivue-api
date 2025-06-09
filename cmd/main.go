package main

import (
	"context"
	"log"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func main() {
	client := client.NewComputerClient("192.168.247.101:24105", "192.168.247.100:80")

	ctxTest, cancelTest := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelTest()

	if err := client.Connect(ctxTest); err != nil {
		log.Fatalf("error when establishing connection: %v", err)
	}

	go client.CollectNumerics()
	// go client.CollectAlarms()

	<-ctxTest.Done()
}
