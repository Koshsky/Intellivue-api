package packages

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures/attributes"
)

type AssociationRequest struct {
	AssocReqSessionHeader       SessionHeader
	AssocReqSessionData         SessionData
	AssocReqPresentationHeader  PresentationHeader
	AssocReqUserData            UserData
	AssocReqPresentationTrailer []byte
}

func NewAssociationRequest() *AssociationRequest {
	sessionData := []byte{0x05, 0x08, 0x13, 0x01, 0x00, 0x16, 0x01, 0x02, 0x80, 0x00, 0x14, 0x02, 0x00, 0x02}
	presentationData := []byte{
		0x31, 0x80, 0xA0, 0x80, 0x80, 0x01, 0x01, 0x00, 0x00, 0xA2, 0x80, 0xA0, 0x03, 0x00, 0x00, 0x01,
		0xA4, 0x80, 0x30, 0x80, 0x02, 0x01, 0x01, 0x06, 0x04, 0x52, 0x01, 0x00, 0x01, 0x30, 0x80, 0x06,
		0x02, 0x51, 0x01, 0x00, 0x00, 0x00, 0x00, 0x30, 0x80, 0x02, 0x01, 0x02, 0x06, 0x0C, 0x2A, 0x86,
		0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x30, 0x80, 0x06, 0x0C, 0x2A, 0x86,
		0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x61, 0x80, 0x30, 0x80, 0x02, 0x01, 0x01, 0xA0, 0x80, 0x60, 0x80, 0xA1, 0x80, 0x06, 0x0C, 0x2A,
		0x86, 0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0xBE, 0x80, 0x28,
		0x80, 0x06, 0x0C, 0x2A, 0x86, 0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x02, 0x01, 0x02, 0x81}
	trailerData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	supportedAprofiles := &attributes.AttributeList{}

	pollProfileExt := attributes.NewPollProfileExt()
	optionalPackages := &attributes.AttributeList{}
	optionalPackages.Append(pollProfileExt.AttributeID, pollProfileExt.Value)
	pollProfileSupport := &attributes.PollProfileSupport{
		PollProfileRevision: base.POLL_PROFILE_REV_0,
		MinPollPeriod:       0x00000fa0, // 4000 ms (4 seconds)
		MaxMtuRx:            0x000005b0, // 1456 bytes
		MaxMtuTx:            0x000005b0, // 1456 bytes
		MaxBwTx:             0xFFFFFFFF, // Unlimited bandwidth
		Options: base.P_OPT_DYN_CREATE_OBJECTS |
			base.P_OPT_DYN_DELETE_OBJECTS,
		OptionalPackages: optionalPackages,
	}
	supportedAprofiles.Append(base.NOM_POLL_PROFILE_SUPPORT, pollProfileSupport)

	mdibObjSupport := attributes.NewMDIBObjSupport()
	supportedAprofiles.Append(base.NOM_MDIB_OBJ_SUPPORT, mdibObjSupport)

	userData := structures.MDSEUserInfoStd{
		ProtocolVersion:     base.MDDL_VERSION1,
		NomenclatureVersion: base.NOMEN_VERSION,
		FunctionalUnits:     0,
		SystemType:          base.SYST_CLIENT,
		StartupMode:         base.HOT_START,
		OptionList:          &attributes.AttributeList{},
		SupportedAprofiles:  supportedAprofiles,
	}

	req := &AssociationRequest{
		AssocReqSessionHeader: SessionHeader{
			Type: 0x0D,
			// need to specify Length in MarshalBinary
		},
		AssocReqSessionData: SessionData{
			Data: sessionData,
		},
		AssocReqPresentationHeader: PresentationHeader{
			Type: 0xC1,
			// need to specify Length in MarshalBinary
			Data: presentationData,
		},
		AssocReqUserData: UserData{
			Data: userData,
		},
		AssocReqPresentationTrailer: trailerData,
	}

	req.AssocReqPresentationHeader.Length = base.LIField(
		uint16(len(req.AssocReqPresentationTrailer)) +
			req.AssocReqUserData.Size() +
			uint16(len(req.AssocReqPresentationHeader.Data)))
	req.AssocReqSessionHeader.Length = base.LIField(
		uint16(len(req.AssocReqPresentationTrailer)) +
			req.AssocReqUserData.Size() +
			req.AssocReqPresentationHeader.Size() +
			uint16(len(req.AssocReqSessionData.Data)))

	return req
}

func (m *AssociationRequest) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	sessionHeader, err := m.AssocReqSessionHeader.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AssocReqSessionHeader: %v", err)
	}
	buf.Write(sessionHeader)
	sessionData, err := m.AssocReqSessionData.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AssocReqSessionData: %v", err)
	}
	buf.Write(sessionData)
	presentationHeader, err := m.AssocReqPresentationHeader.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AssocReqPresentationHeader: %v", err)
	}
	buf.Write(presentationHeader)
	userData, err := m.AssocReqUserData.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AssocReqUserData: %v", err)
	}
	buf.Write(userData)
	buf.Write(m.AssocReqPresentationTrailer)

	return buf.Bytes(), nil
}

func (m *AssociationRequest) ShowInfo() {
	log.Println("--- AssociationRequest ---")
	m.AssocReqSessionHeader.ShowInfo(1)
	m.AssocReqSessionData.ShowInfo(1)
	m.AssocReqPresentationHeader.ShowInfo(1)
	m.AssocReqUserData.ShowInfo(1)
	log.Printf("  AssocReqPresentationTrailer: %v", m.AssocReqPresentationTrailer)
}
