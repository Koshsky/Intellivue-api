package client

import (
	"time"
)

const (
	NUMERIC_POLL_INTERVAL  = 1500 * time.Millisecond
	ALARM_POLL_INTERVAL    = 1500 * time.Millisecond
	WAVEFORM_POLL_INTERVAL = 5000 * time.Second
)

func (c *ComputerClient) CollectNumerics() {
	var invokeID uint16 = 1

	pollInterval := NUMERIC_POLL_INTERVAL
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

func (c *ComputerClient) CollectAlarms() {
	var invokeID uint16 = 100

	pollInterval := ALARM_POLL_INTERVAL
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
