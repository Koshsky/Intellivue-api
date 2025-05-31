package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type NomPartition uint16
type OIDType uint16

const (
	// Partition IDs
	NOM_PART_OBJ         NomPartition = 0x0001 // Object oriented element, device nomenclature
	NOM_PART_SCADA       NomPartition = 0x0002 // Types of measurement and place of the measurement
	NOM_PART_EVT         NomPartition = 0x0003 // Codes for alerts
	NOM_PART_DIM         NomPartition = 0x0004 //  Units of measurement
	NOM_PART_PGRP        NomPartition = 0x0006 // Identification of parameter groups
	NOM_PART_INFRASTRUCT NomPartition = 0x0008 // Infrastructure for Data Export applications
	NOM_PART_EMFC        NomPartition = 0x0401 // EMFC
	NOM_PART_SETTINGS    NomPartition = 0x0402 // Settings

	// Managed Object Classes (MOC)
	NOM_MOC_VMO_METRIC_NU    OIDType = 0x0402 // numerics
	NOM_MOC_VMO_METRIC_SA_RT OIDType = 0x0402 // waves
	NOM_MOC_VMO_AL_MON       OIDType = 0x0402 // alerts
	NOM_MOC_PT_DEMOG         OIDType = 0x0402 // Pat. Demog
	NOM_MOC_VMS_MDS          OIDType = 0x0402 // MDS
)

// Whenever it is not clear from the context,
// from which nomenclature value range the OIDType comes,
// the TYPE data type is used. Here,
// the nomenclature value range (the partition) is explicitly identified.
type TYPE struct {
	Partition NomPartition
	Code      OIDType
}

func (t *TYPE) Size() uint16 {
	return 4
}

func (t *TYPE) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, t.Partition); err != nil {
		return nil, fmt.Errorf("failed to write Partition: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, t.Code); err != nil {
		return nil, fmt.Errorf("failed to write Code: %w", err)
	}

	return buf.Bytes(), nil
}

func (t *TYPE) UnmarshalBinary(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &t.Partition); err != nil {
		return fmt.Errorf("failed to read Partition: %w", err)
	}

	if err := binary.Read(r, binary.BigEndian, &t.Code); err != nil {
		return fmt.Errorf("failed to read Code: %w", err)
	}

	return nil
}
