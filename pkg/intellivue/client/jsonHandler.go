package client

import (
	"encoding/json"
	"fmt"
	"sync"
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

			if c.receiverConn != nil {
				connMu.Lock()
				flatData = append(flatData, '\n')
				if _, err := c.receiverConn.Write(flatData); err != nil {
					c.SafeLog("Failed to send data to receiver: %v", err)
				} else {
					c.SafeLog("Data sent to receiver %s", c.receiverAddr)
				}
				connMu.Unlock()
			}
		}
	}
}
