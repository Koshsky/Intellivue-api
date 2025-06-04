package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

// ComputerClient представляет клиента, взаимодействующего с устройством Intellivue.
type ComputerClient struct {
	// printMu мьютекс для синхронизации вывода в терминал из разных горутин.
	printMu sync.Mutex

	serverAddr string   // Адрес сервера (устройства Intellivue)
	serverPort string   // Порт сервера
	conn       net.Conn // Сетевое соединение
	// reader     *bufio.Reader // Буферизованный читатель для парсинга пакетов - Удалено для UDP

	// Канал для отправки пакетов на устройство.
	sendChan chan []byte
	// Канал для приема обработанных входящих пакетов (пока не используется).
	// receiveChan chan packages.DataPackage

	// Канал для получения пакетов Data Export Protocol для обработки.
	dataExportChan chan []byte

	// Контекст для управления жизненным циклом клиента и его горутин.
	ctx context.Context
	// Функция отмены контекста.
	cancel context.CancelFunc

	// WaitGroup для отслеживания запущенных горутин.
	wg sync.WaitGroup

	isAssociationDone bool // Флаг для отслеживания состояния ассоциации

	closeOnce sync.Once

	assocResponseChan chan struct{}
}

func NewComputerClient(addr, port string) *ComputerClient {
	ctx, cancel := context.WithCancel(context.Background())
	client := &ComputerClient{
		serverAddr:        addr,
		serverPort:        port,
		ctx:               ctx,
		cancel:            cancel,
		sendChan:          make(chan []byte, 10), // если используется
		assocResponseChan: make(chan struct{}),
	}
	client.wg.Add(1)
	go client.runPacketListener()
	return client
}

// Connect устанавливает UDP соединение с устройством Intellivue и выполняет процедуру ассоциации.
// Принимает контекст для управления жизненным циклом операции соединения.
// Возвращает nil при успешной ассоциации или ошибку при неудаче.
func (c *ComputerClient) Connect(ctx context.Context) error {
	c.SafeLog("Попытка установления UDP соединения и ассоциации с %s:%s", c.serverAddr, c.serverPort)

	addr := fmt.Sprintf("%s:%s", c.serverAddr, c.serverPort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("error creating UDP connection: %w", err)
	}
	c.conn = conn
	c.SafeLog("UDP connection established.")

	// Запускаем горутину для прослушивания входящих пакетов
	c.wg.Add(1) // Увеличиваем счетчик WaitGroup для runPacketListener
	go c.runPacketListener()

	// Запускаем горутину для отправки пакетов из sendChan
	c.wg.Add(1) // Увеличиваем счетчик WaitGroup для runPacketSender
	go c.runPacketSender()

	// Выполняем процедуру установления ассоциации
	// Создание сообщения AssociationRequest
	assocReq, err := packages.NewAssocReqMessage().MarshalBinary()
	if err != nil {
		c.Close() // Закрываем соединение при ошибке создания сообщения
		return fmt.Errorf("error creating AssociationRequest message: %w", err)
	}

	// Отправка запроса ассоциации
	if err := c.sendData(assocReq); err != nil {
		c.Close() // Закрываем соединение при ошибке отправки
		return fmt.Errorf("error sending AssociationRequest: %w", err)
	}

	c.SafeLog("AssociationRequest sent. Waiting for response...")

	select {
	case <-c.assocResponseChan:
		c.SafeLog("Association Response received.")
	case <-time.After(5 * time.Second):
		c.Close()
		return fmt.Errorf("timeout waiting for Association Response")
	}

	return nil
}

// Close gracefully закрывает соединение и останавливает все горутины клиента.
func (c *ComputerClient) Close() error {
	c.closeOnce.Do(func() {
		c.SafeLog("Starting client shutdown procedure...")

		c.cancel()
		c.wg.Wait()
		c.SafeLog("All goroutines completed.")

		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
			c.SafeLog("Network connection closed.")
		}

		// Закрываем только реально существующий канал
		if c.sendChan != nil {
			close(c.sendChan)
		}

		c.SafeLog("Client shutdown procedure completed.")
	})
	return nil
}

