package intellivue

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
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

	// Канал для сигнализации о завершении рукопожатия ассоциации
	associationDone chan error
}

// NewComputerClient creates a new instance of ComputerClient
func NewComputerClient(host, port string) *ComputerClient {
	log.Printf("Initializing ComputerClient with host=%s, port=%s", host, port)
	return &ComputerClient{
		host:            host,
		port:            port,
		streams:         make(map[string]*DataStream),
		associationDone: make(chan error, 1),
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
	go c.listenForResponses(ctx, c.associationDone)

	return nil
}

// Connect устанавливает UDP соединение с сервером и выполняет запрос ассоциации.
// Теперь использует EstablishUDPConnection для установки соединения.
func (c *ComputerClient) Connect(ctx context.Context) error {
	// Устанавливаем базовое UDP соединение и запускаем слушателя
	if err := c.EstablishUDPConnection(ctx); err != nil {
		return fmt.Errorf("ошибка при установке UDP соединения: %w", err)
	}

	// Установка ассоциации: создание и отправка AssociationRequest
	msg := NewAssocReqMessage()
	data, err := msg.MarshalBinary()
	if err != nil {
		c.Close()
		return fmt.Errorf("ошибка при создании сообщения AssociationRequest: %v", err)
	}

	log.Printf("\n=== Отправляемое сообщение ===\n")
	PrintHexDump("AssociationRequest", data)

	if err := c.sendData(data); err != nil {
		c.Close()
		return fmt.Errorf("ошибка при отправке AssociationRequest: %v", err)
	}

	log.Println("AssociationRequest отправлен. Ожидание ответа...")

	// Ожидаем завершения рукопожатия ассоциации
	select {
	case <-ctx.Done():
		// Если контекст отменен во время ожидания ответа ассоциации
		c.Close()
		return ctx.Err()
	case err := <-c.associationDone:
		// Получен сигнал о завершении рукопожатия ассоциации
		// associationDone будет закрыт, когда handleAssociationHandshake завершится
		if err != nil {
			// Если рукопожатие завершилось с ошибкой (например, Abort)
			c.Close()
			return fmt.Errorf("рукопожатие ассоциации завершилось с ошибкой: %w", err)
		}
		// Рукопожатие ассоциации успешно завершено (получен Association Response или MDS Create Request)
		log.Println("Рукопожатие ассоциации успешно завершено.")

		// После успешной ассоциации запускаем рутину для прослушивания всех пакетов данных в отдельной горутине.
		// Connect теперь возвращает nil и не ждет завершения этого слушателя.
		dataListenerErrChan := make(chan error, 1) // Канал ошибок слушателя данных (может использоваться для внешней реакции)
		go c.runDataPacketListener(ctx, c.conn.(*net.UDPConn), dataListenerErrChan)

		// TODO: Решить, как использовать dataListenerErrChan для реакции вне Connect
		// Например, передать его в Run или обрабатывать в основной логике приложения.

		return nil // Connect успешно завершен после установления ассоциации и запуска слушателя данных
	}
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

// sendData отправляет данные по UDP
func (c *ComputerClient) sendData(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("соединение не установлено")
	}

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

// listenForResponses запускает горутину для обработки входящих пакетов во время ассоциации.
// Эта функция предназначена для начального этапа обмена сообщениями (рукопожатие ассоциации).
// Теперь принимает канал для сигнализации о завершении ассоциации.
func (c *ComputerClient) listenForResponses(ctx context.Context, associationDone chan error) error {
	// Используем контекст для остановки прослушивания
	go func() {
		<-ctx.Done()
		if c.conn != nil {
			c.conn.Close() // Закрытие соединения приведет к ошибке чтения и завершению цикла в handleAssociationHandshake
		}
	}()

	errChan := make(chan error, 1)
	// Запускаем горутину для обработки пакетов рукопожатия ассоциации
	go c.handleAssociationHandshake(ctx, c.conn.(*net.UDPConn), errChan, associationDone)

	// listenForResponses сама по себе не блокируется на ожидании.
	// Обработка входящих пакетов и сигналы об ошибках/завершении
	// происходят в handleAssociationHandshake.
	return nil
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
	PrintHexDump("Send MDS Create Result", resultBytes)

	return nil
}

func (c *ComputerClient) handleAssociationAbort(data []byte) error {
	log.Println("Обработка Association Abort")
	return fmt.Errorf("ассоциация прервана сервером")
}

// handleAssociationHandshake обрабатывает входящие пакеты во время рукопожатия ассоциации.
// Она завершается после получения определенного типа пакета (Response, Create, Abort).
// Принимает контекст для возможности внешней отмены и канал для сигнализации о завершении.
func (c *ComputerClient) handleAssociationHandshake(ctx context.Context, conn *net.UDPConn, errChan chan error, associationDone chan error) {
	defer close(errChan)
	// Убедимся, что канал associationDone будет закрыт в любом случае при выходе из рутины
	// Это нужно, даже если произошла ошибка, чтобы вызывающая сторона не ждала вечно.
	// При успехе мы отправим nil, при ошибке - саму ошибку.
	defer func() {
		if r := recover(); r != nil {
			// Если произошла паника, передаем ее как ошибку
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic occurred: %v", r)
			}
			select {
			case associationDone <- err:
			default:
				log.Printf("Предупреждение: Не удалось отправить ошибку паники по каналу associationDone: %v", err)
			}
		} else {
			// Если рутина завершилась без паники (успех или обработанная ошибка),
			// проверяем, не была ли уже отправлена ошибка через errChan.
			// Если нет, сигнализируем об успехе (nil).
			select {
			case err := <-errChan:
				// Ошибка уже была отправлена, отправляем ее же через associationDone
				select {
				case associationDone <- err:
				default:
					log.Printf("Предупреждение: Не удалось отправить ошибку по каналу associationDone: %v", err)
				}
			default:
				// Ошибки не было, сигнализируем об успехе
				select {
				case associationDone <- nil:
				default:
					log.Println("Предупреждение: Не удалось отправить сигнал успеха по каналу associationDone")
				}
			}
		}
	}()

	buffer := make([]byte, 4096)

	for {
		// Проверяем контекст на отмену перед чтением
		select {
		case <-ctx.Done():
			log.Println("handleAssociationHandshake: Контекст отменен, завершение рутины.")
			return
		default:
			// Продолжаем чтение
		}

		// Читаем пакеты с установленным таймаутом (если есть)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Таймаут чтения - ожидаемое поведение при отсутствии трафика, продолжаем цикл
				continue
			} else {
				log.Printf("handleAssociationHandshake: Ошибка при чтении UDP: %v", err)
				errChan <- fmt.Errorf("ошибка при чтении UDP: %v", err)
				return // Выходим при других ошибках
			}
		}

		if n == 0 {
			continue // Пропускаем пустые пакеты
		}

		data := buffer[:n]
		firstByte := data[0]
		PrintHexDump(fmt.Sprintf("Получен пакет (0x%02X)", firstByte), data)

		// Диспетчеризация обработки пакетов по первому байту (MessageType)
		switch firstByte {
		case 0x0E:
			// Association Response
			if err := c.handleAssociationResponse(data); err != nil {
				errChan <- err
				return
			}
			// Сигнализируем об успешном завершении рукопожатия
			select {
			case associationDone <- nil:
			default:
				log.Println("Предупреждение: Не удалось отправить сигнал успеха по каналу associationDone (Response)")
			}
			return // Выходим после обработки Association Response
		case 0xE1:
			// MDS Create Request
			if err := c.handleMDSCreateRequest(data); err != nil {
				errChan <- err
				return
			}
			// Сигнализируем об успешном завершении рукопожатия
			select {
			case associationDone <- nil:
			default:
				log.Println("Предупреждение: Не удалось отправить сигнал успеха по каналу associationDone (Create)")
			}
			return // Выходим после обработки MDS Create Request
		case 0x19:
			// Association Abort
			if err := c.handleAssociationAbort(data); err != nil {
				errChan <- err
				return
			}
			// Сигнализируем об ошибке прерывания ассоциации
			select {
			case associationDone <- fmt.Errorf("ассоциация прервана сервером"): // TODO: Уточнить тип ошибки
			default:
				log.Println("Предупреждение: Не удалось отправить сигнал Abort по каналу associationDone")
			}
			return // Выходим после обработки Association Abort
		default:
			log.Printf("Неизвестный тип сообщения во время рукопожатия ассоциации: 0x%02X\n", firstByte)
			// В случае неизвестного пакета на этапе ассоциации,
			// мы не получили подтверждения успеха или ошибки ассоциации.
			// В зависимости от спецификации, можно либо считать это ошибкой,
			// либо игнорировать и продолжать ждать ожидаемые ответы.
			// Пока просто логируем и продолжаем ждать.
		}
	}
}

