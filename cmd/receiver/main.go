package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Определяем флаги командной строки
	port := flag.Int("port", 80, "Порт для прослушивания TCP соединений (по умолчанию: 80)")
	help := flag.Bool("help", false, "Показать справку")
	flag.Parse()

	// Показываем справку если запрошено
	if *help {
		fmt.Println("Использование: receiver [опции]")
		fmt.Println("\nОпции:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Проверяем корректность порта
	if *port < 1 || *port > 65535 {
		log.Fatalf("Некорректный номер порта: %d. Порт должен быть в диапазоне 1-65535", *port)
	}

	// Создаем TCP сервер
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Ошибка при создании TCP сервера: %v", err)
	}
	defer listener.Close()

	log.Printf("Ресивер запущен на порту %d", *port)
	log.Printf("Для завершения работы нажмите Ctrl+C")

	// Канал для обработки сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для завершения работы
	done := make(chan struct{})

	// Горутина для обработки сигналов завершения
	go func() {
		<-sigChan
		log.Println("Получен сигнал завершения")
		close(done)
	}()

	// Горутина для принятия соединений
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-done:
					return
				default:
					log.Printf("Ошибка при принятии соединения: %v", err)
					continue
				}
			}

			go handleConnection(conn)
		}
	}()

	// Ожидаем сигнал завершения
	<-done
	log.Println("Завершение работы ресивера")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Новое соединение от %s", clientAddr)

	reader := bufio.NewReader(conn)
	for {
		// Читаем данные до символа новой строки
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Ошибка при чтении данных от %s: %v", clientAddr, err)
			return
		}

		// Пропускаем пустые строки
		if len(data) <= 1 {
			continue
		}

		// Форматируем JSON для читаемого вывода
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, data[:len(data)-1], "", "    "); err != nil {
			log.Printf("Ошибка при форматировании JSON от %s: %v", clientAddr, err)
			continue
		}

		// Выводим отформатированный JSON
		fmt.Printf("\nПолучены данные от %s:\n%s\n", clientAddr, prettyJSON.String())
	}
}
