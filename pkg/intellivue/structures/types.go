package structures

type AbsoluteTime uint64

// 				OR
// type AbsoluteTime struct {
// 	Century uint8
// 	Year uint8
// 	Month uint8
// 	Day uint8
// 	Hour uint8
// 	Minute uint8
// 	Second uint8
// 	SecFractions uint8
// }
// СТРУКТУРА ТРЕБУЕТ РЕАЛИЗАЦИИ SIZE() MARSHAL() UNMARSHAL()

// TODO: КОГДА ВСЕ СТРУКТУРЫ БУДУТ РЕАЛИЗОВАНЫ
// ЕСЛИ:
//
//	ВСЕ ТИПЫ ВСТРЕЧАЮТСЯ ТОЛЬКО В ОПРЕДЕЛЕННЫХ ФАЙЛАХ
//
// ТО:
//
//	ПЕРЕМЕСТИТЬ ТИПЫ В СООТВЕТСТВУЮЩИЕ ФАЙЛЫ
//
// ИНАЧЕ:
//
//	РАЗДЕЛИТЬ ПУСТЫМИ СТРОКАМИ ТИПА ПО ГРУППАМ
type RelativeTime uint32
type PrivateOID uint16
type Handle uint16
type MdsContext uint16
type TextID uint32
type FLOATType uint32
type ErrorStatus uint16
type CMDType uint16
type ModifyOperator uint16
type Nomenclature uint16
type ConnectIndInfo = AttributeList
type ProtocolVersion uint32
type NomenclatureVersion uint32
type FunctionalUnits uint32
type SystemType uint32
type StartupMode uint32
type MeasurementState uint16
type SimpleColour uint16
type MetricSpec uint16
type MetricAccess uint16
type MetricRelevance uint16
type MetricModality uint16
type SaFlags uint16
type SaFixedValId uint16
type MetricState uint16
type MeasureMode uint16
type Language uint16
type StringFormat uint16
type MDSStatus uint16
type ApplicationArea uint16
type LineFrequency uint16
type AlertState uint16
type AlertPriority uint16
type AlertFlags uint16
type PatDmgState uint16
type PatientType uint16
type PatPacedMode uint16
type PatientSex uint16
type PtBsaFormula uint16
type ApplProtoId uint16
type TransProtoId uint16
type ProtoOptions uint16

// МОЖНО В МЕТОДАХ MARSHAL() UNMARSHAL() КАСТОВАТЬ ЭТИ ТИПЫ В ИХ АЛИАСЫ.
