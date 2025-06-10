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
			// Если элемент - это массив, добавляем все его элементы
			flatData = append(flatData, v...)
		case map[string]interface{}:
			// Если элемент - это объект, добавляем его как есть
			flatData = append(flatData, v)
		default:
			// Для всех остальных типов (строки, числа и т.д.) добавляем как есть
			flatData = append(flatData, v)
		}
	}

	if len(flatData) == 0 {
		// Если после обработки массив пустой, возвращаем пустой массив в JSON
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
			c.SafeLog("JSON: %s", string(flatData))

			if c.receiverConn != nil {
				connMu.Lock()
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
