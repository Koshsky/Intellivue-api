package constants

// Association Control
const (
	MDDL_VERSION1 uint32 = 0x80000000 // Data Export Protocol Version (Used)
	NOMEN_VERSION uint32 = 0x40000000 // Nomenclature Version (Used)

	SYST_CLIENT uint32 = 0x80000000 // System Type Client (Used)
	SYST_SERVER uint32 = 0x00800000 // System Type Server (Used)

	HOT_START  uint32 = 0x80000000 // Hot Start (Used)
	WARM_START uint32 = 0x40000000 // Warm Start (Used)
	COLD_START uint32 = 0x20000000 // Cold Start (Used)

	POLL_PROFILE_REV_0 uint32 = 0x80000000 // Poll Profile Revision (Used)
)

// Protocol Identification
const (
	NOM_POLL_PROFILE_SUPPORT  uint16 = 1     // id for polling profile (Used)
	NOM_MDIB_OBJ_SUPPORT      uint16 = 258   // supported objects for the active profile (Unused)
	NOM_ATTR_POLL_PROFILE_EXT uint16 = 61441 // id for poll profile extensions opt. package (Used)
)

// APDU Types
const (
	ROIV_APDU  uint16 = 1 // RO Initiate and Vote APDU (Used)
	RORS_APDU  uint16 = 2 // RO Reply APDU (Used)
	ROER_APDU  uint16 = 3 // RO Error APDU (Used)
	ROLRS_APDU uint16 = 5 // RO Linked Reply APDU (Unused)
)

// Command Types
const (
	CMD_CONFIRMED_ACTION       uint16 = 0x0007 // Confirmed Action (Used)
	CMD_CONFIRMED_EVENT_REPORT uint16 = 0x0001 // Confirmed Event Report (Used)
)

// Action Types
const (
	NOM_ACT_POLL_MDIB_DATA     uint16 = 3094  // Poll MDIB Data Action (Used)
	NOM_ACT_POLL_MDIB_DATA_EXT uint16 = 61755 // Poll MDIB Data Extended Action (Unused)
)

// Managed Object Classes (MOC)
const (
	NOM_MOC_VMS_MDS          uint16 = 0x0021 // Virtual Medical Systems MDS (Used)
	NOM_MOC_VMO_AL_MON       uint16 = 54     // alarms (Used)
	NOM_MOC_VMO_METRIC_NU    uint16 = 36     // numeric (Used)
	NOM_MOC_VMO_METRIC_SA_RT uint16 = 9      // waves (Used)
	NOM_MOC_PT_DEMOG         uint16 = 42     // Pat.Demog (Used)
)

// Poll Profile Extension Options
const (
	P_OPT_DYN_CREATE_OBJECTS uint32 = 0x40000000 // Dynamic Create Objects (Used)
	P_OPT_DYN_DELETE_OBJECTS uint32 = 0x20000000 // Dynamic Delete Objects (Used)

	POLL_EXT_PERIOD_NU_1SEC       uint32 = 0x80000000 // 1 sec Real-time Numerics (Used)
	POLL_EXT_PERIOD_NU_AVG_12SEC  uint32 = 0x40000000 // 12 sec averaged Numerics (Used)
	POLL_EXT_PERIOD_NU_AVG_60SEC  uint32 = 0x20000000 // 60 sec averaged Numerics (Used)
	POLL_EXT_PERIOD_NU_AVG_300SEC uint32 = 0x10000000 // 300 sec averaged Numerics (Used)
	POLL_EXT_PERIOD_RTSA          uint32 = 0x08000000 // Real-time Status and Alarms (Used)
	POLL_EXT_ENUM                 uint32 = 0x04000000 // allow enumeration objects (Used)
	POLL_EXT_NU_PRIO_LIST         uint32 = 0x02000000 // allow numeric priority list to be set (Used)
	POLL_EXT_DYN_MODALITIES       uint32 = 0x01000000 // send timestamps for numerics with dynamic modalities (Used)
)

// Event Types
const (
	NOM_NOTI_MDS_CREAT uint16 = 0x0d06 // MDS Create Notification (Used)
)

// Пример неиспользуемой константы (для демонстрации)
// const UNUSED_CONSTANT uint16 = 999
