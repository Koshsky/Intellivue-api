package packages

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/structures/attributes"
)

type SessionHeader struct {
	Type   byte
	Length uint16
}

type AssocReqSessionData struct {
	Data []byte
}

type AssocReqPresentationHeader struct {
	Prefix byte
	Length uint16
	Data   []byte
}

type AssocReqUserData struct {
	Length   base.ASNLength
	UserData MDSEUserInfoStd
}

type MDSEUserInfoStd struct {
	ProtocolVersion     uint32
	NomenclatureVersion uint32
	FunctionalUnits     uint32
	SystemType          uint32
	StartupMode         uint32
	OptionList          *attributes.AttributeList
	SupportedAprofiles  *attributes.AttributeList
}

type AssocReqMessage struct {
	SessionHeader       SessionHeader
	SessionData         AssocReqSessionData
	PresentationHeader  AssocReqPresentationHeader
	UserData            AssocReqUserData
	PresentationTrailer []byte
}

func NewAssocReqMessage() *AssocReqMessage {
	sessionData := []byte{0x05, 0x08, 0x13, 0x01, 0x00, 0x16, 0x01, 0x02, 0x80, 0x00, 0x14, 0x02, 0x00, 0x02}
	presentationData := []byte{
		0x31, 0x80, 0xA0, 0x80, 0x80, 0x01, 0x01, 0x00, 0x00, 0xA2, 0x80, 0xA0, 0x03, 0x00, 0x00, 0x01,
		0xA4, 0x80, 0x30, 0x80, 0x02, 0x01, 0x01, 0x06, 0x04, 0x52, 0x01, 0x00, 0x01, 0x30, 0x80, 0x06,
		0x02, 0x51, 0x01, 0x00, 0x00, 0x00, 0x00, 0x30, 0x80, 0x02, 0x01, 0x02, 0x06, 0x0C, 0x2A, 0x86,
		0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x30, 0x80, 0x06, 0x0C, 0x2A, 0x86,
		0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x61, 0x80, 0x30, 0x80, 0x02, 0x01, 0x01, 0xA0, 0x80, 0x60, 0x80, 0xA1, 0x80, 0x06, 0x0C, 0x2A,
		0x86, 0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0xBE, 0x80, 0x28,
		0x80, 0x06, 0x0C, 0x2A, 0x86, 0x48, 0xCE, 0x14, 0x02, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x02, 0x01, 0x02, 0x81,
	}
	trailerData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	supportedAprofiles := &attributes.AttributeList{}

	pollProfileExt := &attributes.PollProfileSupport{
		PollProfileRevision: base.POLL_PROFILE_REV_0,
		MinPollPeriod:       0x00000fa0, // 4000 ms (4 seconds)
		MaxMtuRx:            0x000005b0, // 1456 bytes
		MaxMtuTx:            0x000005b0, // 1456 bytes
		MaxBwTx:             0xFFFFFFFF, // Unlimited bandwidth
		Options: base.P_OPT_DYN_CREATE_OBJECTS |
			base.P_OPT_DYN_DELETE_OBJECTS,
		OptionalPackages: &attributes.AttributeList{
			Value: []attributes.AVAType{
				attributes.NewPollProfileExtensionAVAType(),
			},
		},
	}

	mdibObjSupport := attributes.NewMDIBObjSupport()

	dataPollProfileExt, _ := pollProfileExt.MarshalBinary()
	hbPollProfileExt := attributes.HexBytes(dataPollProfileExt)
	supportedAprofiles.Value = append(
		supportedAprofiles.Value,
		attributes.AVAType{
			AttributeID: base.NOM_POLL_PROFILE_SUPPORT,
			Value:       &hbPollProfileExt,
		},
	)
	dataMdibObjSupport, _ := mdibObjSupport.MarshalBinary()
	hbMdibObjSupport := attributes.HexBytes(dataMdibObjSupport)
	supportedAprofiles.Value = append(
		supportedAprofiles.Value,
		attributes.AVAType{
			AttributeID: base.NOM_MDIB_OBJ_SUPPORT,
			Value:       &hbMdibObjSupport,
		},
	)

	userData := &MDSEUserInfoStd{
		ProtocolVersion:     base.MDDL_VERSION1,
		NomenclatureVersion: base.NOMEN_VERSION,
		FunctionalUnits:     0,
		SystemType:          base.SYST_CLIENT,
		StartupMode:         base.HOT_START,
		OptionList:          &attributes.AttributeList{},
		SupportedAprofiles:  supportedAprofiles,
	}

	userDataBytes, err := userData.MarshalBinary()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &AssocReqMessage{
		SessionHeader: SessionHeader{
			Type: 0x0D,
		},
		SessionData: AssocReqSessionData{
			Data: sessionData,
		},
		PresentationHeader: AssocReqPresentationHeader{
			Prefix: 0xC1,
			Data:   presentationData,
		},
		UserData: AssocReqUserData{
			Length:   base.ASNLength(len(userDataBytes)),
			UserData: *userData,
		},
		PresentationTrailer: trailerData,
	}
}

func (m *AssocReqMessage) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	userData, err := m.UserData.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize user data: %v", err)
	}

	presentationHeaderLen := len(m.PresentationHeader.Data) + len(userData) + 16
	sessionDataLen := len(m.SessionData.Data)
	totalLength := sessionDataLen + presentationHeaderLen

	if presentationHeaderLen > 255 {
		totalLength += 4 // 3 bytes for length field
	} else {
		totalLength += 2 // 1 byte for length field
	}

	buf.WriteByte(m.SessionHeader.Type)

	// writeLIField(&buf, LIField(totalLength)) // Length - Заменяем на MarshalBinary
	liFieldTotalLength, err := base.LIField(totalLength).MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LIField для общей длины: %w", err)
	}
	buf.Write(liFieldTotalLength)

	buf.Write(m.SessionData.Data)

	buf.WriteByte(m.PresentationHeader.Prefix)

	// writeLIField(&buf, LIField(pdataAPDULen)) // Length - Заменяем на MarshalBinary
	liFieldPdataLen, err := base.LIField(presentationHeaderLen).MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LIField для длины P-DATA: %w", err)
	}
	buf.Write(liFieldPdataLen)

	buf.Write(m.PresentationHeader.Data)

	buf.Write(userData)

	buf.Write(m.PresentationTrailer)

	return buf.Bytes(), nil
}

func (m MDSEUserInfoStd) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, m.ProtocolVersion); err != nil {
		return nil, fmt.Errorf("failed to marshal ProtocolVersion: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.NomenclatureVersion); err != nil {
		return nil, fmt.Errorf("failed to marshal NomenclatureVersion: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.FunctionalUnits); err != nil {
		return nil, fmt.Errorf("failed to marshal FunctionalUnits: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.SystemType); err != nil {
		return nil, fmt.Errorf("failed to marshal SystemType: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.StartupMode); err != nil {
		return nil, fmt.Errorf("failed to marshal StartupMode: %v", err)
	}

	optionListData, err := m.OptionList.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OptionList: %v", err)
	}
	buf.Write(optionListData)

	profileData, err := m.SupportedAprofiles.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SupportedAprofiles: %v", err)
	}
	buf.Write(profileData)

	return buf.Bytes(), nil
}

func (u *AssocReqUserData) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Сериализуем Length как ASN.1 длину
	lengthBytes, err := u.Length.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Length: %v", err)
	}
	buf.Write(lengthBytes)

	// Сериализуем UserData
	userData, err := u.UserData.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UserData: %v", err)
	}
	buf.Write(userData)

	return buf.Bytes(), nil
}