// Run запускает основные рутины сбора данных клиента.
// Принимает контекст для управления их жизненным циклом.
// Эта функция должна быть вызвана только после успешного Connect.
func (c *ComputerClient) Run(ctx context.Context) error {
	// Connect() должен быть вызван отдельно для установления соединения и ассоциации перед вызовом Run.

	log.Printf("Starting ComputerClient data collection routines (host=%s, port=%s)", c.host, c.port)

	// Запускаем различные рутины сбора данных
	go c.collectVitalSigns(ctx)
	go c.collectWaveforms(ctx)
	go c.collectAlarms(ctx)

	log.Println("All data collection routines started successfully")
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
func (c *ComputerClient) UpdateStream(streamID string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.streams[streamID] = &DataStream{
		ID:    streamID,
		Value: value,
	}
	log.Printf("Updated stream '%s' with new data", streamID)
}

// collectVitalSigns рутина сбора данных жизненных показателей (пример)
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
			// Симуляция сбора данных жизненных показателей
			c.UpdateStream("vitals", map[string]interface{}{
				"heartRate": 75,
				"spO2":      98,
			})
		}
	}
}

// collectWaveforms рутина сбора данных волноформ (пример)
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
			// Симуляция сбора данных волноформ
			c.UpdateStream("waveforms", map[string]interface{}{
				"ecg": []float64{1.0, 1.2, 0.8, 1.1},
			})
		}
	}
}

