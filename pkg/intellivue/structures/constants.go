package structures

// APDU Types
const (
	ROIV_APDU  uint16 = 0x0001 // RO Initiate and Vote APDU
	RORS_APDU  uint16 = 0x0002 // RO Reply APDU
	ROER_APDU  uint16 = 0x0003 // RO Error APDU
	ROLRS_APDU uint16 = 0x0005 // RO Linked Reply APDU
)

// Protocol Identification
const (
	NOM_POLL_PROFILE_SUPPORT  OIDType = 0x0001 // id for polling profile (Used)
	NOM_MDIB_OBJ_SUPPORT      OIDType = 0x0102 // supported objects for the active profile (Unused)
	NOM_ATTR_POLL_PROFILE_EXT OIDType = 0xF001 // id for poll profile extensions opt. package (Used)
)

// Event Types
const (
	NOM_NOTI_MDS_CREAT OIDType = 0x0d06 // MDS Create Notification (Used)
)

// Command Types
const (
	CMD_EVENT_REPORT           CMDType = 0x0000 // Event Report
	CMD_CONFIRMED_EVENT_REPORT CMDType = 0x0001 // Confirmed Event Report
	CMD_GET                    CMDType = 0x0003 // Get
	CMD_SET                    CMDType = 0x0004 // Set
	CMD_CONFIRMED_SET          CMDType = 0x0005 // Confirmed Set
	CMD_CONFIRMED_ACTION       CMDType = 0x0007 // Confirmed Action
)

// Association Control
const (
	MDDL_VERSION1 uint32 = 0x80000000 // Data Export Protocol Version (Used)
	NOMEN_VERSION uint32 = 0x40000000 // Nomenclature Version (Used)

	SYST_CLIENT uint32 = 0x80000000 // System Type Client (Used)
	SYST_SERVER uint32 = 0x00800000 // System Type Server (Used)

	HOT_START  uint32 = 0x80000000 // Hot Start (Used)
	WARM_START uint32 = 0x40000000 // Warm Start (Used)
	COLD_START uint32 = 0x20000000 // Cold Start (Used)

)
