// Package address_test contains tests for the address package.
package address_test

import (
	"encoding/hex"
	"testing"

	"github.com/janniks/stacks-go/lib/address"
)

func TestBase58Encode(t *testing.T) {
	// Basics
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0}, "1"},
		{[]byte{1}, "2"},
		{[]byte{58}, "21"},
		{[]byte{13, 36}, "211"},
		// Leading zeroes
		{[]byte{0, 13, 36}, "1211"},
		{[]byte{0, 0, 0, 0, 13, 36}, "1111211"},
	}

	for _, tc := range testCases {
		result := address.EncodeBase58(tc.input)
		if result != tc.expected {
			t.Errorf("EncodeBase58(%x) = %s, expected %s", tc.input, result, tc.expected)
		}
	}

	// Address test case
	addr, _ := hex.DecodeString("00f8917303bfa8ef24f292e8fa1419b20460ba064d")
	result := address.EncodeBase58Check(addr)
	expected := "1PfJpZsjreyVrqeoAfabrRwwjQyoSQMmHH"
	if result != expected {
		t.Errorf("EncodeBase58Check(%x) = %s, expected %s", addr, result, expected)
	}
}

func TestBase58Decode(t *testing.T) {
	// Basics
	testCases := []struct {
		input    string
		expected []byte
	}{
		{"1", []byte{0}},
		{"2", []byte{1}},
		{"21", []byte{58}},
		{"211", []byte{13, 36}},
		// Leading zeroes
		{"1211", []byte{0, 13, 36}},
		{"111211", []byte{0, 0, 0, 13, 36}},
	}

	for _, tc := range testCases {
		result, err := address.DecodeBase58(tc.input)
		if err != nil {
			t.Errorf("DecodeBase58(%s) returned error: %v", tc.input, err)
			continue
		}

		// Compare bytes
		if len(result) != len(tc.expected) {
			t.Errorf("DecodeBase58(%s) = %x, expected %x (length mismatch)", tc.input, result, tc.expected)
			continue
		}

		for i := range result {
			if result[i] != tc.expected[i] {
				t.Errorf("DecodeBase58(%s) = %x, expected %x", tc.input, result, tc.expected)
				break
			}
		}
	}

	// Address test case
	expected, _ := hex.DecodeString("00f8917303bfa8ef24f292e8fa1419b20460ba064d")
	result, err := address.DecodeBase58Check("1PfJpZsjreyVrqeoAfabrRwwjQyoSQMmHH")
	if err != nil {
		t.Errorf("DecodeBase58Check returned error: %v", err)
	} else {
		// Compare bytes
		if len(result) != len(expected) {
			t.Errorf("DecodeBase58Check returned %x, expected %x (length mismatch)", result, expected)
		} else {
			for i := range result {
				if result[i] != expected[i] {
					t.Errorf("DecodeBase58Check returned %x, expected %x", result, expected)
					break
				}
			}
		}
	}
}

func TestBase58Roundtrip(t *testing.T) {
	// Test the same roundtrip case from the Rust tests
	s := "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs"

	// Decode
	decoded, err := address.DecodeBase58Check(s)
	if err != nil {
		t.Fatalf("DecodeBase58Check(%s) returned error: %v", s, err)
	}

	// Re-encode
	reencoded := address.EncodeBase58Check(decoded)

	// Compare
	if reencoded != s {
		t.Errorf("Roundtrip failed: original=%s, reencoded=%s", s, reencoded)
	}

	// Another roundtrip
	redecoded, err := address.DecodeBase58Check(reencoded)
	if err != nil {
		t.Fatalf("DecodeBase58Check(%s) returned error: %v", reencoded, err)
	}

	// Compare bytes
	if len(redecoded) != len(decoded) {
		t.Errorf("Decoded byte arrays have different lengths: %d vs %d", len(redecoded), len(decoded))
	} else {
		for i := range redecoded {
			if redecoded[i] != decoded[i] {
				t.Errorf("Decoded byte arrays differ at position %d: %x vs %x", i, redecoded[i], decoded[i])
				break
			}
		}
	}
}
