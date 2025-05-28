package intellivue

// Association Control
const (
	MDDL_VERSION1                    uint32 = 0x80000000  // Data Export Protocol Version
	NOMEN_VERSION                    uint32 = 0x40000000  // Nomenclature Version

	SYST_CLIENT                      uint32 = 0x80000000  // System Type Client
	SYST_SERVER                      uint32 = 0x00800000  // System Type Server

	HOT_START                        uint32 = 0x80000000  // Hot Start
	WARM_START                       uint32 = 0x40000000  // Warm Start
	COLD_START                       uint32 = 0x20000000  // Cold Start

	POLL_PROFILE_REV_0               PollProfileRevision = 0x80000000 // Poll Profile Revision
)

// Protocol Identification
const (
	NOM_POLL_PROFILE_SUPPORT  uint16 = 1      // id for polling profile
	NOM_MDIB_OBJ_SUPPORT      uint16 = 258    // supported objects for the active profile
	NOM_ATTR_POLL_PROFILE_EXT uint16 = 61441  //id for poll profile extensions opt. package
)