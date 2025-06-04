package intellivue

import (
	"bytes"
	"sync"
	"testing"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/base"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

var testMutex sync.Mutex // Мьютекс для синхронизации вывода в тестах

func TestAssocReqMessage(t *testing.T) {
	msg := packages.NewAssocReqMessage()
	{
		userData, err := msg.UserData.MarshalBinary()
		if err != nil {
			t.Fatalf("Failed to marshal user data: %v", err)
		}

		totalLength := len(msg.SessionData.Data) + len(msg.PresentationHeader.Data) +
			len(userData) + len(msg.PresentationTrailer)

		var sessionHeaderBuf bytes.Buffer
		sessionHeaderBuf.WriteByte(msg.SessionHeader.Type)
		liFieldTotalLength, err := base.LIField(totalLength).MarshalBinary()
		if err != nil {
			t.Fatalf("failed to marshal LIField для общей длины: %v", err)
		}
		sessionHeaderBuf.Write(liFieldTotalLength)

		var presentationHeaderBuf bytes.Buffer
		presentationHeaderBuf.WriteByte(msg.PresentationHeader.Prefix)
		liFieldUserDataLen, err := base.LIField(len(userData)).MarshalBinary()
		if err != nil {
			t.Fatalf("failed to marshal LIField для длины пользовательских данных: %v", err)
		}
		presentationHeaderBuf.Write(liFieldUserDataLen)

		presentationHeaderBuf.Write(msg.PresentationHeader.Data)

	}

	fullMessage, err := msg.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	utils.PrintHexDump(&testMutex, "Complete Message", fullMessage)

}
