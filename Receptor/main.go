package main

import (
	"fmt"
)

func xor(a, b byte) byte {
	if a == b {
		return '0'
	}
	return '1'
}

func checkCRC32(msgWithCRC string, poly string) bool {
	msg := []byte(msgWithCRC)
	p := []byte(poly)
	msgLen := len(msg)
	polyLen := len(p)

	// Copy to avoid mutating the original
	working := make([]byte, msgLen)
	copy(working, msg)

	for i := 0; i <= msgLen-polyLen; i++ {
		if working[i] == '1' {
			for j := 0; j < polyLen; j++ {
				working[i+j] = xor(working[i+j], p[j])
			}
		}
	}

	// Check if remainder is all zeros
	for i := msgLen - (polyLen - 1); i < msgLen; i++ {
		if working[i] != '0' {
			return false
		}
	}
	return true
}

func main() {
	var messageWithCRC string
	fmt.Print("Enter binary message with CRC-32: ")
	_, err := fmt.Scanln(&messageWithCRC)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Estandar CRC-32 defined from IEEE
	poly := "100000100110000010001110110110111"

	if checkCRC32(messageWithCRC, poly) {
		fmt.Println("✅ No error detected in the message.")
	} else {
		fmt.Println("❌ Error detected in the message.")
	}
}