// collectAlarms рутина сбора данных тревог: отправляет запрос MDSPollAction
// Обработка ответа происходит в handleIncomingPackets.
func (c *ComputerClient) collectAlarms(ctx context.Context) {
	log.Println("Starting alarms collection routine: Sending MDSPollAction periodically")

	// TODO: Определить подходящий интервал для отправки запросов тревог
	ticker := time.NewTicker(30 * time.Second) // Пример: отправлять запрос каждые 30 секунд
	defer ticker.Stop()

	// Отправляем первый запрос сразу при запуске рутины
	c.sendPollActionAlarm()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping alarms collection routine")
			return
		case <-ticker.C:
			// Отправляем запрос MDSPollAction по таймеру
			c.sendPollActionAlarm()
		}
	}
}

// sendPollActionAlarm создает и отправляет пакет MDSPollAction
func (c *ComputerClient) sendPollActionAlarm() {
	mdsPollAction := NewMDSPollAction(2, 1, NOM_MOC_VMO_AL_MON) // TODO: InvokeID должен быть уникальным
	data, err := mdsPollAction.MarshalBinary()
	if err != nil {
		log.Printf("Ошибка при маршалинге MDSPollAction: %v", err)
		return
	}

	if err := c.sendData(data); err != nil {
		log.Printf("Ошибка при отправке MDSPollAction: %v", err)
		// TODO: Обработать ошибку отправки, возможно, переподключиться или повторить попытку
		return
	}
}

// runDataPacketListener постоянно читает входящие пакеты и печатает их Hex Dump.
// Запускается после успешного рукопожатия ассоциации.
// Также отслеживает пакеты Association Abort (0x19) для определения разрыва соединения со стороны устройства.
// Сигнализирует об ошибке (включая Abort) через предоставленный канал.
func (c *ComputerClient) runDataPacketListener(ctx context.Context, conn *net.UDPConn, errChan chan error) {
	defer close(errChan)
	log.Println("Запуск рутины прослушивания пакетов данных...")
	buffer := make([]byte, 4096)

	for {
		// Проверяем контекст на отмену перед чтением
		select {
		case <-ctx.Done():
			log.Println("runDataPacketListener: Контекст отменен, завершение рутины.")
			// Сигнализируем вызывающей стороне, что рутина завершена из-за отмены контекста (с ошибкой контекста)
			select {
			case errChan <- ctx.Err():
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку контекста по каналу errChan")
			}
			return
		default:
			// Продолжаем чтение
		}

		// Читаем пакеты.
		// Таймаут может быть установлен на уровне соединения в EstablishUDPConnection, если требуется,
		// или можно использовать Select с таймаутом, если нужно периодически делать что-то еще в этом цикле.
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			// Проверяем, не является ли ошибка результатом закрытия соединения из-за отмены контекста или явного закрытия (например, при Abort)
			if ctx.Err() != nil || c.conn == nil {
				log.Println("runDataPacketListener: Ошибка чтения после отмены контекста или закрытия соединения, завершение.")
				// Сигнализируем об ошибке, если она не вызвана отменой контекста (т.е. если c.conn == nil из-за Abort или явного Close)
				if ctx.Err() == nil {
					select {
					case errChan <- fmt.Errorf("соединение закрыто: %w", err):
					default:
						log.Println("Предупреждение: Не удалось отправить ошибку закрытия соединения по каналу errChan")
					}
				}
				return
			}
			// Игнорируем ошибки таймаута, если они настроены на уровне соединения
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue // Просто продолжаем ждать
			}
			// Логируем другие ошибки чтения и, возможно, завершаем рутину
			log.Printf("runDataPacketListener: Ошибка при чтении UDP: %v", err)
			// Сигнализируем об ошибке чтения
			select {
			case errChan <- fmt.Errorf("ошибка при чтении UDP: %w", err):
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку чтения по каналу errChan")
			}
			return
		}

		if n == 0 {
			continue // Пропускаем пустые пакеты
		}

		// Печатаем полученный пакет в HEX-формате и проверяем тип
		data := buffer[:n]
		firstByte := data[0]

		switch firstByte {
		case 0x19:
			// Association Abort
			log.Println("runDataPacketListener: Получен пакет Association Abort (0x19). Соединение разорвано устройством.")
			PrintHexDump("Association Abort Packet", data)
			// Закрываем соединение, чтобы остановить рутину чтения и другие рутины Run.
			c.Close()
			// Сигнализируем об ошибке Abort
			select {
			case errChan <- fmt.Errorf("соединение разорвано устройством (Association Abort)"):
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку Association Abort по каналу errChan")
			}
			return // Завершаем рутину прослушивания
		default:
			// Для всех остальных пакетов данных просто печатаем Hex Dump
			PrintHexDump(fmt.Sprintf("Получен пакет данных (0x%02X)", firstByte), data)
			// TODO: Добавить логику обработки конкретных типов пакетов данных, если необходимо
		}

	}
}
