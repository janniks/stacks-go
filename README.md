# `stacks-go`

A Go library for the Stacks blockchain.

## Overview

This library aims to provide a pure Go implementation for interacting with the Stacks blockchain ecosystem.

## Installation

```bash
go get github.com/hirosystems/stacks-go
```

## Usage

```go
package main

import (
	"fmt"
	stacks "github.com/hirosystems/stacks-go"
)

func main() {
	// Get the library version
	fmt.Println("Library version:", stacks.Version())

	// Basic hello function
	fmt.Println(stacks.Hello())
}
```

## Project Structure

The project follows a clean structure:

```
stacks-go/
├── *.go           # Core library code
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
