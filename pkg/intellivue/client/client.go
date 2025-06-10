package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
)

type ComputerClient struct {
	conn         *net.UDPConn
	receiverConn *net.UDPConn
	monitorAddr  string
	receiverAddr string
	printMu      sync.Mutex

	roivChan      chan []byte
	rorsChan      chan []byte
	rolrsChan     chan []byte
	roerChan      chan []byte
	jsonChan      chan []byte
	mdsCreateChan chan struct{}

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	isAssociationDone bool
	closeOnce         sync.Once
	assocResponseChan chan struct{}
}

func NewComputerClient(monitorAddr, receiverAddr string) *ComputerClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &ComputerClient{
		monitorAddr:       monitorAddr,
		receiverAddr:      receiverAddr,
		roivChan:          make(chan []byte, 100),
		rorsChan:          make(chan []byte, 100),
		rolrsChan:         make(chan []byte, 100),
		roerChan:          make(chan []byte, 100),
		jsonChan:          make(chan []byte, 100),
		ctx:               ctx,
		cancel:            cancel,
		assocResponseChan: make(chan struct{}, 1),
		mdsCreateChan:     make(chan struct{}, 1),
	}
}

func (c *ComputerClient) Connect(ctx context.Context) error {
	// Используем переданный контекст только для установки соединения
	connectCtx, connectCancel := context.WithTimeout(ctx, 5*time.Second)
	defer connectCancel()

	if err := c.connectUDP(); err != nil {
		return fmt.Errorf("failed to establish UDP connection: %w", err)
	}

	// Запускаем обработчики пакетов
	c.StartPacketHandlers()

	// Отправляем Association Request
	if err := c.sendAssociationRequest(); err != nil {
		return fmt.Errorf("failed to send Association Request: %w", err)
	}

	// Ждем ответа на Association Request
	assocReceived := false
	mdsReceived := false

	for !assocReceived || !mdsReceived {
		select {
		case <-c.assocResponseChan:
			assocReceived = true
			c.SafeLog("Association Response received")
		case <-c.mdsCreateChan:
			mdsReceived = true
			c.SafeLog("MDSCreateEvent received")
		case <-connectCtx.Done():
			return fmt.Errorf("timeout waiting for Association Response or MDSCreateEvent")
		}
	}

	c.SafeLog("Association completed successfully")
	return nil
}

func (c *ComputerClient) Wait() {
	c.wg.Wait()
}

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

		if c.receiverConn != nil {
			c.receiverConn.Close()
			c.receiverConn = nil
			c.SafeLog("Receiver connection closed.")
		}

		c.SafeLog("Client shutdown procedure completed.")
	})
	return nil
}

func (c *ComputerClient) sendData(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("connection not established")
	}
	if !c.isAssociationDone {
		return fmt.Errorf("association not done, cannot send data")
	}
	_, err := c.conn.Write(data)
	return err
}

func (c *ComputerClient) runPacketListener() {
	defer c.wg.Done()
	c.SafeLog("Starting Packet Listener goroutine...")

	buffer := make([]byte, 4096)
	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("Packet Listener: Context canceled, goroutine ending.")
			return
		default:
			c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, _, err := c.conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				c.SafeLog("runPacketListener: failed to read UDP: %v", err)
				return
			}
			if n == 0 {
				continue
			}
			firstByte := buffer[0]
			switch firstByte {
			case base.AC_SPDU_SI:
				c.SafeLog("runPacketListener: Received Association Response (0x%02X)", firstByte)
				c.isAssociationDone = true
				c.assocResponseChan <- struct{}{}
			case base.AB_SPDU_SI:
				c.SafeLog("runPacketListener: Received Association Abort (0x%02X).", firstByte)
				c.isAssociationDone = false
				c.Close()
				return
			case base.RF_SPDU_SI:
				c.SafeLog("runPacketListener: Received Association Refuse (0x%02X).", firstByte)
				c.isAssociationDone = false
				c.Close()
				return
			case 0xE1:
				c.SafeLog("runPacketListener: Received Data Export Protocol packet (0x%02X).", firstByte)
				c.handleDataExportPacket(buffer[:n])
			default:
				c.SafeLog("runPacketListener: Received unknown packet (0x%02X). Ignoring.", firstByte)
			}
		}
	}
}

func (c *ComputerClient) SafeLog(format string, args ...interface{}) {
	c.printMu.Lock()
	defer c.printMu.Unlock()
	log.Printf(format, args...)
}

func (c *ComputerClient) connectUDP() error {
	c.SafeLog("Trying to establish UDP connection with %s", c.monitorAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", c.monitorAddr)
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %w", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("error creating UDP connection: %w", err)
	}
	c.conn = conn
	c.SafeLog("UDP connection established.")

	receiverAddr, err := net.ResolveUDPAddr("udp", c.receiverAddr)
	if err != nil {
		c.Close()
		return fmt.Errorf("error resolving receiver UDP address: %w", err)
	}
	receiverConn, err := net.DialUDP("udp", nil, receiverAddr)
	if err != nil {
		c.Close()
		return fmt.Errorf("error creating receiver UDP connection: %w", err)
	}
	c.receiverConn = receiverConn
	c.SafeLog("Receiver UDP connection established.")

	c.wg.Add(1)
	go c.runPacketListener()
	c.wg.Add(1)

	return nil
}

func (c *ComputerClient) sendAssociationRequest() error {
	assocReq, err := packages.NewAssociationRequest().MarshalBinary()
	if err != nil {
		c.Close()
		return fmt.Errorf("error creating AssociationRequest message: %w", err)
	}
	_, err = c.conn.Write(assocReq)
	if err != nil {
		c.Close()
		return fmt.Errorf("error sending AssociationRequest: %w", err)
	}
	c.SafeLog("AssociationRequest sent. Waiting for response...")
	return nil
}
