package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func (c *ComputerClient) flattenJSONData(data []byte) ([]byte, error) {
	var nestedData []interface{}
	if err := json.Unmarshal(data, &nestedData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var flatData []interface{}
	for _, item := range nestedData {
		switch v := item.(type) {
		case []interface{}:
			flatData = append(flatData, v...)
		case map[string]interface{}:
			flatData = append(flatData, v)
		default:
			flatData = append(flatData, v)
		}
	}

	if len(flatData) == 0 {
		return []byte("[]"), nil
	}

	result, err := json.Marshal(flatData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal flattened data: %w", err)
	}

	return result, nil
}

func (c *ComputerClient) jsonHandler() {
	c.SafeLog("JSON handler started")

	const maxBufferSize = 1024 * 1024 // 1MB
	buffer := make([]byte, 0, maxBufferSize)

	flushTimer := time.NewTimer(100 * time.Millisecond)
	defer flushTimer.Stop()

	var connMu sync.Mutex
	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("JSON handler: context done")
			return
		case data := <-c.jsonChan:
			c.SafeLog("JSON handler: received data, length: %d", len(data))

			flatData, err := c.flattenJSONData(data)
			if err != nil {
				c.SafeLog("Failed to flatten JSON data: %v", err)
				continue
			}

			if len(flatData) <= 2 || string(flatData) == "[]" || string(flatData) == "{}" {
				c.SafeLog("Skipping empty JSON data")
				continue
			}

			// Выводим JSON в читаемом формате
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, flatData, "", "    "); err != nil {
				c.SafeLog("Failed to format JSON: %v", err)
			} else {
				c.SafeLog("JSON to send:\n%s", prettyJSON.String())
			}

			if len(buffer)+len(flatData) > maxBufferSize {
				if c.receiverConn != nil {
					connMu.Lock()
					if _, err := c.receiverConn.Write(buffer); err != nil {
						c.SafeLog("Failed to send buffered data to receiver: %v", err)
					} else {
						c.SafeLog("Buffered data sent to receiver %s", c.receiverAddr)
					}
					connMu.Unlock()
				}
				buffer = buffer[:0]
			}

			buffer = append(buffer, flatData...)
			buffer = append(buffer, '\n')

			if !flushTimer.Stop() {
				select {
				case <-flushTimer.C:
				default:
				}
			}
			flushTimer.Reset(100 * time.Millisecond)

		case <-flushTimer.C:
			if len(buffer) > 0 && c.receiverConn != nil {
				connMu.Lock()
				if _, err := c.receiverConn.Write(buffer); err != nil {
					c.SafeLog("Failed to send buffered data to receiver: %v", err)
				} else {
					c.SafeLog("Buffered data sent to receiver %s", c.receiverAddr)
				}
				connMu.Unlock()
				buffer = buffer[:0]
			}
			flushTimer.Reset(100 * time.Millisecond)
		}
	}
}
