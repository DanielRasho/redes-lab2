package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	SOCKET_PORT = ":8081"
	CRC32_POLY  = "100000100110000010001110110110111"
)
const (
	VITERBI_ALG = iota
	CRC32_ALG
)

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Should provide algorithm (viterbi/crc32)")
		return
	}

	algorithm := 0

	if args[1] == "viterbi" {
		algorithm = VITERBI_ALG
	} else if args[1] == "crc32" {
		algorithm = CRC32_ALG
	} else {
		fmt.Println("Invalid algorithm (viterbi/crc32)")
		return
	}

	// Listen on port 8080
	listener, err := net.Listen("tcp", SOCKET_PORT)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Go server listening on port 8080...")
	fmt.Println("Waiting for Python client connections...")
	fmt.Println("===========================================")

	// Accept connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn, algorithm)
	}
}

func handleConnection(conn net.Conn, algorithm int) {
	defer conn.Close()

	fmt.Printf("New connection from: %s\n", conn.RemoteAddr())

	// Receive the message
	message, err := receiveMessage(conn)
	if err != nil {
		fmt.Printf("Error receiving message: %v\n", err)
		return
	}

	fmt.Printf("Message received: %s\n", message)
	response := decodeWithAlgorithm(message, algorithm)

	// Optionally send a response back to Python client
	err = sendMessage(conn, response)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}

	fmt.Printf("\nConnection closed: %s\n\n", conn.RemoteAddr())
}

func decodeWithAlgorithm(input string, algorithm int) string {
	response := ""

	switch algorithm {
	case VITERBI_ALG:
		if err := checkViterbi(input); err != nil {
			// err.Error() ya es "sequence corrupted: detected N bit errors at positions [..]"
			response = fmt.Sprintf("❌ %s", err.Error())
			fmt.Println(err)
			fmt.Println("\t" + response)
		} else {
			response = "✅ No error detected in the message."
			fmt.Println("\t" + response)
		}
	case CRC32_ALG:
		if checkCRC32(input) {
			response = "✅ No error detected in the message."
			fmt.Println("\t" + response)
			decoded, err := binaryToText(input[:len(input)-32])
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Printf("\t Decoded msg: %s\n", decoded)
		} else {
			response = "❌ Error detected in the message."
			fmt.Println("\t" + response)
		}
	}

	return response
}

func receiveMessage(conn net.Conn) (string, error) {
	// First, read the 4-byte length prefix (big-endian unsigned int)
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read length prefix: %v", err)
	}

	// Convert bytes to uint32 (big-endian, thats the way python encode them)
	messageLength := binary.BigEndian.Uint32(lengthBytes)
	fmt.Printf("Expecting message of %d bytes\n", messageLength)

	// Now read the actual message
	messageBytes := make([]byte, messageLength)
	_, err = io.ReadFull(conn, messageBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read message: %v", err)
	}

	// Convert bytes to string
	message := string(messageBytes)

	return message, nil
}

func sendMessage(conn net.Conn, message string) error {
	messageBytes := []byte(message)

	// Create 4-byte big-endian length prefix
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(messageBytes)))

	// Send length prefix + message
	_, err := conn.Write(append(lengthBytes, messageBytes...))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func binaryToText(binaryStr string) (string, error) {
	// Remove any spaces or whitespace
	binaryStr = strings.ReplaceAll(binaryStr, " ", "")
	binaryStr = strings.TrimSpace(binaryStr)

	// Check if length is multiple of 8
	if len(binaryStr)%8 != 0 {
		return "", fmt.Errorf("binary string length (%d) is not a multiple of 8", len(binaryStr))
	}

	var result strings.Builder

	// Process each 8-bit chunk
	for i := 0; i < len(binaryStr); i += 8 {
		// Extract 8-bit chunk
		binaryByte := binaryStr[i : i+8]

		// Convert binary string to integer
		asciiValue, err := strconv.ParseInt(binaryByte, 2, 32)
		if err != nil {
			return "", fmt.Errorf("invalid binary sequence '%s' at position %d: %v", binaryByte, i, err)
		}

		// Check if it's a valid ASCII value
		if asciiValue < 0 || asciiValue > 127 {
			return "", fmt.Errorf("invalid ASCII value %d at position %d", asciiValue, i/8)
		}

		// Convert to character and append
		result.WriteByte(byte(asciiValue))
	}

	return result.String(), nil
}
