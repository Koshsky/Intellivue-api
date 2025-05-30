package intellivue

import (
	"bytes"
	"testing"

	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/structures"
	. "github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

func TestAssocReqMessage(t *testing.T) {
	msg := NewAssocReqMessage()
	{
		userData, err := msg.UserData.MarshalBinary()
		if err != nil {
			t.Fatalf("Failed to marshal user data: %v", err)
		}

		totalLength := len(msg.SessionData.Data) + len(msg.PresentationHeader.Data) +
			len(userData) + len(msg.PresentationTrailer)

		var sessionHeaderBuf bytes.Buffer
		sessionHeaderBuf.WriteByte(msg.SessionHeader.Type)
		liFieldTotalLength, err := LIField(totalLength).MarshalBinary()
		if err != nil {
			t.Fatalf("Ошибка маршалинга LIField для общей длины: %v", err)
		}
		sessionHeaderBuf.Write(liFieldTotalLength)

		var presentationHeaderBuf bytes.Buffer
		presentationHeaderBuf.WriteByte(msg.PresentationHeader.Prefix)
		liFieldUserDataLen, err := LIField(len(userData)).MarshalBinary()
		if err != nil {
			t.Fatalf("Ошибка маршалинга LIField для длины пользовательских данных: %v", err)
		}
		presentationHeaderBuf.Write(liFieldUserDataLen)

		presentationHeaderBuf.Write(msg.PresentationHeader.Data)

	}

	fullMessage, err := msg.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	PrintHexDump("Complete Message", fullMessage)

}
