package packages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/constants"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
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
	Length   ASNLength
	UserData MDSEUserInfoStd
}

type MDSEUserInfoStd struct {
	ProtocolVersion     uint32
	NomenclatureVersion uint32
	FunctionalUnits     uint32
	SystemType          uint32
	StartupMode         uint32
	OptionList          *AttributeList
	SupportedAprofiles  *AttributeList
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

	supportedAprofiles := NewAttributeList()

	pollProfileExt := NewPollProfileSupport()
	mdibObjSupport := NewMDIBObjSupport()

	supportedAprofiles.Value = append(
		supportedAprofiles.Value,
		AVAType{
			AttributeID: NOM_POLL_PROFILE_SUPPORT,
			Value:       pollProfileExt,
		},
	)
	supportedAprofiles.Value = append(
		supportedAprofiles.Value,
		AVAType{
			AttributeID: NOM_MDIB_OBJ_SUPPORT,
			Value:       mdibObjSupport,
		},
	)
	supportedAprofiles.Count = uint16(len(supportedAprofiles.Value))

	userData := &MDSEUserInfoStd{
		ProtocolVersion:     MDDL_VERSION1,
		NomenclatureVersion: NOMEN_VERSION,
		FunctionalUnits:     0,
		SystemType:          SYST_CLIENT,
		StartupMode:         HOT_START,
		OptionList:          NewAttributeList(),
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
			Length:   ASNLength(len(userDataBytes)),
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
	liFieldTotalLength, err := LIField(totalLength).MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга LIField для общей длины: %w", err)
	}
	buf.Write(liFieldTotalLength)

	buf.Write(m.SessionData.Data)

	buf.WriteByte(m.PresentationHeader.Prefix)

	// writeLIField(&buf, LIField(pdataAPDULen)) // Length - Заменяем на MarshalBinary
	liFieldPdataLen, err := LIField(presentationHeaderLen).MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга LIField для длины P-DATA: %w", err)
	}
	buf.Write(liFieldPdataLen)

	buf.Write(m.PresentationHeader.Data)

	buf.Write(userData)

	buf.Write(m.PresentationTrailer)

	return buf.Bytes(), nil
}

