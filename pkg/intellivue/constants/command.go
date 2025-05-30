package constants

// Command Types
const (
	CMD_EVENT_REPORT           uint16 = 0x0000 // Event Report
	CMD_CONFIRMED_EVENT_REPORT uint16 = 0x0001 // Confirmed Event Report
	CMD_GET                    uint16 = 0x0003 // Get
	CMD_SET                    uint16 = 0x0004 // Set
	CMD_CONFIRMED_SET          uint16 = 0x0005 // Confirmed Set
	CMD_CONFIRMED_ACTION       uint16 = 0x0007 // Confirmed Action
)
