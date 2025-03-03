package address

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

// C32 characters used for encoding
const c32Chars = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// c32CharMap provides O(1) lookups for C32 character values
// It normalizes uppercase and lowercase letters, and handles special characters:
// O/o → 0, L/l → 1, I/i → 1
var c32CharMap [128]int

func init() {
	// Initialize c32CharMap with -1 (invalid)
	for i := 0; i < len(c32CharMap); i++ {
		c32CharMap[i] = -1
	}

	// Set values for valid C32 characters
	for i, c := range c32Chars {
		c32CharMap[c] = i
		// Also add lowercase mapping
		if c >= 'A' && c <= 'Z' {
			c32CharMap[c+32] = i // lowercase version
		}
	}

	// Special character mappings
	specialChars := [][2]rune{
		{'O', '0'},
		{'L', '1'},
		{'I', '1'},
	}

	for _, pair := range specialChars {
		val := c32CharMap[pair[1]]
		c32CharMap[pair[0]] = val
		// Also add lowercase mapping
		c32CharMap[pair[0]+32] = val
	}
}

// GetMaxC32EncodeOutputLen calculates the maximum C32 encoded output size given an input size.
// Each C32 character encodes 5 bits.
func GetMaxC32EncodeOutputLen(inputLen int) int {
	capacity := float64(inputLen+(inputLen%5)) / 5.0 * 8.0
	return int(capacity)
}

// EncodeC32 encodes a byte slice using C32 encoding.
func EncodeC32(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	capacity := GetMaxC32EncodeOutputLen(len(input))
	buffer := make([]byte, capacity)
	bytesWritten, _ := EncodeC32ToBuffer(input, buffer)
	return string(buffer[:bytesWritten])
}

// EncodeC32ToBuffer encodes input bytes into a C32 encoded output buffer.
// Returns the number of bytes written to the output buffer.
func EncodeC32ToBuffer(input []byte, output []byte) (int, error) {
	minLen := GetMaxC32EncodeOutputLen(len(input))
	if len(output) < minLen {
		return 0, fmt.Errorf("C32 encode output buffer is too small, given size %d, need minimum size %d",
			len(output), minLen)
	}

	var carry byte
	var carryBits byte
	position := 0

	// Process bytes in reverse order
	for i := len(input) - 1; i >= 0; i-- {
		currentValue := input[i]
		lowBitsToTake := 5 - carryBits
		lowBits := currentValue & ((1 << lowBitsToTake) - 1)
		c32Value := (lowBits << carryBits) + carry

		output[position] = c32Chars[c32Value]
		position++

		carryBits = (8 + carryBits) - 5
		carry = currentValue >> (8 - carryBits)

		if carryBits >= 5 {
			c32Value = carry & ((1 << 5) - 1)
			output[position] = c32Chars[c32Value]
			position++

			carryBits = carryBits - 5
			carry = carry >> 5
		}
	}

	if carryBits > 0 {
		output[position] = c32Chars[carry]
		position++
	}

	// Remove leading zeros from c32 encoding
	for position > 0 && output[position-1] == c32Chars[0] {
		position--
	}

	// Add leading zeros from input
	for _, currentValue := range input {
		if currentValue == 0 {
			output[position] = c32Chars[0]
			position++
		} else {
			break
		}
	}

	// Reverse the buffer
	for i, j := 0, position-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}

	return position, nil
}

// DecodeC32 decodes a C32 encoded string back to bytes.
func DecodeC32(input string) ([]byte, error) {
	// Must be ASCII
	for i := 0; i < len(input); i++ {
		if input[i] >= 128 {
			return nil, errors.New("invalid c32 string: must be ASCII")
		}
	}
	return DecodeC32Bytes([]byte(input))
}

