package intellivue

import (
	"bytes"
	"fmt"
	"testing"
)

func printHexDump(title string, data []byte) {
	fmt.Printf("\n=== %s ===\n", title)
	fmt.Printf("Length: %d bytes\n", len(data))
	fmt.Println("Hex dump:")
	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		// fmt.Printf("%04x: ", i)

		for j := 0; j < 16; j++ {
			if j+i < len(data) {
				fmt.Printf("%02x ", chunk[j])
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println()
	}
}

func TestAssocReqMessage(t *testing.T) {
	msg := NewAssocReqMessage()
	{
		userData, err := msg.serializeUserData()
		if err != nil {
			t.Fatalf("Failed to serialize user data: %v", err)
		}

		totalLength := len(msg.SessionData.Data) + len(msg.PresentationHeader.Data) +
			len(userData) + len(msg.PresentationTrailer)

		var sessionHeaderBuf bytes.Buffer
		sessionHeaderBuf.WriteByte(msg.SessionHeader.Type)
		writeLIField(&sessionHeaderBuf, LIField(totalLength))
		printHexDump("Session Header", sessionHeaderBuf.Bytes())

		printHexDump("Session Data", msg.SessionData.Data)

		var presentationHeaderBuf bytes.Buffer
		presentationHeaderBuf.WriteByte(msg.PresentationHeader.Prefix)
		writeLIField(&presentationHeaderBuf, LIField(len(userData)))
		presentationHeaderBuf.Write(msg.PresentationHeader.Data)
		printHexDump("Complete Presentation Header", presentationHeaderBuf.Bytes())

		printHexDump("User Data", userData)

		printHexDump("Presentation Trailer", msg.PresentationTrailer)
	}

	// СЕРИАЛИЗАЦИЯ ПОЛНОГО СООБЩЕНИЯ
	fullMessage, err := msg.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to serialize message: %v", err)
	}
	printHexDump("Complete Message", fullMessage)

	if err := msg.ShowInfo(); err != nil {
		t.Fatalf("Failed to show message info: %v", err)
	}
}
