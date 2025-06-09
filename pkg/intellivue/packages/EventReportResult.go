package packages

import (
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
)

type EventReportResult struct {
	ManagedObject base.ManagedObjectId `json:"managed_object"`
	CurrentTime   base.RelativeTime    `json:"current_time"`
	EventType     base.OIDType         `json:"event_type"`
	Length        uint16               `json:"length"`
	EventData     []byte               `json:"event_data"`
}
