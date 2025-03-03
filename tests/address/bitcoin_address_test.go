package address_test

import (
	"testing"

	"github.com/janniks/stacks-go/lib/address"
)

func TestBitcoinAddressEncoding(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		addrType  address.BitcoinAddressType
		networkID address.BitcoinNetworkType
	}{
		{
			name:      "mainnet singlesig address",
			address:   "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", // Satoshi's address
			addrType:  address.PublicKeyHash,
			networkID: address.Mainnet,
		},
		{
			name:      "mainnet multisig address",
			address:   "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", // Common example address
			addrType:  address.ScriptHash,
			networkID: address.Mainnet,
		},
		{
			name:      "testnet singlesig address",
			address:   "mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef", // Example testnet address
			addrType:  address.PublicKeyHash,
			networkID: address.Testnet,
		},
		{
			name:      "testnet multisig address",
			address:   "2MzQwSSnBHWHqSAqtTVQ6v47XtaisrJa1Vc", // Example testnet P2SH address
			addrType:  address.ScriptHash,
			networkID: address.Testnet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test decoding
			decoded, err := address.DecodeBitcoinAddress(tt.address)
			if err != nil {
				t.Fatalf("DecodeBitcoinAddress(%s) failed: %v", tt.address, err)
			}

			if decoded.AddrType != tt.addrType {
				t.Errorf("DecodeBitcoinAddress(%s) got address type %v, want %v", tt.address, decoded.AddrType, tt.addrType)
			}

			if decoded.NetworkID != tt.networkID {
				t.Errorf("DecodeBitcoinAddress(%s) got network ID %v, want %v", tt.address, decoded.NetworkID, tt.networkID)
			}

			// Test encoding
			encoded := address.EncodeBitcoinAddress(decoded)
			if encoded != tt.address {
				t.Errorf("EncodeBitcoinAddress() got %s, want %s", encoded, tt.address)
			}
		})
	}
}

func TestInvalidAddresses(t *testing.T) {
	invalidAddresses := []string{
		"",                                    // Empty string
		"1",                                   // Too short
		"1QJQxDas5JhdiXhEbNS14iNjgZMGDweis",   // Invalid checksum
		"1QJQxDas5JhdiXhEbNS14iNjgZMGDweisss", // Too long
		"1QJQxDas5JhdiXhEbNS14iNjgZMGDweiO0",  // Invalid character 'O'
		"9QJQxDas5JhdiXhEbNS14iNjgZMGDweiss",  // Invalid version byte
	}

	for _, addr := range invalidAddresses {
		t.Run("invalid:"+addr, func(t *testing.T) {
			_, err := address.DecodeBitcoinAddress(addr)
			if err == nil {
				t.Errorf("DecodeBitcoinAddress(%s) should have failed", addr)
			}
		})
	}
}

func TestVersionByteConversion(t *testing.T) {
	tests := []struct {
		addrType  address.BitcoinAddressType
		networkID address.BitcoinNetworkType
		version   uint8
	}{
		{address.PublicKeyHash, address.Mainnet, address.AddressVersionMainnetSinglesig},
		{address.ScriptHash, address.Mainnet, address.AddressVersionMainnetMultisig},
		{address.PublicKeyHash, address.Testnet, address.AddressVersionTestnetSinglesig},
		{address.ScriptHash, address.Testnet, address.AddressVersionTestnetMultisig},
		{address.PublicKeyHash, address.Regtest, address.AddressVersionTestnetSinglesig},
		{address.ScriptHash, address.Regtest, address.AddressVersionTestnetMultisig},
	}

	for _, tt := range tests {
		t.Run("version_conversion", func(t *testing.T) {
			version := address.AddressTypeToVersionByte(tt.addrType, tt.networkID)
			if version != tt.version {
				t.Errorf("AddressTypeToVersionByte(%v, %v) = %v, want %v",
					tt.addrType, tt.networkID, version, tt.version)
			}

			addrType, networkID, valid := address.VersionByteToAddressType(tt.version)
			if !valid {
				t.Errorf("VersionByteToAddressType(%v) returned not valid", tt.version)
			}

			if addrType != tt.addrType {
				t.Errorf("VersionByteToAddressType(%v) got address type %v, want %v",
					tt.version, addrType, tt.addrType)
			}

			// For Regtest, we expect Testnet to be returned when decoding
			expectedNetworkID := tt.networkID
			if tt.networkID == address.Regtest {
				expectedNetworkID = address.Testnet
			}

			if networkID != expectedNetworkID {
				t.Errorf("VersionByteToAddressType(%v) got network ID %v, want %v",
					tt.version, networkID, expectedNetworkID)
			}
		})
	}
}
