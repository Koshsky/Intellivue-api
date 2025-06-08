package intellivue

import (
	"sync"
	"testing"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
	"github.com/Koshsky/Intellivue-api/pkg/intellivue/utils"
)

var testMutex sync.Mutex // Мьютекс для синхронизации вывода в тестах

func TestAssocReqMessage(t *testing.T) {
	msg := packages.NewAssociationRequest()
	msg.ShowInfo()

	fullMessage, err := msg.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	utils.PrintHexDump(&testMutex, "Complete Message", fullMessage)

}
