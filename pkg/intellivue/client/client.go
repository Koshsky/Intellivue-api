package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
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
}

func NewComputerClient(addr, port string) *ComputerClient {
	ctx, cancel := context.WithCancel(context.Background())
	client := &ComputerClient{
		serverAddr: addr,
		serverPort: port,
		ctx:        ctx,
		cancel:     cancel,
		sendChan:   make(chan []byte, 10), // если используется
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
		return fmt.Errorf("ошибка разрешения UDP адреса: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("ошибка создания UDP соединения: %w", err)
	}
	c.conn = conn
	c.SafeLog("UDP соединение установлено.")

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
		return fmt.Errorf("ошибка при создании сообщения AssociationRequest: %w", err)
	}

	// Отправка запроса ассоциации
	if err := c.sendData(assocReq); err != nil {
		c.Close() // Закрываем соединение при ошибке отправки
		return fmt.Errorf("ошибка отправки запроса ассоциации: %w", err)
	}

	c.SafeLog("AssociationRequest отправлен. Ожидание ответа...")

	// TODO: Реализовать ожидание Association Response в runPacketListener
	// Сейчас Connect возвращает успешно после отправки запроса, что неверно.
	// Нужен механизм сигнализации от runPacketListener об успешной ассоциации.
	// Временно возвращаем nil для продолжения, но это требует доработки.

	return nil // Временно: успешное возвращение после отправки запроса
}

// Close gracefully закрывает соединение и останавливает все горутины клиента.
func (c *ComputerClient) Close() error {
	c.closeOnce.Do(func() {
		c.SafeLog("Запуск процедуры закрытия клиента...")

		c.cancel()
		c.wg.Wait()
		c.SafeLog("Все горутины завершены.")

		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
			c.SafeLog("Сетевое соединение закрыто.")
		}

		// Закрываем только реально существующий канал
		if c.sendChan != nil {
			close(c.sendChan)
		}

		c.SafeLog("Процедура закрытия клиента завершена.")
	})
	return nil
}

// runPacketSender отправляет пакеты из sendChan на устройство.
func (c *ComputerClient) runPacketSender() {
	defer c.wg.Done() // Уменьшаем счетчик WaitGroup при завершении горутины
	c.SafeLog("Запуск Packet Sender рутины...")

	for {
		select {
		case data, ok := <-c.sendChan:
			if !ok { // Канал закрыт
				c.SafeLog("Packet Sender: Канал отправки закрыт, завершение рутины.")
				return
			}
			// TODO: Добавить мьютекс для записи в соединение, если multiple goroutines могут писать
			if err := c.sendData(data); err != nil {
				c.SafeLog("Ошибка при отправке данных: %v", err)
				// TODO: Обработка ошибок отправки (переподключение?)
			}
		case <-c.ctx.Done():
			c.SafeLog("Packet Sender: Контекст отменен, завершение рутины.")
			return
		}
	}
}

// sendData отправляет байты данных в сетевое соединение.
func (c *ComputerClient) sendData(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("соединение не установлено")
	}
	// TODO: Установить таймаут записи для UDP
	_, err := c.conn.Write(data)
	return err
}

// runPacketListener прослушивает входящие пакеты от устройства по UDP и сразу обрабатывает их.
func (c *ComputerClient) runPacketListener() {
	defer c.wg.Done()
	c.SafeLog("Запуск Packet Listener рутины...")

	conn, ok := c.conn.(*net.UDPConn)
	if !ok {
		c.SafeLog("runPacketListener: Ошибка: соединение не является UDP.")
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
		c.SafeLog("runPacketListener: Получен пакет, первый байт: 0x%02X", data[0])
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
		c.SafeLog("runPacketListener: Получен Association Response (0x0E)")
		c.isAssociationDone = true
	case 0x19:
		c.SafeLog("runPacketListener: Получен пакет Association Abort (0x19).")
		c.Close()
		return
	case 0x0C:
		c.SafeLog("runPacketListener: Получен пакет Association Refuse (0x0C).")
		c.Close()
		return
	case 0xE1:
		c.SafeLog("runPacketListener: Получен Data Export Protocol пакет (0xE1). Обработка...")
		if len(data) < 6 {
			c.SafeLog("Data Export Protocol: слишком короткий пакет для определения ro_type")
			return
		}
		roType := binary.BigEndian.Uint16(data[4:6])
		switch roType {
		case ROIV_APDU: // ROIV_APDU
			c.SafeLog("Data Export Protocol: ROIV_APDU")
			// обработка ROIV_APDU
		case RORS_APDU: // RORS_APDU
			c.SafeLog("Data Export Protocol: RORS_APDU")
			result := &packages.SinglePollDataResult{}
			if err := result.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("Ошибка анмаршалинга SinglePollDataResult: %v", err)
				return
			}
			result.ShowInfo(&c.printMu, 0)
		case ROLRS_APDU: // ROLRS_APDU
			c.SafeLog("Data Export Protocol: ROLRS_APDU")
			linkedResult := &packages.SinglePollDataResultLinked{}
			if err := linkedResult.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("Ошибка анмаршалинга SinglePollDataResultLinked: %v", err)
				return
			}
			linkedResult.ShowInfo(&c.printMu, 0)
		case ROER_APDU: // ROER_APDU
			c.SafeLog("Data Export Protocol: ROER_APDU")
			// обработка ROER_APDU
		default:
			c.SafeLog("Data Export Protocol: Неизвестный ro_type 0x%04X", roType)
		}
	default:
		c.SafeLog("runPacketListener: Получен неизвестный пакет (0x%02X). Игнорируем.", firstByte)
		utils.PrintHexDump(&c.printMu, "Неизвестный пакет", data)
	}
}
