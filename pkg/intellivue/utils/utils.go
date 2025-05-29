package utils

import (
	"fmt"
)

// PrintHexDump выводит бинарные данные в HEX-формате.
func PrintHexDump(title string, data []byte) {
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
}
