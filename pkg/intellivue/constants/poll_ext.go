package constants

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
