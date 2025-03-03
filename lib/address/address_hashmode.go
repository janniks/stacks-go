package address

import (
	"fmt"
)

// AddressHashMode represents serialization modes for public keys to addresses.
// We support different modes due to legacy compatibility with Stacks v1 addresses.
type AddressHashMode byte

// Address hash mode constants
const (
	// SerializeP2PKH - hash160(public-key), same as bitcoin's p2pkh
	SerializeP2PKH AddressHashMode = 0x00
	// SerializeP2SH - hash160(multisig-redeem-script), same as bitcoin's multisig p2sh
	SerializeP2SH AddressHashMode = 0x01
	// SerializeP2WPKH - hash160(segwit-program-00(p2pkh)), same as bitcoin's p2sh-p2wpkh
	SerializeP2WPKH AddressHashMode = 0x02
	// SerializeP2WSH - hash160(segwit-program-00(public-keys)), same as bitcoin's p2sh-p2wsh
	SerializeP2WSH AddressHashMode = 0x03
	// SerializeP2SHNonSequential - hash160(multisig-redeem-script), same as bitcoin's multisig p2sh (non-sequential signing)
	SerializeP2SHNonSequential AddressHashMode = 0x05
	// SerializeP2WSHNonSequential - hash160(segwit-program-00(public-keys)), same as bitcoin's p2sh-p2wsh (non-sequential signing)
	SerializeP2WSHNonSequential AddressHashMode = 0x07
)

// These constants define address versions for mainnet and testnet
const (
	C32AddressVersionMainnetSinglesig = 22 // 'P'
	C32AddressVersionMainnetMultisig  = 20 // 'M'
	C32AddressVersionTestnetSinglesig = 26 // 'T'
	C32AddressVersionTestnetMultisig  = 21 // 'N'
)

// IsValid returns true if the AddressHashMode is a valid value
func (m AddressHashMode) IsValid() bool {
	switch m {
	case SerializeP2PKH, SerializeP2SH, SerializeP2WPKH, SerializeP2WSH,
		SerializeP2SHNonSequential, SerializeP2WSHNonSequential:
		return true
	}
	return false
}

// ToVersionMainnet converts an AddressHashMode to its corresponding mainnet version
// Returns an error if the mode is not valid
func (m AddressHashMode) ToVersionMainnet() (byte, error) {
	switch m {
	case SerializeP2PKH, SerializeP2WPKH:
		return C32AddressVersionMainnetSinglesig, nil
	case SerializeP2SH, SerializeP2WSH, SerializeP2SHNonSequential, SerializeP2WSHNonSequential:
		return C32AddressVersionMainnetMultisig, nil
	}
	return 0, fmt.Errorf("invalid address hash mode for mainnet conversion: %d", m)
}

// ToVersionTestnet converts an AddressHashMode to its corresponding testnet version
// Returns an error if the mode is not valid
func (m AddressHashMode) ToVersionTestnet() (byte, error) {
	switch m {
	case SerializeP2PKH, SerializeP2WPKH:
		return C32AddressVersionTestnetSinglesig, nil
	case SerializeP2SH, SerializeP2WSH, SerializeP2SHNonSequential, SerializeP2WSHNonSequential:
		return C32AddressVersionTestnetMultisig, nil
	}
	return 0, fmt.Errorf("invalid address hash mode for testnet conversion: %d", m)
}

// FromByte converts a byte value to an AddressHashMode
// Returns an error if the value is invalid
func FromByte(value byte) (AddressHashMode, error) {
	mode := AddressHashMode(value)
	if !mode.IsValid() {
		return 0, fmt.Errorf("invalid address hash mode: %d", value)
	}
	return mode, nil
}
