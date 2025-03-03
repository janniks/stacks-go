package clarity_value_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/janniks/stacks-go/lib/clarity_value"
)

func TestDecodeClarityName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Valid clarity name",
			input:    "056e616d653f", // Length 5, "name?"
			expected: "name?",
			hasError: false,
		},
		{
			name:     "Empty name",
			input:    "00",
			hasError: true,
		},
		{
			name:     "Too long name",
			input:    "81" + hex.EncodeToString(bytes.Repeat([]byte("a"), 129)), // Length 129, exceeds max
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Decode the hex string to bytes
			inputBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Fatalf("Failed to decode hex input: %v", err)
			}

			// Create a bytes.Reader from the input bytes
			reader := bytes.NewReader(inputBytes)

			// Call the DecodeClarityName function
			result, err := clarity_value.DecodeClarityName(reader)

			// Check if the error matches expectations
			if tc.hasError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error is expected, check the result
			if !tc.hasError {
				if string(result) != tc.expected {
					t.Errorf("Expected result %q, got %q", tc.expected, result)
				}
			}
		})
	}
}

func TestDecodeContractName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Valid contract name",
			input:    "08636f6e7472616374", // Length 8, "contract"
			expected: "contract",
			hasError: false,
		},
		{
			name:     "Empty contract name",
			input:    "00",
			hasError: true,
		},
		{
			name:     "Too long contract name",
			input:    "29" + hex.EncodeToString(bytes.Repeat([]byte("a"), 41)), // Length 41, exceeds max
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Decode the hex string to bytes
			inputBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Fatalf("Failed to decode hex input: %v", err)
			}

			// Create a bytes.Reader from the input bytes
			reader := bytes.NewReader(inputBytes)

			// Call the DecodeContractName function
			result, err := clarity_value.DecodeContractName(reader)

			// Check if the error matches expectations
			if tc.hasError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error is expected, check the result
			if !tc.hasError {
				if string(result) != tc.expected {
					t.Errorf("Expected result %q, got %q", tc.expected, result)
				}
			}
		})
	}
}

func TestDecodeStandardPrincipalData(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedVer  byte
		expectedHash string
		hasError     bool
	}{
		{
			name:         "Valid standard principal data",
			input:        "002b3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d", // Version 0, 20 bytes hash
			expectedVer:  0,
			expectedHash: "2b3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d",
			hasError:     false,
		},
		{
			name:     "Incomplete data",
			input:    "00112233", // Only version and 3 bytes
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Decode the hex string to bytes
			inputBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Fatalf("Failed to decode hex input: %v", err)
			}

			// Create a bytes.Reader from the input bytes
			reader := bytes.NewReader(inputBytes)

			// Call the DecodeStandardPrincipalData function
			result, err := clarity_value.DecodeStandardPrincipalData(reader)

			// Check if the error matches expectations
			if tc.hasError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error is expected, check the result
			if !tc.hasError && err == nil {
				if result.Version != tc.expectedVer {
					t.Errorf("Expected version %d, got %d", tc.expectedVer, result.Version)
				}

				hashHex := hex.EncodeToString(result.Hash[:])
				if hashHex != tc.expectedHash {
					t.Errorf("Expected hash %s, got %s", tc.expectedHash, hashHex)
				}
			}
		})
	}
}

