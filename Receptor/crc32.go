package main

func xor(a, b byte) byte {
	if a == b {
		return '0'
	}
	return '1'
}

func checkCRC32(msgWithCRC string) bool {
	msg := []byte(msgWithCRC)
	p := []byte(CRC32_POLY)
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
