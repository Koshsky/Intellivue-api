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
