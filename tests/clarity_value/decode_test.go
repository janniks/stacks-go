package clarity_value_test

import (
	"encoding/hex"
	"testing"

	"github.com/janniks/stacks-go/lib/clarity_value"
)

func TestDecodeClarityValueToObject(t *testing.T) {
	testCases := []struct {
		name       string
		clarityVal func() *clarity_value.ClarityValue
		bytes      []byte
		deep       bool
		validate   func(t *testing.T, result *clarity_value.DecodedClarityValue)
	}{
		{
			name: "Int value",
			clarityVal: func() *clarity_value.ClarityValue {
				val := clarity_value.IntValue(42)
				return &clarity_value.ClarityValue{
					Value:           val,
					SerializedBytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a}, // 42 in big-endian
				}
			},
			bytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a},
			deep:  true,
			validate: func(t *testing.T, result *clarity_value.DecodedClarityValue) {
				if result.Repr != clarity_value.IntValue(42).ReprString() {
					t.Errorf("Expected repr %s, got %s", clarity_value.IntValue(42).ReprString(), result.Repr)
				}
				if result.TypeID != int(clarity_value.PrefixInt) {
					t.Errorf("Expected type_id %d, got %d", int(clarity_value.PrefixInt), result.TypeID)
				}
				if result.Hex != "000000000000002a" {
					t.Errorf("Expected hex %s, got %s", "000000000000002a", result.Hex)
				}
				if result.Value != "42" {
					t.Errorf("Expected value %s, got %v", "42", result.Value)
				}
			},
		},
		{
			name: "Bool value (true)",
			clarityVal: func() *clarity_value.ClarityValue {
				val := clarity_value.BoolValue(true)
				return &clarity_value.ClarityValue{
					Value:           val,
					SerializedBytes: []byte{},
				}
			},
			bytes: []byte{},
			deep:  true,
			validate: func(t *testing.T, result *clarity_value.DecodedClarityValue) {
				if result.Repr != clarity_value.BoolValue(true).ReprString() {
					t.Errorf("Expected repr %s, got %s", clarity_value.BoolValue(true).ReprString(), result.Repr)
				}
				if result.TypeID != int(clarity_value.PrefixBoolTrue) {
					t.Errorf("Expected type_id %d, got %d", int(clarity_value.PrefixBoolTrue), result.TypeID)
				}
				if value, ok := result.Value.(bool); !ok || !value {
					t.Errorf("Expected boolean true, got %v", result.Value)
				}
			},
		},
		{
			name: "String ASCII value",
			clarityVal: func() *clarity_value.ClarityValue {
				val := clarity_value.StringASCIIValue([]byte("hello"))
				return &clarity_value.ClarityValue{
					Value:           val,
					SerializedBytes: []byte{0x05, 'h', 'e', 'l', 'l', 'o'}, // 5 bytes + "hello"
				}
			},
			bytes: []byte{0x05, 'h', 'e', 'l', 'l', 'o'},
			deep:  true,
			validate: func(t *testing.T, result *clarity_value.DecodedClarityValue) {
				if result.Repr != clarity_value.StringASCIIValue([]byte("hello")).ReprString() {
					t.Errorf("Expected repr %s, got %s", clarity_value.StringASCIIValue([]byte("hello")).ReprString(), result.Repr)
				}
				if result.TypeID != int(clarity_value.PrefixStringASCII) {
					t.Errorf("Expected type_id %d, got %d", int(clarity_value.PrefixStringASCII), result.TypeID)
				}
				if result.Hex != "0568656c6c6f" {
					t.Errorf("Expected hex %s, got %s", "0568656c6c6f", result.Hex)
				}
				if data, ok := result.Data.(string); !ok || data != "hello" {
					t.Errorf("Expected data 'hello', got %v", result.Data)
				}
			},
		},
		{
			name: "When deep is false",
			clarityVal: func() *clarity_value.ClarityValue {
				val := clarity_value.IntValue(42)
				return &clarity_value.ClarityValue{
					Value:           val,
					SerializedBytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a},
				}
			},
			bytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a},
			deep:  false,
			validate: func(t *testing.T, result *clarity_value.DecodedClarityValue) {
				if result.Repr != clarity_value.IntValue(42).ReprString() {
					t.Errorf("Expected repr %s, got %s", clarity_value.IntValue(42).ReprString(), result.Repr)
				}
				if result.TypeID != int(clarity_value.PrefixInt) {
					t.Errorf("Expected type_id %d, got %d", int(clarity_value.PrefixInt), result.TypeID)
				}
				if result.Hex != "000000000000002a" {
					t.Errorf("Expected hex %s, got %s", "000000000000002a", result.Hex)
				}
				if result.Value != nil {
					t.Errorf("Expected value to be nil since deep is false, got %v", result.Value)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := tc.clarityVal()
			result, err := clarity_value.DecodeClarityValueToObject(val, tc.deep, tc.bytes)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			tc.validate(t, result)
		})
	}
}

func TestDecodeClarityValueToObjectWithSerializedBytes(t *testing.T) {
	// Create a test case where we use the SerializedBytes from the ClarityValue
	intVal := clarity_value.IntValue(42)
	serializedBytes := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a}

	clarityVal := clarity_value.ClarityValue{
		Value:           intVal,
		SerializedBytes: serializedBytes,
	}

	result, err := clarity_value.DecodeClarityValueToObject(&clarityVal, true, clarityVal.SerializedBytes)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Repr != intVal.ReprString() {
		t.Errorf("Expected repr %s, got %s", intVal.ReprString(), result.Repr)
	}

	if result.TypeID != int(clarity_value.PrefixInt) {
		t.Errorf("Expected type_id %d, got %d", int(clarity_value.PrefixInt), result.TypeID)
	}

	if result.Hex != hex.EncodeToString(serializedBytes) {
		t.Errorf("Expected hex %s, got %s", hex.EncodeToString(serializedBytes), result.Hex)
	}

	if result.Value != "42" {
		t.Errorf("Expected value %s, got %v", "42", result.Value)
	}
}
