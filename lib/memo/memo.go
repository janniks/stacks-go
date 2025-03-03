package memo

import (
	"strings"
	"unicode"
)

// DecodeMemo normalizes the input bytes into a readable string.
// It handles several cases:
// - Converts non-printable characters to spaces
// - Preserves special character sequences like emoji with modifiers
// - Collapses multiple spaces into a single space
// - Trims leading and trailing spaces
//
// This is used to convert raw memo bytes from transaction data into
// a human-readable string representation.
func DecodeMemo(input []byte) string {
	// Handle empty input
	if len(input) == 0 {
		return ""
	}

	// Convert input bytes to string using UTF-8 decoding
	memoStr := string(input)
	runeData := []rune(memoStr)

	// Create a mask to preserve special character sequences
	preserveMask := make([]bool, len(runeData))
	markForPreservation(runeData, preserveMask)

	// Process the string by runes with awareness of special sequences
	var resultBuilder strings.Builder
	resultBuilder.Grow(len(memoStr))

	for i, r := range runeData {
		if preserveMask[i] {
			// Preserve this character as part of a special sequence
			resultBuilder.WriteRune(r)
		} else if unicode.IsPrint(r) {
			// Keep printable characters
			resultBuilder.WriteRune(r)
		} else {
			// Replace non-printable characters with a space
			resultBuilder.WriteRune(' ')
		}
	}

	// Collapse multiple spaces into one and trim
	return collapseAndTrimSpaces(resultBuilder.String())
}

// markForPreservation identifies characters that should be preserved as-is
// (like characters in emoji sequences with zero-width joiners)
func markForPreservation(runeData []rune, preserveMask []bool) {
	for i := 0; i < len(runeData); i++ {
		// Preserve zero-width joiners and adjacent characters
		if runeData[i] == '\u200D' { // zero-width joiner
			preserveMask[i] = true

			// Also preserve characters around ZWJ to keep emoji sequences intact
			if i > 0 {
				preserveMask[i-1] = true
			}
			if i < len(runeData)-1 {
				preserveMask[i+1] = true
			}
		}

		// Preserve combining marks and their base characters
		if i > 0 && isCombiningMark(runeData[i]) {
			preserveMask[i] = true
			preserveMask[i-1] = true
		}
	}
}

// collapseAndTrimSpaces collapses multiple consecutive spaces into a single space
// and trims leading/trailing spaces
func collapseAndTrimSpaces(s string) string {
	wasSpace := false
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		isSpace := unicode.IsSpace(r) || r == '\uFFFD' // Space or replacement character

		if isSpace {
			if !wasSpace {
				builder.WriteRune(' ')
				wasSpace = true
			}
		} else {
			builder.WriteRune(r)
			wasSpace = false
		}
	}

	return strings.TrimSpace(builder.String())
}

// isCombiningMark returns true if the rune is a combining mark
func isCombiningMark(r rune) bool {
	// Combining diacritical marks (U+0300–U+036F)
	if r >= 0x0300 && r <= 0x036F {
		return true
	}

	// Combining spacing marks (various ranges)
	if r >= 0x0900 && r <= 0x097F {
		return true
	}

	// Variation selectors (U+FE00–U+FE0F)
	if r >= 0xFE00 && r <= 0xFE0F {
		return true
	}

	return false
}
