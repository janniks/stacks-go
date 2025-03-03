// Package address implements Stacks address-related functionality
package address

import (
	"fmt"
)

// Bitcoin mainnet and testnet address version bytes
const (
	AddressVersionMainnetSinglesig uint8 = 0
	AddressVersionMainnetMultisig  uint8 = 5
	AddressVersionTestnetSinglesig uint8 = 111
	AddressVersionTestnetMultisig  uint8 = 196
)

// BitcoinAddressType represents the type of Bitcoin address
type BitcoinAddressType int

const (
	// PublicKeyHash represents a pay-to-pubkey-hash address
	PublicKeyHash BitcoinAddressType = iota
	// ScriptHash represents a pay-to-script-hash address
	ScriptHash
)

// BitcoinNetworkType represents the Bitcoin network
type BitcoinNetworkType int

const (
	// Mainnet is the main Bitcoin network
	Mainnet BitcoinNetworkType = iota
	// Testnet is the Bitcoin test network
	Testnet
	// Regtest is the Bitcoin regression test network
	Regtest
)

// BitcoinAddress represents a Bitcoin address
type BitcoinAddress struct {
	AddrType     BitcoinAddressType
	NetworkID    BitcoinNetworkType
	Hash160Bytes [20]byte
}

// VersionByteToAddressType converts a version byte to address type and network
func VersionByteToAddressType(version uint8) (BitcoinAddressType, BitcoinNetworkType, bool) {
	switch version {
	case AddressVersionMainnetSinglesig:
		return PublicKeyHash, Mainnet, true
	case AddressVersionMainnetMultisig:
		return ScriptHash, Mainnet, true
	case AddressVersionTestnetSinglesig:
		return PublicKeyHash, Testnet, true
	case AddressVersionTestnetMultisig:
		return ScriptHash, Testnet, true
	default:
		return 0, 0, false
	}
}

// AddressTypeToVersionByte converts address type and network to a version byte
func AddressTypeToVersionByte(addrType BitcoinAddressType, networkID BitcoinNetworkType) uint8 {
	switch {
	case addrType == PublicKeyHash && networkID == Mainnet:
		return AddressVersionMainnetSinglesig
	case addrType == ScriptHash && networkID == Mainnet:
		return AddressVersionMainnetMultisig
	case addrType == PublicKeyHash && (networkID == Testnet || networkID == Regtest):
		return AddressVersionTestnetSinglesig
	case addrType == ScriptHash && (networkID == Testnet || networkID == Regtest):
		return AddressVersionTestnetMultisig
	default:
		// This shouldn't happen with valid inputs
		return 0
	}
}

// DecodeBitcoinAddress decodes a base58check Bitcoin address string
func DecodeBitcoinAddress(addrb58 string) (*BitcoinAddress, error) {
	bytes, err := DecodeBase58Check(addrb58)
	if err != nil {
		return nil, err
	}

	if len(bytes) != 21 {
		return nil, fmt.Errorf("invalid address: %d bytes", len(bytes))
	}

	version := bytes[0]
	addrType, networkID, valid := VersionByteToAddressType(version)
	if !valid {
		return nil, fmt.Errorf("invalid address: unrecognized version %d", version)
	}

	var payload [20]byte
	copy(payload[:], bytes[1:21])

	return &BitcoinAddress{
		AddrType:     addrType,
		NetworkID:    networkID,
		Hash160Bytes: payload,
	}, nil
}

// EncodeBitcoinAddress encodes a Bitcoin address as a base58check string
func EncodeBitcoinAddress(addr *BitcoinAddress) string {
	version := AddressTypeToVersionByte(addr.AddrType, addr.NetworkID)

	data := make([]byte, 21)
	data[0] = version
	copy(data[1:], addr.Hash160Bytes[:])

	return EncodeBase58Check(data)
}
