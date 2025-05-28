package intellivue

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
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
}

// NewComputerClient creates a new instance of ComputerClient
func NewComputerClient(host, port string) *ComputerClient {
	log.Printf("Initializing ComputerClient with host=%s, port=%s", host, port)
	return &ComputerClient{
		host:    host,
		port:    port,
		streams: make(map[string]*DataStream),
	}
}

// Connect устанавливает UDP соединение с сервером
func (c *ComputerClient) Connect() error {
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
	return nil
}

// Close закрывает соединение
func (c *ComputerClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		log.Println("Соединение закрыто")
		return err
	}
	return nil
}

func (c *ComputerClient) SendAssociationRequest() error {
    if err := c.ensureConnection(); err != nil {
        return err
    }
    defer c.Close()

    msg := NewAssocReqMessage()
    data, err := msg.MarshalBinary()
    if err != nil {
        return fmt.Errorf("ошибка при создании сообщения: %v", err)
    }

    log.Printf("\n=== Отправляемое сообщение ===\n")
    printHexDump("AssociationRequest", data)

    if err := c.sendData(data); err != nil {
        return err
    }

    return c.listenForResponses()
}

// Вспомогательные методы
func (c *ComputerClient) ensureConnection() error {
    if c.conn == nil {
        return c.Connect()
    }
    return nil
}

func (c *ComputerClient) sendData(data []byte) error {
    if err := c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
        return fmt.Errorf("ошибка установки таймаута записи: %v", err)
    }

    n, err := c.conn.Write(data)
    if err != nil {
        return fmt.Errorf("ошибка при отправке: %v", err)
    }
    log.Printf("Отправлено %d байт по UDP\n", n)
    return nil
}

func (c *ComputerClient) listenForResponses() error {
    if err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
        return fmt.Errorf("ошибка установки таймаута чтения: %v", err)
    }

    udpConn, ok := c.conn.(*net.UDPConn)
    if !ok {
        return fmt.Errorf("соединение не является UDP соединением")
    }

    errChan := make(chan error, 1)
    go c.handleIncomingPackets(udpConn, errChan)
    return <-errChan
}

// Обработчики конкретных типов сообщений
func (c *ComputerClient) handleAssociationResponse(data []byte) error {
    log.Println("Обработка Association Response")
    return nil
}

func (c *ComputerClient) handleMDSCreateRequest(data []byte) error {
    log.Println("Обработка MDS Create Request")

    result := NewMDSCreateResult()
    resultBytes, err := result.MarshalBinary()
    if err != nil {
        return fmt.Errorf("ошибка маршалинга MDSCreateResult: %v", err)
    }

    if err := c.sendData(resultBytes); err != nil {
        return fmt.Errorf("ошибка отправки данных: %v", err)
    }
	printHexDump("Send MDS Create Result", resultBytes)

    return nil
}


func (c *ComputerClient) handleAssociationAbort(data []byte) error {
    log.Println("Обработка Association Abort")
    return fmt.Errorf("ассоциация прервана сервером")
}

// Основной обработчик входящих пакетов
func (c *ComputerClient) handleIncomingPackets(conn *net.UDPConn, errChan chan<- error) {
    defer close(errChan)
    
    for {
        buffer := make([]byte, 4096)
        n, _, err := conn.ReadFromUDP(buffer)
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                errChan <- nil
            } else {
                errChan <- fmt.Errorf("ошибка при чтении UDP: %v", err)
            }
            return
        }

        if n == 0 {
            continue
        }

        data := buffer[:n]
        firstByte := data[0]
        printHexDump(fmt.Sprintf("Получен пакет (0x%02X)", firstByte), data)

        switch firstByte {
        case 0x0E:
            if err := c.handleAssociationResponse(data); err != nil {
                errChan <- err
                return
            }
        case 0xE1:
            if err := c.handleMDSCreateRequest(data); err != nil {
                errChan <- err
                return
            }
        case 0x19:
            if err := c.handleAssociationAbort(data); err != nil {
                errChan <- err
                return
            }
        default:
            log.Printf("Неизвестный тип сообщения: 0x%02X\n", firstByte)
        }
    }
}

func (c *ComputerClient) Run(ctx context.Context) error {
	// Сначала устанавливаем ассоциацию
	if err := c.SendAssociationRequest(); err != nil {
		return fmt.Errorf("ошибка при установке ассоциации: %v", err)
	}

	log.Printf("Starting ComputerClient data collection routines (host=%s, port=%s)", c.host, c.port)

	// Start three different data collection goroutines
	go c.collectVitalSigns(ctx)
	go c.collectWaveforms(ctx)
	go c.collectAlarms(ctx)

	log.Println("All data collection routines started successfully")
	return nil
}

// GetLatestData returns the latest data for a specific stream
func (c *ComputerClient) GetLatestData(streamID string) (*DataStream, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, exists := c.streams[streamID]
	if !exists {
		log.Printf("Warning: Requested stream '%s' not found", streamID)
	}
	return data, exists
}

// UpdateStream updates the value of a specific stream
func (c *ComputerClient) UpdateStream(streamID string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.streams[streamID] = &DataStream{
		ID:    streamID,
		Value: value,
	}
	log.Printf("Updated stream '%s' with new data", streamID)
}

// Internal methods for collecting different types of data
func (c *ComputerClient) collectVitalSigns(ctx context.Context) {
	log.Println("Starting vital signs collection routine")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping vital signs collection routine")
			return
		case <-ticker.C:
			// Simulate collecting vital signs
			c.UpdateStream("vitals", map[string]interface{}{
				"heartRate": 75,
				"spO2":     98,
			})
		}
	}
}

func (c *ComputerClient) collectWaveforms(ctx context.Context) {
	log.Println("Starting waveforms collection routine")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping waveforms collection routine")
			return
		case <-ticker.C:
			// Simulate collecting waveform data
			c.UpdateStream("waveforms", map[string]interface{}{
				"ecg": []float64{1.0, 1.2, 0.8, 1.1},
			})
		}
	}
}

func (c *ComputerClient) collectAlarms(ctx context.Context) {
	log.Println("Starting alarms collection routine")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping alarms collection routine")
			return
		case <-ticker.C:
			// Simulate collecting alarm states
			c.UpdateStream("alarms", map[string]interface{}{
				"technical": []string{},
				"clinical":  []string{},
			})
		}
	}
}