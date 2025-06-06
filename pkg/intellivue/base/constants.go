package base

const (
	NOM_ATTR_AL_MON_P_AL_LIST OIDType = 0x0902
	NOM_ATTR_AL_MON_T_AL_LIST OIDType = 0x0904
	NOM_ATTR_DEV_AL_COND      OIDType = 0x0916

	NOM_ATTR_ID_HANDLE       OIDType = 0x0921
	NOM_ATTR_ID_TYPE         OIDType = 0x092f
	NOM_ATTR_METRIC_SPECN    OIDType = 0x093f
	NOM_ATTR_ID_LABEL        OIDType = 0x0924
	NOM_ATTR_ID_LABEL_STRING OIDType = 0x0927
	NOM_ATTR_COLOR           OIDType = 0x0911
	NOM_ATTR_DISP_RES        OIDType = 0x0917
	NOM_ATTR_NU_VAL_OBS      OIDType = 0x0950
	NOM_ATTR_TIME_STAMP_ABS  OIDType = 0x0990
	NOM_ATTR_NU_CMPD_VAL_OBS OIDType = 0x094B
)

// APDU Types
const (
	ROIV_APDU  uint16 = 0x0001 // RO Initiate and Vote APDU
	RORS_APDU  uint16 = 0x0002 // RO Reply APDU
	ROER_APDU  uint16 = 0x0003 // RO Error APDU
	ROLRS_APDU uint16 = 0x0005 // RO Linked Reply APDU
)

// Protocol Identification
const (
	NOM_PART_OBJ         NomPartition = 0x0001
	NOM_PART_SCADA       NomPartition = 0x0002
	NOM_PART_EVT         NomPartition = 0x0003
	NOM_PART_DIM         NomPartition = 0x0004
	NOM_PART_PGRP        NomPartition = 0x0006
	NOM_PART_INFRASTRUCT NomPartition = 0x0008
	NOM_PART_EMFC        NomPartition = 0x0401
	NOM_PART_SETTINGS    NomPartition = 0x0402

	NOM_MOC_VMO_METRIC_NU    OIDType = 0x0006
	NOM_MOC_VMO_METRIC_SA_RT OIDType = 0x0009
	NOM_MOC_VMO_AL_MON       OIDType = 0x0036
	NOM_MOC_PT_DEMOG         OIDType = 0x002a
	NOM_MOC_VMS_MDS          OIDType = 0x0021

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

const (
	CMD_EVENT_REPORT           CMDType = 0x0000
	CMD_CONFIRMED_EVENT_REPORT CMDType = 0x0001
	CMD_GET                    CMDType = 0x0003
	CMD_SET                    CMDType = 0x0004
	CMD_CONFIRMED_SET          CMDType = 0x0005
	CMD_CONFIRMED_ACTION       CMDType = 0x0007
)

const (
	NOM_ACT_POLL_MDIB_DATA     OIDType = 0x0c16
	NOM_ACT_POLL_MDIB_DATA_EXT OIDType = 0xf13b
)

const (
	POLL_PROFILE_REV_0 PollProfileRevision = 0x80000000

	P_OPT_DYN_CREATE_OBJECTS PollProfileOptions = 0x40000000
	P_OPT_DYN_DELETE_OBJECTS PollProfileOptions = 0x20000000
)
