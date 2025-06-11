package client

import (
	"time"
)

func (c *ComputerClient) CollectNumerics(pollInterval time.Duration) {
	var invokeID uint16 = 1
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	if err := c.SendPollNumericAction(invokeID); err != nil {
		c.SafeLog("Failed to send initial SinglePollDataRequest: %v", err)
		return
	}
	invokeID++

	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("CollectNumerics: Context canceled, goroutine ending.")
			return
		case <-ticker.C:
			c.SafeLog("Sending next SinglePollDataRequest by timer.")
			if err := c.SendPollNumericAction(invokeID); err != nil {
				c.SafeLog("Error sending SinglePollDataRequest with InvokeID %d: %v", invokeID, err)
			}
			invokeID++
		}
	}
}

func (c *ComputerClient) CollectAlarms(pollInterval time.Duration) {
	var invokeID uint16 = 100
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	if err := c.SendPollAlarmAction(invokeID); err != nil {
		c.SafeLog("Failed to process SinglePollDataResult: %v", err)
	}
	invokeID++

	for {
		select {
		case <-c.ctx.Done():
			c.SafeLog("CollectAlarms: Context canceled, goroutine ending.")
			return
		case <-ticker.C:
			c.SafeLog("Sending next SinglePollDataRequest by timer.")
			if err := c.SendPollAlarmAction(invokeID); err != nil {
				c.SafeLog("Error sending SinglePollDataRequest with InvokeID %d: %v", invokeID, err)
			}
			invokeID++
		}
	}
}
