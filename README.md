# `stacks-go`

A Go library for the Stacks blockchain.

## Overview

This library aims to provide a pure Go implementation for interacting with the Stacks blockchain ecosystem.

## Installation

```bash
go get github.com/janniks/stacks-go
```

## Usage

```go
package main

import (
	"fmt"
	"encoding/hex"
	"github.com/janniks/stacks-go/address"
)

func main() {
	// Encode byte data to base58
	encoded := address.EncodeBase58([]byte{0, 0, 42})
	fmt.Println("Base58 encoded:", encoded) // "112f"

	// Decode from base58
	decoded, err := address.DecodeBase58("112f")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("Decoded bytes: %v\n", decoded) // [0 0 42]

	// Base58Check encode (with checksum)
	addr, _ := hex.DecodeString("00f8917303bfa8ef24f292e8fa1419b20460ba064d")
	checkEncoded := address.EncodeBase58Check(addr)
	fmt.Println("Address:", checkEncoded) // "1PfJpZsjreyVrqeoAfabrRwwjQyoSQMmHH"
}
```

## Packages

### Address Package

The `address` package provides utilities for working with Stacks addresses and related encodings.

- `base58.go`: Implementation of Base58 and Base58Check encoding/decoding
  - `EncodeBase58`: Encodes bytes to Base58 string
  - `DecodeBase58`: Decodes Base58 string to bytes
  - `EncodeBase58Check`: Encodes bytes with a checksum (Base58Check)
  - `DecodeBase58Check`: Decodes Base58Check string and verifies checksum

## Project Structure

The project follows a clean structure:

```
stacks-go/
├── *.go           # Core library code
├── address/       # Address-related functionality
│   └── base58.go  # Base58 encoding/decoding implementation
├── tests/         # Test files
│   └── *_test.go  # Tests for the library
└── examples/      # Example applications
    ├── basic/     # Basic usage example
    └── hello/     # Hello world example
```

## Development

### Prerequisites

- Go 1.19 or higher

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

## Examples

Check out the examples directory for more usage examples.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