func hex2bytes(hexStr string) ([]byte, error) {
	hexStr = bytes.NewBuffer(bytes.ReplaceAll(
		bytes.ReplaceAll([]byte(hexStr), []byte(" "), []byte("")),
		[]byte("\n"), []byte(""),
	)).String()

	var result []byte
	for i := 0; i < len(hexStr); i += 2 {
		b, err := hex2byte(hexStr[i : i+2])
		if err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

func hex2byte(hexStr string) (byte, error) {
	var result byte
	_, err := fmt.Sscanf(hexStr, "%02x", &result)
	return result, err
}

func (m *AssocReqMessage) ShowInfo() error {
	fmt.Printf("\n=== User Data Details ===\n")
	fmt.Printf("Protocol Version: 0x%08x\n", m.UserData.UserData.ProtocolVersion)
	fmt.Printf("Nomenclature Version: 0x%08x\n", m.UserData.UserData.NomenclatureVersion)
	fmt.Printf("Functional Units: 0x%08x\n", m.UserData.UserData.FunctionalUnits)
	fmt.Printf("System Type: 0x%08x\n", m.UserData.UserData.SystemType)
	fmt.Printf("Startup Mode: 0x%08x\n", m.UserData.UserData.StartupMode)

	fmt.Printf("\n=== Optional Packages Details ===\n")
	fmt.Printf("Count: %d\n", m.UserData.UserData.SupportedAprofiles.Count)
	fmt.Printf("Length: %d\n", m.UserData.UserData.SupportedAprofiles.Size())
	for i, ava := range m.UserData.UserData.SupportedAprofiles.Value {
		fmt.Printf("\nAVA Type #%d:\n", i+1)
		fmt.Printf("  Attribute ID: 0x%04x\n", ava.AttributeID)
		fmt.Printf("  Length: %d\n", ava.Size())
		if pollProfile, ok := ava.Value.(*PollProfileSupport); ok {
			fmt.Printf("  Poll Profile Revision: 0x%08x\n", pollProfile.PollProfileRevision)
			fmt.Printf("  Min Poll Period: %d\n", pollProfile.MinPollPeriod)
			fmt.Printf("  Max MTU RX: %d\n", pollProfile.MaxMtuRx)
			fmt.Printf("  Max MTU TX: %d\n", pollProfile.MaxMtuTx)
			fmt.Printf("  Max BW TX: 0x%08x\n", pollProfile.MaxBwTx)
			fmt.Printf("  Options: 0x%08x\n", pollProfile.Options)
		}
	}

	return nil
}

func (m MDSEUserInfoStd) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, m.ProtocolVersion); err != nil {
		return nil, fmt.Errorf("failed to write ProtocolVersion: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.NomenclatureVersion); err != nil {
		return nil, fmt.Errorf("failed to write NomenclatureVersion: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.FunctionalUnits); err != nil {
		return nil, fmt.Errorf("failed to write FunctionalUnits: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.SystemType); err != nil {
		return nil, fmt.Errorf("failed to write SystemType: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, m.StartupMode); err != nil {
		return nil, fmt.Errorf("failed to write StartupMode: %v", err)
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

func (m *AssocReqMessage) serializeUserData() ([]byte, error) {
	return m.UserData.MarshalBinary()
}

// AAREMessage представляет структуру сообщения AARE.
type AAREMessage struct {
	// TODO: Реализовать поля AAREMessage
}

// Size возвращает общую длину AAREMessage в байтах.
func (a *AAREMessage) Size() uint16 {
	// TODO: Реализовать расчет реального размера
	return 0
}

// MarshalBinary кодирует структуру AAREMessage в бинарный формат.
func (a *AAREMessage) MarshalBinary() ([]byte, error) {
	// TODO: Реализовать маршалинг
	return []byte{}, nil
}

// AARQMessage представляет структуру сообщения AARQ.
type AARQMessage struct {
	// TODO: Реализовать поля AARQMessage
}

// Size возвращает общую длину AARQMessage в байтах.
func (a *AARQMessage) Size() uint16 {
	// TODO: Реализовать расчет реального размера
	return 0
}

// MarshalBinary кодирует структуру AARQMessage в бинарный формат.
func (a *AARQMessage) MarshalBinary() ([]byte, error) {
	// TODO: Реализовать маршалинг
	return []byte{}, nil
}

// RLRQMessage представляет структуру сообщения RLRQ.
type RLRQMessage struct {
	// TODO: Реализовать поля RLRQMessage
}

// Size возвращает общую длину RLRQMessage в байтах.
func (r *RLRQMessage) Size() uint16 {
	// TODO: Реализовать расчет реального размера
	return 0
}

// MarshalBinary кодирует структуру RLRQMessage в бинарный формат.
func (r *RLRQMessage) MarshalBinary() ([]byte, error) {
	// TODO: Реализовать маршалинг
	return []byte{}, nil
}

// RLREMessage представляет структуру сообщения RLRE.
type RLREMessage struct {
	// TODO: Реализовать поля RLREMessage
}

// Size возвращает общую длину RLREMessage в байтах.
func (r *RLREMessage) Size() uint16 {
	// TODO: Реализовать расчет реального размера
	return 0
}

// MarshalBinary кодирует структуру RLREMessage в бинарный формат.
func (r *RLREMessage) MarshalBinary() ([]byte, error) {
	// TODO: Реализовать маршалинг
	return []byte{}, nil
}

// ABSEMessage представляет структуру сообщения ABSE.
type ABSEMessage struct {
	// TODO: Реализовать поля ABSEMessage
}

// Size возвращает общую длину ABSEMessage в байтах.
func (a *ABSEMessage) Size() uint16 {
	// TODO: Реализовать расчет реального размера
	return 0
}

// MarshalBinary кодирует структуру ABSEMessage в бинарный формат.
func (a *ABSEMessage) MarshalBinary() ([]byte, error) {
	// TODO: Реализовать маршалинг
	return []byte{}, nil
}

// Ассоциация - это логическое соединение между двумя приложениями.
type Association struct {
	RemoteAddr *net.TCPAddr
	conn       *net.TCPConn
	id         uint16
}

// NewAssociation создает новую ассоциацию с удаленным адресом.
func NewAssociation(remoteAddr *net.TCPAddr) *Association {
	// Генерация уникального ID для ассоциации (пример)
	id := uint16(time.Now().UnixNano()%65535) + 1
	return &Association{
		RemoteAddr: remoteAddr,
		id:         id,
	}
}

// Connect устанавливает TCP-соединение с удаленным узлом.
func (a *Association) Connect() error {
	conn, err := net.DialTCP("tcp", nil, a.RemoteAddr)
	if err != nil {
		return fmt.Errorf("не удалось установить соединение: %w", err)
	}
	a.conn = conn
	return nil
}

// SendAssocRequest отправляет запрос ассоциации.
func (a *Association) SendAssocRequest() error {
	assocReqMsg := NewAssocReqMessage()
	assocReqMsgData, err := assocReqMsg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("ошибка маршалинга AssocReqMessage: %w", err)
	}

	// Отправляем сообщение
	_, err = a.conn.Write(assocReqMsgData)
	if err != nil {
		return fmt.Errorf("не удалось отправить запрос ассоциации: %w", err)
	}

	return nil
}

// ReadResponse читает ответ от удаленного узла.
func (a *Association) ReadResponse() ([]byte, error) {
	// TODO: Реализовать чтение ответа с учетом структуры сообщений
	// Это потребует парсинга заголовков и содержимого APDU.
	// Пока просто читаем некоторое количество байт.
	buffer := make([]byte, 1024) // Примерный размер буфера
	n, err := a.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	return buffer[:n], nil
}

// Close закрывает соединение ассоциации.
func (a *Association) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
