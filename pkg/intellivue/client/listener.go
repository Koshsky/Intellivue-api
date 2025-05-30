package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

// runPacketListener постоянно читает входящие пакеты и диспетчеризирует их.
// Запускается единожды в EstablishUDPConnection.
// Сигнализирует о всех причинах завершения рутины (успех ассоциации, ошибки, таймаут, Abort/Refuse)
// через предоставленный канал listenerOutcome.
func (c *ComputerClient) runPacketListener(ctx context.Context, conn *net.UDPConn, listenerOutcome chan error) {
	defer close(listenerOutcome)

	log.Println("Запуск основной рутины прослушивания пакетов...")
	buffer := make([]byte, 4096) // Буфер для чтения пакетов

	isAssociationDone := false

	inactivityTimeout := 10 * time.Second
	timer := time.NewTimer(inactivityTimeout)

	defer timer.Stop() // Останавливаем таймер при выходе из рутины

	for {
		select {
		case <-ctx.Done():
			log.Println("runPacketListener: Контекст отменен, завершение рутины.")
			select {
			case listenerOutcome <- ctx.Err():
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку контекста по каналу listenerOutcome")
			}
			c.Close() // Закрываем соединение при отмене контекста
			return    // Завершаем рутину

		case assocErr := <-c.listenerOutcome:
			// Этот случай не должен срабатывать в нормальном цикле рутины, т.к. рутина сама пишет в этот канал.
			// Это может произойти только если что-то пошло не так или канал используется не по назначению.
			log.Printf("runPacketListener: Неожиданное чтение из собственного listenerOutcome канала: %v", assocErr)
			select {
			case listenerOutcome <- fmt.Errorf("неожиданное чтение из listenerOutcome: %w", assocErr):
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку неожиданного чтения по каналу listenerOutcome")
			}
			c.Close()
			return

		case <-timer.C:
			// Сработал таймер неактивности
			log.Printf("runPacketListener: Таймаут неактивности (%s) истек, закрытие соединения.", inactivityTimeout)
			c.Close()
			// Сигнализируем об ошибке неактивности через listenerOutcome
			select {
			case listenerOutcome <- fmt.Errorf("соединение закрыто из-за неактивности (%s)", inactivityTimeout):
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку таймаута по каналу listenerOutcome")
			}
			return // Завершаем рутину

		default:
			// Неблокирующее чтение из каналов завершилось. Теперь переходим к блокирующему чтению из сети.
		}

		// Устанавливаем таймаут для чтения, чтобы не блокироваться вечно и дать возможность сработать контексту или таймеру.
		// Таймаут чтения должен быть меньше таймаута неактивности, чтобы таймер неактивности мог сработать.
		// Установка таймаута 0 или отрицательного значения сбрасывает его.
		readTimeout := inactivityTimeout / 2 // Пример: половина таймаута неактивности
		if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
			log.Printf("runPacketListener: Ошибка установки таймаута чтения: %v", err)
			// В случае ошибки установки таймаута, логируем и продолжаем. Это может повлиять на своевременное
			// обнаружение таймаута неактивности, но не должно приводить к завершению рутины сразу.
			// TODO: Решить, является ли ошибка установки таймаута чтения критической.
		}

		// Читаем пакеты. Блокируется до получения данных, истечения таймаута чтения или закрытия соединения.
		n, _, err := conn.ReadFromUDP(buffer)

		// Сбрасываем таймаут чтения после попытки чтения (успешной или нет), чтобы не влиять на следующий цикл select
		conn.SetReadDeadline(time.Time{}) // Сброс таймаута

		if err != nil {
			// Проверяем тип ошибки
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Это таймаут чтения. Ожидаемо, если нет входящих данных.
				// Цикл select вернется к ожиданию на каналах (ctx.Done, timer.C).
				continue
			} else if ctx.Err() != nil {
				// Ошибка чтения вызвана отменой контекста. Рутина завершится в верхнем select case.
				log.Println("runPacketListener: Ошибка чтения после отмены контекста.")
				continue // Позволяем верхнему select case обработать ctx.Done()
			} else if c.conn == nil {
				// Ошибка чтения, потому что соединение уже закрыто (например, в результате Abort/Refuse, обработанного ниже)
				log.Println("runPacketListener: Ошибка чтения, соединение закрыто.") // Рутина завершится после обработки Abort/Refuse ниже или в следующем цикле, когда select наткнется на закрытый канал.
				continue
			} else {
				// Неожиданная ошибка чтения
				log.Printf("runPacketListener: Неожиданная ошибка при чтении UDP: %v", err)
				// Сигнализируем об ошибке чтения через listenerOutcome и завершаем рутину
				select {
				case listenerOutcome <- fmt.Errorf("ошибка при чтении UDP: %w", err):
				default:
					log.Println("Предупреждение: Не удалось отправить ошибку чтения по каналу listenerOutcome")
				}
				c.Close() // Закрываем соединение при неожиданной ошибке
				return
			}
		}

		if n == 0 {
			// Пустые пакеты игнорируем, но сбрасываем таймер неактивности
			timer.Reset(inactivityTimeout)
			continue
		}

		// При успешном чтении пакета (n > 0) сбрасываем таймер неактивности
		timer.Reset(inactivityTimeout)

		// Копируем полученные данные, чтобы они не были перезаписаны следующим чтением
		data := make([]byte, n)
		copy(data, buffer[:n])

		firstByte := data[0]

		// Диспетчеризация обработки пакетов по первому байту (MessageType)
		switch firstByte {
		case 0x0E:
			// Association Response - Ожидается только во время рукопожатия
			if !isAssociationDone { // Обрабатываем только если еще не завершено рукопожатие
				log.Println("runPacketListener: Получен Association Response (0x0E)")
				select {
				case listenerOutcome <- nil:
				default:
					log.Println("Предупреждение: Не удалось отправить сигнал nil по каналу listenerOutcome после получения Association Response")
				}
				isAssociationDone = true
			} else {
				log.Printf("runPacketListener: Получен неожиданный Association Response (0x0E) после завершения рукопожатия. Игнорируем.")
			}
		case 0x19:
			// Association Abort - Может прийти в любое время
			log.Println("runPacketListener: Получен пакет Association Abort (0x19).")
			c.Close()
			return // Завершаем рутину прослушивания при Abort
		case 0x0C:
			// Association Refuse - Может прийти во время или после рукопожатия
			log.Println("runPacketListener: Получен пакет Association Refuse (0x0C).")
			// Закрываем соединение
			c.Close()
			// Сигнализируем об ошибке через listenerOutcome
			select {
			case listenerOutcome <- fmt.Errorf("получен пакет Association Refuse (0x0C), соединение закрыто"):
				// Успешно отправили
			default:
				log.Println("Предупреждение: Не удалось отправить ошибку Association Refuse по каналу listenerOutcome")
			}
			return // Завершаем рутину прослушивания при Refuse
		case 0xE1:
			// Data Export Protocol - Ожидается после завершения рукопожатия
			if isAssociationDone { // Обрабатываем только если рукопожатие завершено
				log.Println("runPacketListener: Получен пакет Data Export Protocol (0xE1).")
				// handleDataExportPacket пока только логирует, но здесь будет дальнейшая диспетчеризация.
				if err := c.handleDataExportPacket(ctx, data); err != nil {
					log.Printf("runPacketListener: Ошибка при обработке Data Export Protocol пакета: %v", err)
					// Сигнализируем об ошибке обработки пакета данных. Не обязательно завершать рутину,
					// если это не критическая ошибка парсинга или внутренней логики.
					select {
					case listenerOutcome <- fmt.Errorf("ошибка обработки Data Export Protocol пакета: %w", err):
						// Успешно отправили
					default:
						log.Println("Предупреждение: Не удалось отправить ошибку обработки Data Export Protocol по каналу listenerOutcome")
					}
					// TODO: Решить, нужно ли завершать рутину при ошибке обработки пакета данных.
				}
			} else {
				log.Printf("runPacketListener: Получен неожиданный Data Export Protocol пакет (0xE1) до завершения рукопожатия. Игнорируем.")
			}

		default:
			// Неизвестный или неинтересный пакет
			log.Printf("runPacketListener: Получен неизвестный пакет (0x%02X). Игнорируем.", firstByte)
			utils.PrintHexDump("Неизвестный пакет", data) // TODO: Возможно, стоит убрать в продакшене
		}
	}

}
