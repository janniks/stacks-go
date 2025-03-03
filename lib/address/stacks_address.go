package address

import (
	"encoding/binary"
	"fmt"
	"io"
)

// StacksAddress represents a Stacks blockchain address
type StacksAddress struct {
	Version byte
	Hash160 [20]byte
}

// NewStacksAddress creates a new StacksAddress with the given version and hash160
func NewStacksAddress(version byte, hash160 [20]byte) StacksAddress {
	return StacksAddress{
		Version: version,
		Hash160: hash160,
	}
}

// FromString creates a StacksAddress from a C32-encoded string
func FromString(s string) (StacksAddress, error) {
	version, bytes, err := DecodeC32Address(s)
	if err != nil {
		return StacksAddress{}, fmt.Errorf("error decoding c32 address: %w", err)
	}

	var hash160 [20]byte
	copy(hash160[:], bytes)

	return StacksAddress{
		Version: version,
		Hash160: hash160,
	}, nil
}

// DecodeStacksAddress decodes a StacksAddress from a byte stream
func DecodeStacksAddress(r io.Reader) (StacksAddress, error) {
	var version byte
	var hash160 [20]byte

	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return StacksAddress{}, fmt.Errorf("failed to read address version: %w", err)
	}

	if _, err := io.ReadFull(r, hash160[:]); err != nil {
		return StacksAddress{}, fmt.Errorf("failed to read address hash160: %w", err)
	}

	return StacksAddress{
		Version: version,
		Hash160: hash160,
	}, nil
}

// String returns the C32-encoded string representation of the address
func (a StacksAddress) String() string {
	addr, err := EncodeC32Address(a.Version, a.Hash160[:])
	if err != nil {
		return fmt.Sprintf("invalid address: %v", err)
	}
	return addr
}
