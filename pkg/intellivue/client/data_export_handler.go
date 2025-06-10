package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
)

const (
	RORLS_FIRST              uint8 = 0x0001 // set in the first message
	RORLS_NOT_FIRST_NOT_LAST uint8 = 0x0002
	RORLS_LAST               uint8 = 0x0003 // last RORLSapdu, one RORSapdu to follow
)

type rolrsGroup struct {
	invokeID uint16
	results  []*packages.SinglePollDataResultLinked
}

func (c *ComputerClient) StartPacketHandlers() {
	go c.roivHandler()
	go c.rorsHandler()
	go c.rolrsHandler()
	go c.roerHandler()
	go c.jsonHandler()
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
			if len(data) < 6 {
				continue
			}
			roType := binary.BigEndian.Uint16(data[4:6])
			if roType == base.ROER_APDU {
				select {
				case c.assocResponseChan <- struct{}{}:
				case <-c.ctx.Done():
					return
				}
			}
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
				select {
				case c.mdsCreateChan <- struct{}{}:
				case <-c.ctx.Done():
					c.SafeLog("Context done while sending MDSCreateEvent signal")
				}
				createResult := packages.NewMdsCreateEventResult()
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

			jsonBytes, err := json.Marshal(result.GetPollInfoList())
			if err != nil {
				c.SafeLog("Failed to marshal PollInfoList to JSON: %v", err)
				continue
			}
			c.jsonChan <- jsonBytes

		}
	}
}

func (c *ComputerClient) rolrsHandler() {
	c.SafeLog("ROLRS handler started")
	groups := make(map[uint16]*rolrsGroup)

	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("ROLRS handler: context done")
			return
		case data := <-c.rolrsChan:
			c.SafeLog("ROLRS handler: received data, length: %d", len(data))
			c.SafeLog("Data Export Protocol: ROLRS_APDU")

			if len(data) < 12 {
				c.SafeLog("ROLRS_APDU packet too short")
				continue
			}
			invokeID := binary.BigEndian.Uint16(data[10:12])

			result := &pollDataResultLinkedWrapper{result: &packages.SinglePollDataResultLinked{}}
			if err := result.UnmarshalBinary(bytes.NewReader(data)); err != nil {
				c.SafeLog("Failed to unmarshal ROLRS_APDU: %v", err)
				continue
			}

			rolrsType := result.result.ROLRSapdu.LinkedID.State
			c.SafeLog("ROLRS type: %d for invoke_id: %d", rolrsType, invokeID)

			group, exists := groups[invokeID]
			if !exists {
				group = &rolrsGroup{
					invokeID: invokeID,
					results:  make([]*packages.SinglePollDataResultLinked, 0),
				}
				groups[invokeID] = group
			}

			group.results = append(group.results, result.result)

			if rolrsType == RORLS_LAST {
				c.SafeLog("Received last ROLRS packet for invoke_id: %d", invokeID)

				var combinedAttributes []interface{}
				for _, res := range group.results {
					if res.PollMdibDataReply.PollInfoList != nil {
						jsonBytes, err := json.Marshal(res.PollMdibDataReply.PollInfoList)
						if err != nil {
							c.SafeLog("Failed to marshal PollInfoList: %v", err)
							continue
						}

						var nestedData []interface{}
						if err := json.Unmarshal(jsonBytes, &nestedData); err != nil {
							c.SafeLog("Failed to unmarshal for flattening: %v", err)
							continue
						}

						for _, item := range nestedData {
							if itemList, ok := item.([]interface{}); ok {
								combinedAttributes = append(combinedAttributes, itemList...)
							}
						}
					}
				}

				jsonBytes, err := json.Marshal(combinedAttributes)
				if err != nil {
					c.SafeLog("Failed to marshal combined PollInfoList to JSON: %v", err)
					continue
				}
				c.jsonChan <- jsonBytes

				delete(groups, invokeID)
			}
		}
	}
}
