package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/client"
)

func main() {
	monitorAddr := flag.String("monitor", "", "Адрес монитора (IP:порт)")
	receiverAddr := flag.String("receiver", "", "Адрес получателя (IP:порт)")

	collectNumeric := flag.Bool("numeric", false, "Включить сбор числовых данных")
	numericPeriod := flag.Duration("numeric-period", 1500*time.Millisecond, "Период сбора числовых данных")

	collectAlarm := flag.Bool("alarm", false, "Включить сбор тревог")
	alarmPeriod := flag.Duration("alarm-period", 1500*time.Millisecond, "Период сбора тревог")

	flag.Parse()

	if *monitorAddr == "" {
		log.Fatal("Необходимо указать адрес монитора (--monitor)")
	}
	if *receiverAddr == "" {
		log.Fatal("Необходимо указать адрес получателя (--receiver)")
	}

	if !*collectNumeric && !*collectAlarm {
		log.Fatal("Необходимо включить хотя бы один тип сбора данных (--numeric или --alarm)")
	}

	client := client.NewComputerClient(*monitorAddr, *receiverAddr)

	ctxTest, cancelTest := context.WithCancel(context.Background())
	defer cancelTest()

	if err := client.Connect(ctxTest); err != nil {
		log.Fatalf("error when establishing connection: %v", err)
	}

	if *collectNumeric {
		log.Printf("Запуск сбора числовых данных с периодом %v", *numericPeriod)
		go client.CollectNumerics(*numericPeriod)
	}

	if *collectAlarm {
		log.Printf("Запуск сбора тревог с периодом %v", *alarmPeriod)
		go client.CollectAlarms(*alarmPeriod)
	}

	<-ctxTest.Done()
}
