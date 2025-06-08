package utils

import (
	"fmt"
	"strings"
	"sync"
)

// PrintHexDump выводит шестнадцатеричный дамп байтов с форматированием.
// Принимает указатель на sync.Mutex для синхронизации вывода в терминал.
func PrintHexDump(mu *sync.Mutex, title string, data []byte) {
	fmt.Printf("\n=== %s ===\n", title)
	fmt.Printf("Length: %d bytes\n", len(data))
	fmt.Println("Hex dump:")
	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		for j := 0; j < 16; j++ {
			if j+i < len(data) {
				fmt.Printf("%02x ", chunk[j])
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func PrintHex(data []byte) string {
	bytesPerLine := 16
	var builder strings.Builder

	for i := 0; i < len(data); i += bytesPerLine {
		chunk := data[i:min(i+bytesPerLine, len(data))]
		builder.WriteString(fmt.Sprintf("% x\n", chunk)) // пробелы между байтами
	}

	return strings.TrimSuffix(builder.String(), "\n") // убираем последний \n
}
