// Package address implements Stacks address-related functionality
package address

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
)

// Base58 alphabet used for encoding and decoding
const base58Chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Pre-computed base58 digit values
var base58Digits [128]int

func init() {
	// Initialize base58Digits lookup table
	for i := 0; i < len(base58Digits); i++ {
		base58Digits[i] = -1
	}
	for i := 0; i < len(base58Chars); i++ {
		base58Digits[base58Chars[i]] = i
	}
}

// DecodeBase58 decodes a base58-encoded string into a byte slice
func DecodeBase58(input string) ([]byte, error) {
	// Quick return for empty input
	if len(input) == 0 {
		return []byte{}, nil
	}

	// Allocate enough space for the decoded data
	// 11/15 is just over log_256(58)
	result := make([]byte, 1+len(input)*11/15)

	// Count leading '1's (base58 encoding of 0)
	var leadingZeros int
	for i := 0; i < len(input) && input[i] == '1'; i++ {
		leadingZeros++
	}

	// Convert from base58 to base256
	for i := 0; i < len(input); i++ {
		c := input[i]
		// Check if character is in valid range
		if c >= 128 {
			return nil, fmt.Errorf("invalid base58 character: %c", c)
		}

		// Get digit value using lookup table
		digit := base58Digits[c]
		if digit == -1 {
			return nil, fmt.Errorf("invalid base58 character: %c", c)
		}

		// Multiply existing result by 58 and add the new digit
		carry := digit
		for j := len(result) - 1; j >= 0; j-- {
			carry += int(result[j]) * 58
			result[j] = byte(carry & 0xff)
			carry >>= 8
		}
	}

	// Skip leading zeros in result and prepend any leading 1s from input
	i := 0
	for i < len(result) && result[i] == 0 {
		i++
	}

	// Create the final result with leading zeros
	final := make([]byte, leadingZeros+(len(result)-i))
	copy(final[leadingZeros:], result[i:])

	return final, nil
}

// DecodeBase58Check decodes a base58check-encoded string
func DecodeBase58Check(input string) ([]byte, error) {
	decoded, err := DecodeBase58(input)
	if err != nil {
		return nil, err
	}

	if len(decoded) < 4 {
		return nil, errors.New("base58check data too short for checksum")
	}

	// Split data and checksum
	ckStart := len(decoded) - 4
	data := decoded[:ckStart]
	checksum := decoded[ckStart:]

	// Calculate expected checksum
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	expected := hash2[:4]

	// Verify checksum
	for i := 0; i < 4; i++ {
		if checksum[i] != expected[i] {
			return nil, errors.New("base58check checksum mismatch")
		}
	}

	return data, nil
}

// EncodeBase58 encodes a byte slice as a base58 string
func EncodeBase58(data []byte) string {
	// Quick return for empty data
	if len(data) == 0 {
		return ""
	}

	// Count leading zeros
	var leadingZeros int
	for i := 0; i < len(data) && data[i] == 0; i++ {
		leadingZeros++
	}

	// Allocate enough space for the encoded data
	// 7/5 is just over log_58(256)
	result := make([]byte, 1+len(data)*7/5)
	var resultLen int

	// Convert from base256 to base58
	for i := 0; i < len(data); i++ {
		carry := int(data[i])

		j := 0
		for ; j < resultLen || carry != 0; j++ {
			if j < resultLen {
				carry += 256 * int(result[j])
			}
			result[j] = byte(carry % 58)
			carry /= 58
		}
		resultLen = j
	}

	// Skip leading zeros in result
	i := resultLen - 1

	// Convert to base58 characters and prepend any leading 1s
	output := strings.Repeat("1", leadingZeros)
	for ; i >= 0; i-- {
		output += string(base58Chars[result[i]])
	}

	return output
}

// EncodeBase58Check encodes data with a 4-byte checksum
func EncodeBase58Check(data []byte) string {
	// Calculate checksum (double SHA-256)
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	checksum := hash2[:4]

	// Append checksum to data
	withChecksum := make([]byte, len(data)+4)
	copy(withChecksum, data)
	copy(withChecksum[len(data):], checksum)

	// Encode with base58
	return EncodeBase58(withChecksum)
}