// DecodeC32Bytes decodes a C32 encoded byte slice back to bytes.
func DecodeC32Bytes(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return []byte{}, nil
	}

	initialCapacity := len(input)
	result := make([]byte, 0, initialCapacity)
	var carry uint16
	var carryBits byte // Can be up to 5

	c32Digits := make([]byte, len(input))

	// Process in reverse order
	for i := len(input) - 1; i >= 0; i-- {
		if int(input[i]) >= len(c32CharMap) || c32CharMap[input[i]] == -1 {
			return nil, fmt.Errorf("invalid c32 character: %c", input[i])
		}
		c32Digits[len(input)-i-1] = byte(c32CharMap[input[i]])
	}

	for _, current5bit := range c32Digits {
		carry += uint16(current5bit) << carryBits
		carryBits += 5

		if carryBits >= 8 {
			result = append(result, byte(carry&0xFF))
			carryBits -= 8
			carry = carry >> 8
		}
	}

	if carryBits > 0 {
		result = append(result, byte(carry))
	}

	// Remove trailing zeros
	i := len(result)
	for i > 0 && result[i-1] == 0 {
		i--
	}
	result = result[:i]

	// Add leading zeros from input
	for i := len(c32Digits) - 1; i >= 0; i-- {
		if c32Digits[i] == 0 {
			result = append(result, 0)
		} else {
			break
		}
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// C32CheckEncodePrefixed encodes data with a version and checksum, prefixed by the given byte.
func C32CheckEncodePrefixed(version byte, data []byte, prefix byte) ([]byte, error) {
	if version >= 32 {
		return nil, fmt.Errorf("invalid version %d", version)
	}

	dataLen := len(data)
	buffer := make([]byte, dataLen+4)

	// Calculate double SHA256 checksum
	hash1 := sha256.Sum256(append([]byte{version}, data...))
	hash2 := sha256.Sum256(hash1[:])
	checksum := hash2[:4]

	// Copy data and checksum to buffer
	copy(buffer[:dataLen], data)
	copy(buffer[dataLen:], checksum)

	capacity := GetMaxC32EncodeOutputLen(len(buffer)) + 2
	result := make([]byte, capacity)

	result[0] = prefix
	result[1] = c32Chars[version]
	bytesWritten, err := EncodeC32ToBuffer(buffer, result[2:])
	if err != nil {
		return nil, err
	}

	return result[:bytesWritten+2], nil
}

// C32CheckDecode decodes a C32 check-encoded string into version and data.
func C32CheckDecode(input string) (byte, []byte, error) {
	// Must be ASCII
	for i := 0; i < len(input); i++ {
		if input[i] >= 128 {
			return 0, nil, errors.New("invalid c32 string: must be ASCII")
		}
	}

	if len(input) < 2 {
		return 0, nil, errors.New("invalid c32 string: size less than 2")
	}

	// Split version and data
	versionChar := input[0]
	dataPart := input[1:]

	decodedData, err := DecodeC32(dataPart)
	if err != nil {
		return 0, nil, err
	}

	if len(decodedData) < 4 {
		return 0, nil, errors.New("invalid c32 string: decoded byte length less than 4")
	}

	// Split data and checksum
	dataBytes := decodedData[:len(decodedData)-4]
	expectedSum := decodedData[len(decodedData)-4:]

	// Decode version
	versionDecoded, err := DecodeC32(string(versionChar))
	if err != nil {
		return 0, nil, err
	}
	version := versionDecoded[0]

	// Verify checksum
	hash1 := sha256.Sum256(append([]byte{version}, dataBytes...))
	hash2 := sha256.Sum256(hash1[:])
	computedSum := hash2[:4]

	if !byteSliceEqual(computedSum, expectedSum) {
		return 0, nil, fmt.Errorf("checksum mismatch")
	}

	return version, dataBytes, nil
}

// DecodeC32Address decodes a C32 address string into version and address bytes.
func DecodeC32Address(c32AddressStr string) (byte, []byte, error) {
	if len(c32AddressStr) <= 5 {
		return 0, nil, errors.New("invalid c32 address: address string smaller than 5 bytes")
	}

	version, data, err := C32CheckDecode(c32AddressStr[1:])
	if err != nil {
		return 0, nil, err
	}

	return version, data, nil
}

// EncodeC32Address encodes a version and address bytes into a C32 address string.
func EncodeC32Address(version byte, data []byte) (string, error) {
	bytes, err := C32CheckEncodePrefixed(version, data, 'S')
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// byteSliceEqual compares two byte slices for equality
func byteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
