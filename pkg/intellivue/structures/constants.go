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
	NOM_POLL_PROFILE_SUPPORT  OIDType = 0x0001
	NOM_MDIB_OBJ_SUPPORT      OIDType = 0x0102
	NOM_ATTR_POLL_PROFILE_EXT OIDType = 0xF001
)

// Event Types
const (
	NOM_NOTI_MDS_CREAT OIDType = 0x0d06 // MDS Create Notification
)

// Association Control
const (
	MDDL_VERSION1 uint32 = 0x80000000 // Data Export Protocol Version
	NOMEN_VERSION uint32 = 0x40000000 // Nomenclature Version

	SYST_CLIENT uint32 = 0x80000000 // System Type Client
	SYST_SERVER uint32 = 0x00800000 // System Type Server

	HOT_START  uint32 = 0x80000000 // Hot Start
	WARM_START uint32 = 0x40000000 // Warm Start
	COLD_START uint32 = 0x20000000 // Cold Start

)
