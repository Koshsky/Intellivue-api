package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
)

// DataStream represents a single data stream from the computer
type DataStream struct {
	ID    string
	Value interface{}
}

// ComputerClient handles communication with the medical device
type ComputerClient struct {
	host string
	port string
	mu   sync.RWMutex
	// streams stores the latest value for each stream
	streams map[string]*DataStream
	conn    net.Conn

	// Канал для сигнализации о завершении работы рутины runPacketListener и ее результате (включая завершение ассоциации)
	listenerOutcome chan error
}

// NewComputerClient creates a new instance of ComputerClient
func NewComputerClient(host, port string) *ComputerClient {
	log.Printf("Initializing ComputerClient with host=%s, port=%s", host, port)
	return &ComputerClient{
		host:    host,
		port:    port,
		streams: make(map[string]*DataStream),
		// Буферизация 1 нужна, чтобы не блокироваться при отправке nil после успешной ассоциации
		// до того, как Connect начнет читать из канала.
		listenerOutcome: make(chan error, 1),
	}
}

// EstablishUDPConnection устанавливает только UDP соединение с сервером и запускает обработчик входящих пакетов.
// Принимает контекст для управления жизненным циклом обработчика входящих пакетов.
func (c *ComputerClient) EstablishUDPConnection(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	log.Printf("Создание UDP соединения с %s...", addr)

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("ошибка разрешения UDP адреса: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("ошибка создания UDP соединения: %v", err)
	}

	c.conn = conn
	log.Println("UDP соединение установлено")

	// Начинаем слушать ответы после установки соединения, используя переданный контекст
	// runPacketListener будет управлять всем циклом, включая рукопожатие и дальнейшие пакеты.
	go c.runPacketListener(ctx, c.conn.(*net.UDPConn), c.listenerOutcome) // Передаем listenerOutcome для сигнализации

	return nil
}

// Connect устанавливает UDP соединение с сервером и выполняет запрос ассоциации.
// Connect возвращает nil при успешной ассоциации или ошибку при неудаче (включая ошибки в runPacketListener
// во время рукопожатия). После успешного Connect, runPacketListener продолжает работать
// и сигнализирует о дальнейших ошибках/закрытии через канал listenerOutcome.
func (c *ComputerClient) Connect(ctx context.Context) error {
	// Устанавливаем базовое UDP соединение и запускаем слушателя (runPacketListener)
	if err := c.EstablishUDPConnection(ctx); err != nil {
		return fmt.Errorf("ошибка при установке UDP соединения: %w", err)
	}

	// Установка ассоциации: создание и отправка AssociationRequest
	msg := NewAssocReqMessage()
	data, err := msg.MarshalBinary()
	if err != nil {
		c.Close() // Закрываем соединение при ошибке создания сообщения
		return fmt.Errorf("ошибка при создании сообщения AssociationRequest: %v", err)
	}

	// Send the Association Request
	if err := c.sendData(data); err != nil {
		c.Close() // Закрываем соединение при ошибке отправки
		return fmt.Errorf("ошибка отправки запроса ассоциации: %w", err)
	}

	log.Println("AssociationRequest отправлен. Ожидание ответа...")

	// Ожидаем завершения рукопожатия ассоциации от runPacketListener
	select {
	case <-ctx.Done():
		// Если контекст отменен во время ожидания ответа ассоциации
		c.Close()
		// listenerOutcome будет закрыт рутиной слушателя при отмене контекста
		return ctx.Err() // Возвращаем ошибку контекста
	case err := <-c.listenerOutcome:
		// Получен сигнал о завершении рукопожатия или ошибке от runPacketListener
		if err != nil {
			return fmt.Errorf("рукопожатие ассоциации завершилось с ошибкой: %w", err) // Возвращаем полученную ошибку
		}
		log.Println("Рукопожатие ассоциации успешно завершено.")

		// runPacketListener продолжает работать и обрабатывать последующие пакеты.

		return nil // Connect успешно завершен
	}
}

// Close закрывает соединение
func (c *ComputerClient) Close() error {
	c.mu.Lock() // Защита от множественного закрытия
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		log.Println("Соединение закрыто")
		// TODO: Возможно, нужно закрыть канал listenerOutcome здесь, если рутина слушателя завершилась не сама.
		// Но лучше, чтобы канал закрывала сама рутина runPacketListener при своем завершении.
		return err
	}
	return nil
}

// sendData отправляет данные по UDP
func (c *ComputerClient) sendData(data []byte) error {
	c.mu.RLock() // Используем RLock, так как только читаем c.conn
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("соединение не установлено")
	}

	// Устанавливаем таймаут записи, чтобы отправка не блокировалась вечно
	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return fmt.Errorf("ошибка установки таймаута записи: %v", err)
	}

	// Сбрасываем таймаут записи после операции, чтобы не мешать последующим операциям (чтению)
	defer conn.SetWriteDeadline(time.Time{})

	n, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("ошибка при отправке: %v", err)
	}
	log.Printf("Отправлено %d байт по UDP\n", n)
	return nil
}

// GetLatestData возвращает последние данные для указанного потока
func (c *ComputerClient) GetLatestData(streamID string) (*DataStream, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, exists := c.streams[streamID]
	if !exists {
		log.Printf("Warning: Requested stream '%s' not found", streamID)
	}
	return data, exists
}

// UpdateStream обновляет значение указанного потока данных
// Эта функция будет вызываться из обработчиков пакетов данных (например, MDS Poll Result)
func (c *ComputerClient) UpdateStream(streamID string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.streams[streamID] = &DataStream{
		ID:    streamID,
		Value: value,
	}
	log.Printf("Updated stream '%s' with new data", streamID)
}
