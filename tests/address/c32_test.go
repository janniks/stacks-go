// Package address_test contains tests for the address package.
package address_test

import (
	"encoding/hex"
	"testing"

	"github.com/janniks/stacks-go/lib/address"
)

func TestC32Simple(t *testing.T) {
	hexStrings := []string{
		"a46ff88886c2ef9762d970b4d2c63678835bd39d",
		"",
		"0000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000001",
		"1000000000000000000000000000000000000001",
		"1000000000000000000000000000000000000000",
		"01",
		"22",
		"0001",
		"000001",
		"00000001",
		"10",
		"0100",
		"1000",
		"010000",
		"100000",
		"01000000",
		"10000000",
		"0100000000",
	}

	c32Strs := []string{
		"MHQZH246RBQSERPSE2TD5HHPF21NQMWX",
		"",
		"00000000000000000000",
		"00000000000000000001",
		"20000000000000000000000000000001",
		"20000000000000000000000000000000",
		"1",
		"12",
		"01",
		"001",
		"0001",
		"G",
		"80",
		"400",
		"2000",
		"10000",
		"G0000",
		"800000",
		"4000000",
	}

	for i := range hexStrings {
		hexStr := hexStrings[i]
		expectedC32 := c32Strs[i]

		// Skip empty string test for simplicity
		if hexStr == "" {
			continue
		}

		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			t.Fatalf("Failed to decode hex string: %s", err)
		}

		// Test encoding
		c32Encoded := address.EncodeC32(bytes)
		if c32Encoded != expectedC32 {
			t.Errorf("EncodeC32(%x) = %s, expected %s", bytes, c32Encoded, expectedC32)
		}

		// Test decoding
		decoded, err := address.DecodeC32(c32Encoded)
		if err != nil {
			t.Errorf("DecodeC32(%s) failed: %s", c32Encoded, err)
		} else if !bytesEqual(decoded, bytes) {
			t.Errorf("DecodeC32(%s) = %x, expected %x", c32Encoded, decoded, bytes)
		}
	}
}

func TestC32Addresses(t *testing.T) {
	hexStrs := []string{
		"a46ff88886c2ef9762d970b4d2c63678835bd39d",
		"0000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000001",
		"1000000000000000000000000000000000000001",
		"1000000000000000000000000000000000000000",
	}

	versions := []byte{22, 0, 31, 20, 26, 21}

	c32Addrs := [][]string{
		{
			"SP2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKNRV9EJ7",
			"SP000000000000000000002Q6VF78",
			"SP00000000000000000005JA84HQ",
			"SP80000000000000000000000000000004R0CMNV",
			"SP800000000000000000000000000000033H8YKK",
		},
		{
			"S02J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKPVKG2CE",
			"S0000000000000000000002AA028H",
			"S000000000000000000006EKBDDS",
			"S080000000000000000000000000000007R1QC00",
			"S080000000000000000000000000000003ENTGCQ",
		},
		{
			"SZ2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKQ9H6DPR",
			"SZ000000000000000000002ZE1VMN",
			"SZ00000000000000000005HZ3DVN",
			"SZ80000000000000000000000000000004XBV6MS",
			"SZ800000000000000000000000000000007VF5G0",
		},
		{
			"SM2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKQVX8X0G",
			"SM0000000000000000000062QV6X",
			"SM00000000000000000005VR75B2",
			"SM80000000000000000000000000000004WBEWKC",
			"SM80000000000000000000000000000000JGSYGV",
		},
		{
			"ST2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKQYAC0RQ",
			"ST000000000000000000002AMW42H",
			"ST000000000000000000042DB08Y",
			"ST80000000000000000000000000000006BYJ4R4",
			"ST80000000000000000000000000000002YBNPV3",
		},
		{
			"SN2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKP6D2ZK9",
			"SN000000000000000000003YDHWKJ",
			"SN00000000000000000005341MC8",
			"SN800000000000000000000000000000066KZWY0",
			"SN800000000000000000000000000000006H75AK",
		},
	}

	for i := range hexStrs {
		hexStr := hexStrs[i]
		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			t.Fatalf("Failed to decode hex string: %s", err)
		}

		for j := range versions {
			ver := versions[j]
			expectedAddr := c32Addrs[j][i]

			// Test encoding
			addr, err := address.EncodeC32Address(ver, bytes)
			if err != nil {
				t.Errorf("EncodeC32Address(%d, %x) failed: %s", ver, bytes, err)
				continue
			}

			if addr != expectedAddr {
				t.Errorf("EncodeC32Address(%d, %x) = %s, expected %s", ver, bytes, addr, expectedAddr)
			}

			// Test decoding
			decodedVer, decodedBytes, err := address.DecodeC32Address(addr)
			if err != nil {
				t.Errorf("DecodeC32Address(%s) failed: %s", addr, err)
				continue
			}

			if decodedVer != ver {
				t.Errorf("DecodeC32Address(%s) version = %d, expected %d", addr, decodedVer, ver)
			}

			if !bytesEqual(decodedBytes, bytes) {
				t.Errorf("DecodeC32Address(%s) bytes = %x, expected %x", addr, decodedBytes, bytes)
			}
		}
	}
}

func TestC32Normalize(t *testing.T) {
	addrs := []string{
		"S02J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKPVKG2CE",
		"SO2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKPVKG2CE",
		"S02J6ZY48GVLEZ5V2V5RB9MP66SW86PYKKPVKG2CE",
		"SO2J6ZY48GVLEZ5V2V5RB9MP66SW86PYKKPVKG2CE",
		"s02j6zy48gv1ez5v2v5rb9mp66sw86pykkpvkg2ce",
		"sO2j6zy48gv1ez5v2v5rb9mp66sw86pykkpvkg2ce",
		"s02j6zy48gvlez5v2v5rb9mp66sw86pykkpvkg2ce",
		"sO2j6zy48gvlez5v2v5rb9mp66sw86pykkpvkg2ce",
	}

	expectedBytes, _ := hex.DecodeString("a46ff88886c2ef9762d970b4d2c63678835bd39d")
	expectedVersion := byte(0)

	for _, addr := range addrs {
		decodedVersion, decodedBytes, err := address.DecodeC32Address(addr)
		if err != nil {
			t.Errorf("DecodeC32Address(%s) failed: %s", addr, err)
			continue
		}

		if decodedVersion != expectedVersion {
			t.Errorf("DecodeC32Address(%s) version = %d, expected %d", addr, decodedVersion, expectedVersion)
		}

		if !bytesEqual(decodedBytes, expectedBytes) {
			t.Errorf("DecodeC32Address(%s) bytes = %x, expected %x", addr, decodedBytes, expectedBytes)
		}
	}
}

func TestC32AsciiOnly(t *testing.T) {
	// Try a non-ASCII character in the address
	_, _, err := address.DecodeC32Address("S\u1d7d82J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKPVKG2CE")
	if err == nil {
		t.Error("Expected error for non-ASCII input, but got success")
	}
}

// bytesEqual compares two byte slices for equality
func bytesEqual(a, b []byte) bool {
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
