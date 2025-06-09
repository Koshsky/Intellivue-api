package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
)

func (c *ComputerClient) StartPacketHandlers() {
	go c.roivHandler()
	go c.rorsHandler()
	go c.rolrsHandler()
	go c.roerHandler()
}

func (c *ComputerClient) handleDataExportPacket(data []byte) {
	if len(data) < 6 {
		c.SafeLog("Data Export Protocol: too short packet to determine ro_type")
		return
	}

	roType := binary.BigEndian.Uint16(data[4:6])
	c.SafeLog("Data Export Protocol: received packet with ro_type: 0x%04X, length: %d", roType, len(data))

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	switch roType {
	case base.ROIV_APDU:
		select {
		case c.roivChan <- dataCopy:
		case <-c.ctx.Done():
			c.SafeLog("Data Export Protocol: context done while sending to roivChan")
		}
	case base.RORS_APDU:
		select {
		case c.rorsChan <- dataCopy:
		case <-c.ctx.Done():
			c.SafeLog("Data Export Protocol: context done while sending to rorsChan")
		}
	case base.ROLRS_APDU:
		select {
		case c.rolrsChan <- dataCopy:
		case <-c.ctx.Done():
			c.SafeLog("Data Export Protocol: context done while sending to rolrsChan")
		}
	default:
		c.SafeLog("Data Export Protocol: unknown ro_type: 0x%04X", roType)
	}
}

func (c *ComputerClient) roerHandler() {
	c.SafeLog("ROER handler started")
	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("ROER handler: context done")
			return
		case data := <-c.roerChan:
			c.SafeLog("ROER handler: received data, length: %d", len(data))
		}
	}
}

func (c *ComputerClient) roivHandler() {
	var connMu sync.Mutex
	c.SafeLog("ROIV handler started")
	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("ROIV handler: context done")
			return
		case data := <-c.roivChan:
			c.SafeLog("ROIV handler: received data, length: %d", len(data))
			if len(data) < 12 {
				c.SafeLog("ROIV_APDU packet too short")
				continue
			}

			commandType := base.CMDType(binary.BigEndian.Uint16(data[10:12]))
			c.SafeLog("ROIV_APDU: commandType: 0x%04X: length: %d", commandType, len(data))

			switch commandType {
			case base.CMD_CONFIRMED_EVENT_REPORT:
				c.SafeLog("Received MDSCreateEvent")
				if c.mdsCreateHandler != nil {
					c.mdsCreateHandler()
				}
				createResult := packages.NewMDSCreateResult()
				resultBytes, err := createResult.MarshalBinary()
				if err != nil {
					c.SafeLog("Failed to marshal MDSCreateResult: %v", err)
					continue
				}
				connMu.Lock()
				if _, err := c.conn.Write(resultBytes); err != nil {
					c.SafeLog("Failed to send MDSCreateResult: %v", err)
					connMu.Unlock()
					continue
				}
				connMu.Unlock()
				c.SafeLog("MDSCreateResult sent")
			default:
				c.SafeLog("ROIV_APDU: unknown commandType: 0x%04X", commandType)
			}
		}
	}
}

func (c *ComputerClient) rorsHandler() {
	var connMu sync.Mutex
	c.SafeLog("RORS handler started")
	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("RORS handler: context done")
			return
		case data := <-c.rorsChan:
			c.SafeLog("RORS handler: received data, length: %d", len(data))
			c.SafeLog("Data Export Protocol: RORS_APDU")
			result := &pollDataResultWrapper{result: &packages.SinglePollDataResult{}}
			if err := result.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("Failed to unmarshal RORS_APDU: %v", err)
				continue
			}

			jsonBytes, err := json.MarshalIndent(result.GetPollInfoList(), "", "  ")
			if err != nil {
				c.SafeLog("Failed to marshal PollInfoList to JSON: %v", err)
				continue
			}

			if c.receiverConn != nil {
				connMu.Lock()
				if _, err := c.receiverConn.Write(jsonBytes); err != nil {
					c.SafeLog("Failed to send data to receiver: %v", err)
					connMu.Unlock()
					continue
				}
				connMu.Unlock()
				c.SafeLog("Data sent to receiver %s", c.receiverAddr)
			}
		}
	}
}

func (c *ComputerClient) rolrsHandler() {
	var connMu sync.Mutex
	c.SafeLog("ROLRS handler started")
	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("ROLRS handler: context done")
			return
		case data := <-c.rolrsChan:
			c.SafeLog("ROLRS handler: received data, length: %d", len(data))
			c.SafeLog("Data Export Protocol: ROLRS_APDU")
			result := &pollDataResultLinkedWrapper{result: &packages.SinglePollDataResultLinked{}}
			if err := result.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("Failed to unmarshal ROLRS_APDU: %v", err)
				continue
			}

			jsonBytes, err := json.MarshalIndent(result.GetPollInfoList(), "", "  ")
			if err != nil {
				c.SafeLog("Failed to marshal PollInfoList to JSON: %v", err)
				continue
			}

			if c.receiverConn != nil {
				connMu.Lock()
				if _, err := c.receiverConn.Write(jsonBytes); err != nil {
					c.SafeLog("Failed to send data to receiver: %v", err)
					connMu.Unlock()
					continue
				}
				connMu.Unlock()
				c.SafeLog("Data sent to receiver %s", c.receiverAddr)
			}
		}
	}
}