// runPacketSender отправляет пакеты из sendChan на устройство.
func (c *ComputerClient) runPacketSender() {
	defer c.wg.Done() // Уменьшаем счетчик WaitGroup при завершении горутины
	c.SafeLog("Starting Packet Sender goroutine...")

	for {
		select {
		case data, ok := <-c.sendChan:
			if !ok { // Канал закрыт
				c.SafeLog("Packet Sender: Sending channel closed, goroutine ending.")
				return
			}
			// TODO: Добавить мьютекс для записи в соединение, если multiple goroutines могут писать
			if err := c.sendData(data); err != nil {
				c.SafeLog("Error sending data: %v", err)
				// TODO: Обработка ошибок отправки (переподключение?)
			}
		case <-c.ctx.Done():
			c.SafeLog("Packet Sender: Context canceled, goroutine ending.")
			return
		}
	}
}

// sendData отправляет байты данных в сетевое соединение.
func (c *ComputerClient) sendData(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("connection not established")
	}
	// TODO: Установить таймаут записи для UDP
	_, err := c.conn.Write(data)
	return err
}

// runPacketListener прослушивает входящие пакеты от устройства по UDP и сразу обрабатывает их.
func (c *ComputerClient) runPacketListener() {
	defer c.wg.Done()
	c.SafeLog("Starting Packet Listener goroutine...")

	conn, ok := c.conn.(*net.UDPConn)
	if !ok {
		c.SafeLog("runPacketListener: Error: connection is not UDP.")
		return
	}

	buffer := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			c.SafeLog("runPacketListener: failed to unmarshal UDP: %v", err)
			return
		}
		if n == 0 {
			continue
		}
		data := make([]byte, n)
		copy(data, buffer[:n])
		c.SafeLog("runPacketListener: Received packet, first byte: 0x%02X", data[0])
		c.handleDataExportPacket(data)
	}
}

func (c *ComputerClient) SafeLog(format string, args ...interface{}) {
	c.printMu.Lock()
	defer c.printMu.Unlock()
	log.Printf(format, args...)
}

// handleDataExportPacket теперь сразу обрабатывает пакет
func (c *ComputerClient) handleDataExportPacket(data []byte) {
	firstByte := data[0]
	switch firstByte {
	case 0x0E:
		c.SafeLog("runPacketListener: Received Association Response (0x0E)")
		c.isAssociationDone = true
		select {
		case c.assocResponseChan <- struct{}{}:
		default:
		}
	case 0x19:
		c.SafeLog("runPacketListener: Received Association Abort (0x19).")
		c.Close()
		return
	case 0x0C:
		c.SafeLog("runPacketListener: Received Association Refuse (0x0C).")
		c.Close()
		return
	case 0xE1:
		c.SafeLog("runPacketListener: Received Data Export Protocol packet (0xE1). Processing...")
		if len(data) < 6 {
			c.SafeLog("Data Export Protocol: too short packet to determine ro_type")
			return
		}
		roType := binary.BigEndian.Uint16(data[4:6])
		switch roType {
		case base.ROIV_APDU: // ROIV_APDU
			c.SafeLog("Data Export Protocol: ROIV_APDU")
			// обработка ROIV_APDU
		case base.RORS_APDU: // RORS_APDU
			c.SafeLog("Data Export Protocol: RORS_APDU")
			result := &packages.SinglePollDataResult{}
			if err := result.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("failed to unmarshal SinglePollDataResult: %v", err)
				return
			}
			var mu sync.Mutex
			result.ShowInfo(&mu, 0)
			jsonBytes, err := json.MarshalIndent(result.PollMdibDataReply.PollInfoList, "", "  ")
			if err != nil {
				c.SafeLog("failed to marshal PollInfoList to JSON: %v", err)
				return
			}
			c.SafeLog("PollInfoList (JSON):\n%s", string(jsonBytes))
		case base.ROLRS_APDU: // ROLRS_APDU
			c.SafeLog("Data Export Protocol: ROLRS_APDU")
			linkedResult := &packages.SinglePollDataResultLinked{}
			if err := linkedResult.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("failed to unmarshal SinglePollDataResultLinked: %v", err)
				return
			}
			linkedResult.ShowInfo(&c.printMu, 0)
		case base.ROER_APDU: // ROER_APDU
			c.SafeLog("Data Export Protocol: ROER_APDU")
			// обработка ROER_APDU
		default:
			c.SafeLog("Data Export Protocol: unknown ro_type 0x%04X", roType)
		}
	default:
		c.SafeLog("runPacketListener: Received unknown packet (0x%02X). Ignoring.", firstByte)
		utils.PrintHexDump(&c.printMu, "Unknown packet", data)
	}
}
