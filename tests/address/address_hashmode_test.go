package address_test

import (
	"testing"

	"github.com/janniks/stacks-go/lib/address"
)

func TestAddressHashModeIsValid(t *testing.T) {
	tests := []struct {
		mode     address.AddressHashMode
		expected bool
	}{
		{address.SerializeP2PKH, true},
		{address.SerializeP2SH, true},
		{address.SerializeP2WPKH, true},
		{address.SerializeP2WSH, true},
		{address.SerializeP2SHNonSequential, true},
		{address.SerializeP2WSHNonSequential, true},
		{address.AddressHashMode(0x04), false},
		{address.AddressHashMode(0x06), false},
		{address.AddressHashMode(0xFF), false},
	}

	for _, test := range tests {
		result := test.mode.IsValid()
		if result != test.expected {
			t.Errorf("IsValid() for mode %v = %v, expected %v", test.mode, result, test.expected)
		}
	}
}

func TestAddressHashModeToVersionMainnet(t *testing.T) {
	// Test valid conversions
	validTests := []struct {
		mode     address.AddressHashMode
		expected byte
	}{
		{address.SerializeP2PKH, address.C32AddressVersionMainnetSinglesig},
		{address.SerializeP2WPKH, address.C32AddressVersionMainnetSinglesig},
		{address.SerializeP2SH, address.C32AddressVersionMainnetMultisig},
		{address.SerializeP2WSH, address.C32AddressVersionMainnetMultisig},
		{address.SerializeP2SHNonSequential, address.C32AddressVersionMainnetMultisig},
		{address.SerializeP2WSHNonSequential, address.C32AddressVersionMainnetMultisig},
	}

	for _, test := range validTests {
		result, err := test.mode.ToVersionMainnet()
		if err != nil {
			t.Errorf("ToVersionMainnet() for mode %v returned error: %v", test.mode, err)
		}
		if result != test.expected {
			t.Errorf("ToVersionMainnet() for mode %v = %v, expected %v", test.mode, result, test.expected)
		}
	}

	// Test invalid conversion
	invalidMode := address.AddressHashMode(0x04)
	_, err := invalidMode.ToVersionMainnet()
	if err == nil {
		t.Errorf("ToVersionMainnet() for invalid mode %v expected error, got nil", invalidMode)
	}
}

func TestAddressHashModeToVersionTestnet(t *testing.T) {
	// Test valid conversions
	validTests := []struct {
		mode     address.AddressHashMode
		expected byte
	}{
		{address.SerializeP2PKH, address.C32AddressVersionTestnetSinglesig},
		{address.SerializeP2WPKH, address.C32AddressVersionTestnetSinglesig},
		{address.SerializeP2SH, address.C32AddressVersionTestnetMultisig},
		{address.SerializeP2WSH, address.C32AddressVersionTestnetMultisig},
		{address.SerializeP2SHNonSequential, address.C32AddressVersionTestnetMultisig},
		{address.SerializeP2WSHNonSequential, address.C32AddressVersionTestnetMultisig},
	}

	for _, test := range validTests {
		result, err := test.mode.ToVersionTestnet()
		if err != nil {
			t.Errorf("ToVersionTestnet() for mode %v returned error: %v", test.mode, err)
		}
		if result != test.expected {
			t.Errorf("ToVersionTestnet() for mode %v = %v, expected %v", test.mode, result, test.expected)
		}
	}

	// Test invalid conversion
	invalidMode := address.AddressHashMode(0x04)
	_, err := invalidMode.ToVersionTestnet()
	if err == nil {
		t.Errorf("ToVersionTestnet() for invalid mode %v expected error, got nil", invalidMode)
	}
}

func TestFromByte(t *testing.T) {
	// Test valid conversions
	validTests := []struct {
		value    byte
		expected address.AddressHashMode
	}{
		{0x00, address.SerializeP2PKH},
		{0x01, address.SerializeP2SH},
		{0x02, address.SerializeP2WPKH},
		{0x03, address.SerializeP2WSH},
		{0x05, address.SerializeP2SHNonSequential},
		{0x07, address.SerializeP2WSHNonSequential},
	}

	for _, test := range validTests {
		result, err := address.FromByte(test.value)
		if err != nil {
			t.Errorf("FromByte(%v) returned error: %v", test.value, err)
		}
		if result != test.expected {
			t.Errorf("FromByte(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}

	// Test invalid conversion
	invalidValues := []byte{0x04, 0x06, 0xFF}
	for _, value := range invalidValues {
		_, err := address.FromByte(value)
		if err == nil {
			t.Errorf("FromByte(%v) expected error, got nil", value)
		}
	}
}
