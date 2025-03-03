package memo_test

import (
	"encoding/hex"
	"testing"

	"github.com/janniks/stacks-go/lib/memo"
)

// TestDecodeMemo tests the memo package's DecodeMemo function.
// These test cases are directly translated from the Rust implementation
// and the test vectors are kept identical to ensure compatibility.
func TestDecodeMemo(t *testing.T) {
	// Table-driven tests with vectors from the original Rust implementation
	tests := []struct {
		name     string // Test case name
		input    []byte // Input to be normalized
		expected string // Expected normalized output
	}{
		{
			name:     "Whitespace",
			input:    []byte("hello   world"),
			expected: "hello world",
		},
		{
			name:     "Unknown Unicode",
			input:    []byte("hello\uFFFDworld  test part1   goodbye\uFFFDworld  test part2     "),
			expected: "hello world test part1 goodbye world test part2",
		},
		{
			name:     "Misc BTC Coinbase",
			input:    mustDecodeHex("037e180b04956b4e68627463706f6f6c2f3266646575fabe6d6df77973b452568eb2f43593285804dad9d7ef057eada5ff9f2a1634ec43f514b1020000008e9b20aa0ebfd204924b040000000000"),
			expected: "~ kNhbtcpool/2fdeu mm ys RV 5 (X ~ * 4 C K",
		},
		{
			name:     "Misc BTC Coinbase 2",
			input:    mustDecodeHex("037c180b2cfabe6d6d5e0eb001a2eaea9c5e39b7f54edd5c23eb6e684dab1995191f664658064ba7dc10000000f09f909f092f4632506f6f6c2f6500000000000000000000000000000000000000000000000000000000000000000000000500f3fa0200"),
			expected: "| , mm^ ^9 N \\# nhM fFX K 🐟 /F2Pool/e",
		},
		{
			name:     "Grapheme Extended",
			input:    []byte("👩‍👩‍👧‍👦 hello world"),
			expected: "👩‍👩‍👧‍👦 hello world",
		},
		{
			name:     "Unicode",
			input:    mustDecodeHex("f09f87b3f09f87b12068656c6c6f20776f726c64"),
			expected: "🇳🇱 hello world",
		},
		{
			name:     "Padded Start",
			input:    mustDecodeHex("00000000000068656c6c6f20776f726c64"),
			expected: "hello world",
		},
		{
			name:     "Padded End",
			input:    mustDecodeHex("68656c6c6f20776f726c64000000000000"),
			expected: "hello world",
		},
		{
			name:     "Padded Middle",
			input:    mustDecodeHex("68656c6c6f20776f726c6400000000000068656c6c6f20776f726c64"),
			expected: "hello world hello world",
		},
		{
			name:     "Unicode Scalar",
			input:    []byte("hello worldy̆ test"),
			expected: "hello worldy̆ test",
		},
		{
			name:     "Zero Width Joiner",
			input:    []byte("👨\u200D👩"),
			expected: "👨‍👩",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := memo.DecodeMemo(tt.input)
			if output != tt.expected {
				t.Errorf("DecodeMemo() = %q, want %q", output, tt.expected)
			}
		})
	}
}

// mustDecodeHex decodes a hex string or panics.
// This is only used in tests, so panic is acceptable for invalid test vectors.
func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex test vector: " + err.Error())
	}
	return b
}