func TestDecodeClarityValue(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		withBytes   bool
		checkBytes  bool
		expectedHex string
		valueCheck  func(t *testing.T, value clarity_value.Value)
		hasError    bool
	}{
		{
			name:        "Int value",
			input:       "000000000000000000000000000000000a", // Int(10)
			withBytes:   true,
			checkBytes:  true,
			expectedHex: "000000000000000000000000000000000a",
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				intValue, ok := value.(clarity_value.IntValue)
				if !ok {
					t.Fatalf("Expected IntValue, got %T", value)
				}
				if intValue != 10 {
					t.Errorf("Expected int value 10, got %d", intValue)
				}
			},
			hasError: false,
		},
		{
			name:      "UInt value",
			input:     "010000000000000000000000000000000f", // UInt(15)
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				uintValue, ok := value.(clarity_value.UIntValue)
				if !ok {
					t.Fatalf("Expected UIntValue, got %T", value)
				}
				if uintValue != 15 {
					t.Errorf("Expected uint value 15, got %d", uintValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Boolean true",
			input:     "03", // Bool(true)
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				boolValue, ok := value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", value)
				}
				if boolValue != true {
					t.Errorf("Expected bool value true, got %v", boolValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Boolean false",
			input:     "04", // Bool(false)
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				boolValue, ok := value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected BoolValue, got %T", value)
				}
				if boolValue != false {
					t.Errorf("Expected bool value false, got %v", boolValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Buffer",
			input:     "0200000003010203", // Buffer([1, 2, 3])
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				bufferValue, ok := value.(clarity_value.BufferValue)
				if !ok {
					t.Fatalf("Expected BufferValue, got %T", value)
				}
				expected := []byte{1, 2, 3}
				if !bytes.Equal(bufferValue, expected) {
					t.Errorf("Expected buffer %v, got %v", expected, bufferValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Optional none",
			input:     "09", // OptionalNone
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				_, ok := value.(clarity_value.OptionalNoneValue)
				if !ok {
					t.Fatalf("Expected OptionalNoneValue, got %T", value)
				}
			},
			hasError: false,
		},
		{
			name:      "Optional some",
			input:     "0a03", // OptionalSome(true)
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				someValue, ok := value.(clarity_value.OptionalSomeValue)
				if !ok {
					t.Fatalf("Expected OptionalSomeValue, got %T", value)
				}

				innerValue, ok := someValue.Value.Value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected inner BoolValue, got %T", someValue.Value.Value)
				}
				if innerValue != true {
					t.Errorf("Expected inner value true, got %v", innerValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Response ok",
			input:     "0703", // ResponseOk(true)
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				okValue, ok := value.(clarity_value.ResponseOkValue)
				if !ok {
					t.Fatalf("Expected ResponseOkValue, got %T", value)
				}

				innerValue, ok := okValue.Value.Value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected inner BoolValue, got %T", okValue.Value.Value)
				}
				if innerValue != true {
					t.Errorf("Expected inner value true, got %v", innerValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Response err",
			input:     "0804", // ResponseErr(false)
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				errValue, ok := value.(clarity_value.ResponseErrValue)
				if !ok {
					t.Fatalf("Expected ResponseErrValue, got %T", value)
				}

				innerValue, ok := errValue.Value.Value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected inner BoolValue, got %T", errValue.Value.Value)
				}
				if innerValue != false {
					t.Errorf("Expected inner value false, got %v", innerValue)
				}
			},
			hasError: false,
		},
		{
			name:      "List",
			input:     "0b00000002030a03", // List[true, OptionalSome(true)]
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				listValue, ok := value.(clarity_value.ListValue)
				if !ok {
					t.Fatalf("Expected ListValue, got %T", value)
				}
				if len(listValue) != 2 {
					t.Fatalf("Expected list length 2, got %d", len(listValue))
				}

				// Check first element
				firstValue, ok := listValue[0].Value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected first element BoolValue, got %T", listValue[0].Value)
				}
				if firstValue != true {
					t.Errorf("Expected first element true, got %v", firstValue)
				}

				// Check second element
				secondValue, ok := listValue[1].Value.(clarity_value.OptionalSomeValue)
				if !ok {
					t.Fatalf("Expected second element OptionalSomeValue, got %T", listValue[1].Value)
				}

				innerValue, ok := secondValue.Value.Value.(clarity_value.BoolValue)
				if !ok {
					t.Fatalf("Expected inner BoolValue, got %T", secondValue.Value.Value)
				}
				if innerValue != true {
					t.Errorf("Expected inner value true, got %v", innerValue)
				}
			},
			hasError: false,
		},
		{
			name:      "Invalid type prefix",
			input:     "ff00000000", // Invalid prefix
			withBytes: false,
			valueCheck: func(t *testing.T, value clarity_value.Value) {
				// Should not be called due to error
			},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Decode the hex string to bytes
			inputBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Fatalf("Failed to decode hex input: %v", err)
			}

			// Create a bytes.Reader from the input bytes
			reader := bytes.NewReader(inputBytes)

			// Call the DecodeClarityValue function
			result, err := clarity_value.DecodeClarityValue(reader, tc.withBytes)

			// Check if the error matches expectations
			if tc.hasError && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error is expected, check the result
			if !tc.hasError {
				tc.valueCheck(t, result.Value)

				// Check serialized bytes if required
				if tc.checkBytes {
					bytesHex := hex.EncodeToString(result.SerializedBytes)
					if bytesHex != tc.expectedHex {
						t.Errorf("Expected serialized bytes %s, got %s", tc.expectedHex, bytesHex)
					}
				}
			}
		})
	}
}
