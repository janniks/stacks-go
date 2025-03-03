package post_condition_test

import (
	"bufio"
	"compress/gzip"
	"encoding/hex"
	"os"
	"testing"

	"github.com/janniks/stacks-go/lib/post_condition"
)

func TestDecodePostConditionSamples(t *testing.T) {
	// Read the gzipped sample data file
	sampleFile, err := os.Open("../gz/sampled-post-conditions.txt.gz")
	if err != nil {
		t.Fatalf("Failed to open sample file: %v", err)
	}
	defer sampleFile.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(sampleFile)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Read line by line
	scanner := bufio.NewScanner(gzipReader)
	for scanner.Scan() {
		line := scanner.Text()

		// Decode the hex string to bytes
		inputBytes, err := hex.DecodeString(line)
		if err != nil {
			t.Fatalf("Failed to decode hex string: %v", err)
		}

		// Decode the post conditions
		_, err = post_condition.DecodeTxPostConditions(inputBytes)
		if err != nil {
			t.Fatalf("Failed to decode post conditions for input %s: %v", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading sample file: %v", err)
	}
}
