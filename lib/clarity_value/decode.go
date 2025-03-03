package clarity_value

import (
	"encoding/hex"
	"fmt"

	"github.com/janniks/stacks-go/lib/address"
)

// DecodedClarityValue represents a decoded Clarity value with additional metadata
type DecodedClarityValue struct {
	Repr   string      `json:"repr"`
	Hex    string      `json:"hex"`
	TypeID int         `json:"type_id"`
	Value  interface{} `json:"value,omitempty"`
	// Data fields for different types
	Data             interface{}           `json:"data,omitempty"`
	Buffer           string                `json:"buffer,omitempty"`
	List             []DecodedClarityValue `json:"list,omitempty"`
	AddressVersion   int                   `json:"address_version,omitempty"`
	AddressHashBytes string                `json:"address_hash_bytes,omitempty"`
	Address          string                `json:"address,omitempty"`
	ContractName     string                `json:"contract_name,omitempty"`
}

// DecodeClarityValueToObject decodes a ClarityValue into a DecodedClarityValue object
func DecodeClarityValueToObject(val *ClarityValue, deep bool, bytes []byte) (*DecodedClarityValue, error) {
	decoded := &DecodedClarityValue{
		Repr:   val.Value.ReprString(),
		Hex:    hex.EncodeToString(bytes),
		TypeID: int(val.Value.TypePrefix()),
	}

	if deep {
		switch v := val.Value.(type) {
		case IntValue:
			decoded.Value = fmt.Sprintf("%d", v)
		case UIntValue:
			decoded.Value = fmt.Sprintf("%d", v)
		case BoolValue:
			decoded.Value = bool(v)
		case BufferValue:
			decoded.Buffer = hex.EncodeToString(v)
		case ListValue:
			list := make([]DecodedClarityValue, len(v))
			for i, x := range v {
				item, err := DecodeClarityValueToObject(&x, deep, x.SerializedBytes)
				if err != nil {
					return nil, err
				}
				list[i] = *item
			}
			decoded.List = list
		case StringASCIIValue:
			decoded.Data = string(v)
		case StringUTF8Value:
			// Flatten the UTF8 bytes
			var utf8Bytes []byte
			for _, b := range v {
				utf8Bytes = append(utf8Bytes, b...)
			}
			decoded.Data = string(utf8Bytes)
		case PrincipalStandardValue:
			standardPrincipal := StandardPrincipalData(v)
			decoded.AddressVersion = int(standardPrincipal.Version)

			hashSlice := standardPrincipal.Hash[:]
			decoded.AddressHashBytes = hex.EncodeToString(hashSlice)

			addr, err := address.EncodeC32Address(standardPrincipal.Version, hashSlice)
			if err != nil {
				return nil, fmt.Errorf("error converting to C32 address: %w", err)
			}
			decoded.Address = addr
		case PrincipalContractValue:
			contractIdentifier := QualifiedContractIdentifier(v)
			decoded.AddressVersion = int(contractIdentifier.Issuer.Version)

			hashSlice := contractIdentifier.Issuer.Hash[:]
			decoded.AddressHashBytes = hex.EncodeToString(hashSlice)

			addr, err := address.EncodeC32Address(contractIdentifier.Issuer.Version, hashSlice)
			if err != nil {
				return nil, fmt.Errorf("error converting to C32 address: %w", err)
			}
			decoded.Address = addr
			decoded.ContractName = string(contractIdentifier.Name)
		case TupleValue:
			tupleData := make(map[string]DecodedClarityValue)
			for key, value := range v {
				valueDecoded, err := DecodeClarityValueToObject(&value, deep, value.SerializedBytes)
				if err != nil {
					return nil, err
				}
				tupleData[string(key)] = *valueDecoded
			}
			decoded.Data = tupleData
		case OptionalSomeValue:
			optionVal, err := DecodeClarityValueToObject(&v.Value, deep, v.Value.SerializedBytes)
			if err != nil {
				return nil, err
			}
			decoded.Value = optionVal
		case OptionalNoneValue:
			decoded.Value = nil
		case ResponseOkValue:
			respVal, err := DecodeClarityValueToObject(&v.Value, deep, v.Value.SerializedBytes)
			if err != nil {
				return nil, err
			}
			decoded.Value = respVal
		case ResponseErrValue:
			respVal, err := DecodeClarityValueToObject(&v.Value, deep, v.Value.SerializedBytes)
			if err != nil {
				return nil, err
			}
			decoded.Value = respVal
		}
	}

	return decoded, nil
}
