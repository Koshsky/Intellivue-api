package client

import (
	"context"
	"time"
)

const (
	NUMERIC_POLL_INTERVAL  = 15 * time.Second
	ALARM_POLL_INTERVAL    = 10 * time.Second
	WAVEFORM_POLL_INTERVAL = 5 * time.Second
)

func (c *ComputerClient) CollectNumerics(ctx context.Context) {
	var invokeID uint16 = 1

	pollInterval := NUMERIC_POLL_INTERVAL
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	if err := c.SendPollNumericAction(invokeID); err != nil {
		c.SafeLog("Failed to process SinglePollDataResultLinked: %v", err)
	}
	invokeID++

	for {
		select {
		case <-ctx.Done():
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

func (c *ComputerClient) CollectAlarms(ctx context.Context) {
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
		case <-ctx.Done():
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

func (c *ComputerClient) CollectWaveforms(ctx context.Context) {
	c.SafeLog("Starting CollectWaveforms goroutine...")
	// Пример цикла, который будет остановлен при отмене контекста
	for {
		select {
		case <-ctx.Done():
			c.SafeLog("CollectWaveforms: Context canceled, goroutine ending.")
			return
		default:
			time.Sleep(WAVEFORM_POLL_INTERVAL)
		}
	}
}
