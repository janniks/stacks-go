package post_condition

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/janniks/stacks-go/lib/address"
	"github.com/janniks/stacks-go/lib/clarity_value"
)

// Asset info type constants
const (
	AssetInfoSTX         byte = 0
	AssetInfoFungible    byte = 1
	AssetInfoNonfungible byte = 2
)

// Principal type constants
const (
	PrincipalOrigin   byte = 0x01
	PrincipalStandard byte = 0x02
	PrincipalContract byte = 0x03
)

// FungibleConditionCode defines condition codes for fungible assets
type FungibleConditionCode byte

const (
	FCSentEq FungibleConditionCode = 0x01
	FCSentGt FungibleConditionCode = 0x02
	FCSentGe FungibleConditionCode = 0x03
	FCSentLt FungibleConditionCode = 0x04
	FCSentLe FungibleConditionCode = 0x05
)

// NonfungibleConditionCode defines condition codes for non-fungible assets
type NonfungibleConditionCode byte

const (
	NFCSent    NonfungibleConditionCode = 0x10
	NFCNotSent NonfungibleConditionCode = 0x11
)

// AssetInfo contains token asset information
type AssetInfo struct {
	Address      address.StacksAddress
	ContractName clarity_value.ClarityName
	AssetName    clarity_value.ClarityName
}

// Principal represents a transaction principal
type Principal struct {
	Type         byte
	Address      address.StacksAddress     // Used for Standard and Contract
	ContractName clarity_value.ClarityName // Used for Contract only
}

// PostCondition represents a transaction post condition
type PostCondition struct {
	Type          byte // AssetInfoSTX, AssetInfoFungible, or AssetInfoNonfungible
	Principal     Principal
	Asset         AssetInfo // Used for Fungible and Nonfungible types
	ConditionCode byte
	Amount        uint64                     // Used for STX and Fungible types
	AssetValue    clarity_value.ClarityValue // Used for Nonfungible type
}

// PostConditionMode represents the mode for post conditions
type PostConditionMode byte

const (
	// PostConditionModeAllow allows transactions even if post conditions are not met
	PostConditionModeAllow PostConditionMode = 0x01
	// PostConditionModeDeny denies transactions if post conditions are not met
	PostConditionModeDeny PostConditionMode = 0x02
)

// PostConditionsResponse represents the response of post conditions decoding
type PostConditionsResponse struct {
	PostConditionMode PostConditionMode `json:"post_condition_mode"`
	PostConditions    []PostCondition   `json:"post_conditions"`
}

// DecodePostCondition decodes a PostCondition from a byte stream
func DecodePostCondition(r io.Reader) (PostCondition, error) {
	// Read asset info type
	var assetType byte
	if err := binary.Read(r, binary.BigEndian, &assetType); err != nil {
		return PostCondition{}, fmt.Errorf("read asset type: %w", err)
	}

	// Read principal
	principal, err := decodePrincipal(r)
	if err != nil {
		return PostCondition{}, err
	}

	// Process based on asset type
	switch assetType {
	case AssetInfoSTX:
		condCode, amount, err := decodeSTXData(r)
		if err != nil {
			return PostCondition{}, err
		}

		return PostCondition{
			Type:          assetType,
			Principal:     principal,
			ConditionCode: condCode,
			Amount:        amount,
		}, nil

	case AssetInfoFungible:
		asset, condCode, amount, err := decodeFungibleData(r)
		if err != nil {
			return PostCondition{}, err
		}

		return PostCondition{
			Type:          assetType,
			Principal:     principal,
			Asset:         asset,
			ConditionCode: condCode,
			Amount:        amount,
		}, nil

	case AssetInfoNonfungible:
		asset, assetValue, condCode, err := decodeNonfungibleData(r)
		if err != nil {
			return PostCondition{}, err
		}

		return PostCondition{
			Type:          assetType,
			Principal:     principal,
			Asset:         asset,
			ConditionCode: condCode,
			AssetValue:    assetValue,
		}, nil

	default:
		return PostCondition{}, fmt.Errorf("unknown asset type: %d", assetType)
	}
}

// decodePrincipal decodes a Principal from a byte stream
func decodePrincipal(r io.Reader) (Principal, error) {
	var principalType byte
	if err := binary.Read(r, binary.BigEndian, &principalType); err != nil {
		return Principal{}, fmt.Errorf("read principal type: %w", err)
	}

	switch principalType {
	case PrincipalOrigin:
		return Principal{Type: principalType}, nil

	case PrincipalStandard:
		addr, err := address.DecodeStacksAddress(r)
		if err != nil {
			return Principal{}, fmt.Errorf("decode standard address: %w", err)
		}

		return Principal{
			Type:    principalType,
			Address: addr,
		}, nil

	case PrincipalContract:
		addr, err := address.DecodeStacksAddress(r)
		if err != nil {
			return Principal{}, fmt.Errorf("decode contract address: %w", err)
		}

		reader, ok := r.(*bytes.Reader)
		if !ok {
			return Principal{}, fmt.Errorf("expected bytes.Reader")
		}

		name, err := clarity_value.DecodeClarityName(reader)
		if err != nil {
			return Principal{}, fmt.Errorf("decode contract name: %w", err)
		}

		return Principal{
			Type:         principalType,
			Address:      addr,
			ContractName: name,
		}, nil

	default:
		return Principal{}, fmt.Errorf("unknown principal type: %d", principalType)
	}
}

// decodeSTXData decodes STX-specific data (condition code and amount)
func decodeSTXData(r io.Reader) (byte, uint64, error) {
	// Read condition code
	var condCode byte
	if err := binary.Read(r, binary.BigEndian, &condCode); err != nil {
		return 0, 0, fmt.Errorf("read condition code: %w", err)
	}

	if err := validateFungibleConditionCode(condCode); err != nil {
		return 0, 0, err
	}

	// Read amount
	var amount uint64
	if err := binary.Read(r, binary.BigEndian, &amount); err != nil {
		return 0, 0, fmt.Errorf("read amount: %w", err)
	}

	return condCode, amount, nil
}

// decodeFungibleData decodes fungible asset data (asset info, condition code, and amount)
func decodeFungibleData(r io.Reader) (AssetInfo, byte, uint64, error) {
	// Read asset info
	asset, err := decodeAssetInfo(r)
	if err != nil {
		return AssetInfo{}, 0, 0, err
	}

	// Read condition code
	var condCode byte
	if err := binary.Read(r, binary.BigEndian, &condCode); err != nil {
		return AssetInfo{}, 0, 0, fmt.Errorf("read condition code: %w", err)
	}

	if err := validateFungibleConditionCode(condCode); err != nil {
		return AssetInfo{}, 0, 0, err
	}

	// Read amount
	var amount uint64
	if err := binary.Read(r, binary.BigEndian, &amount); err != nil {
		return AssetInfo{}, 0, 0, fmt.Errorf("read amount: %w", err)
	}

	return asset, condCode, amount, nil
}

// decodeNonfungibleData decodes non-fungible asset data (asset info, asset value, and condition code)
func decodeNonfungibleData(r io.Reader) (AssetInfo, clarity_value.ClarityValue, byte, error) {
	// Read asset info
	asset, err := decodeAssetInfo(r)
	if err != nil {
		return AssetInfo{}, clarity_value.ClarityValue{}, 0, err
	}

	// Get bytes.Reader for clarity value operations
	reader, ok := r.(*bytes.Reader)
	if !ok {
		return AssetInfo{}, clarity_value.ClarityValue{}, 0, fmt.Errorf("expected bytes.Reader")
	}

	// Capture starting position for serialized bytes
	startPos := reader.Size() - int64(reader.Len())

	// Decode the clarity value
	val, err := clarity_value.DecodeClarityValue(reader, false)
	if err != nil {
		return AssetInfo{}, clarity_value.ClarityValue{}, 0, fmt.Errorf("decode clarity value: %w", err)
	}

	// Calculate bytes read and capture them
	endPos := reader.Size() - int64(reader.Len())

	// Extract serialized bytes
	serializedBytes, err := extractSerializedBytes(reader, startPos, endPos)
	if err != nil {
		return AssetInfo{}, clarity_value.ClarityValue{}, 0, err
	}

	// Add serialized bytes to the clarity value
	val.SerializedBytes = serializedBytes

	// Read condition code
	var condCode byte
	if err := binary.Read(r, binary.BigEndian, &condCode); err != nil {
		return AssetInfo{}, clarity_value.ClarityValue{}, 0, fmt.Errorf("read condition code: %w", err)
	}

	if err := validateNonfungibleConditionCode(condCode); err != nil {
		return AssetInfo{}, clarity_value.ClarityValue{}, 0, err
	}

	return asset, val, condCode, nil
}

// extractSerializedBytes extracts bytes from a reader between start and end positions
func extractSerializedBytes(reader *bytes.Reader, startPos, endPos int64) ([]byte, error) {
	// Go back to read all the bytes
	if _, err := reader.Seek(startPos, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek to start: %w", err)
	}

	// Read the bytes into a slice
	decoded := make([]byte, endPos-startPos)
	if _, err := io.ReadFull(reader, decoded); err != nil {
		return nil, fmt.Errorf("read value bytes: %w", err)
	}

	// Restore position
	if _, err := reader.Seek(endPos, io.SeekStart); err != nil {
		return nil, fmt.Errorf("restore position: %w", err)
	}

	return decoded, nil
}

// decodeAssetInfo decodes asset info from a byte stream
func decodeAssetInfo(r io.Reader) (AssetInfo, error) {
	addr, err := address.DecodeStacksAddress(r)
	if err != nil {
		return AssetInfo{}, fmt.Errorf("decode address: %w", err)
	}

	reader, ok := r.(*bytes.Reader)
	if !ok {
		return AssetInfo{}, fmt.Errorf("expected bytes.Reader")
	}

	contractName, err := clarity_value.DecodeClarityName(reader)
	if err != nil {
		return AssetInfo{}, fmt.Errorf("decode contract name: %w", err)
	}

	assetName, err := clarity_value.DecodeClarityName(reader)
	if err != nil {
		return AssetInfo{}, fmt.Errorf("decode asset name: %w", err)
	}

	return AssetInfo{
		Address:      addr,
		ContractName: contractName,
		AssetName:    assetName,
	}, nil
}

// validateFungibleConditionCode checks if a condition code is valid for fungible assets
func validateFungibleConditionCode(code byte) error {
	switch FungibleConditionCode(code) {
	case FCSentEq, FCSentGt, FCSentGe, FCSentLt, FCSentLe:
		return nil
	default:
		return fmt.Errorf("invalid fungible condition code: %d", code)
	}
}

// validateNonfungibleConditionCode checks if a condition code is valid for non-fungible assets
func validateNonfungibleConditionCode(code byte) error {
	switch NonfungibleConditionCode(code) {
	case NFCSent, NFCNotSent:
		return nil
	default:
		return fmt.Errorf("invalid non-fungible condition code: %d", code)
	}
}

// DecodeTxPostConditions decodes a transaction's post conditions from bytes
func DecodeTxPostConditions(data []byte) (*PostConditionsResponse, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for post conditions")
	}

	resp := &PostConditionsResponse{
		PostConditionMode: PostConditionMode(data[0]),
		PostConditions:    []PostCondition{},
	}

	if len(data) > 4 {
		// Next 4 bytes are array length but we don't need to use it
		// as we'll just read until we run out of data

		// Next bytes are serialized post condition items
		postConditionBytes := data[5:]
		reader := bytes.NewReader(postConditionBytes)

		for reader.Len() > 0 {
			postCondition, err := DecodePostCondition(reader)
			if err != nil {
				return nil, fmt.Errorf("error deserializing post condition: %w", err)
			}
			resp.PostConditions = append(resp.PostConditions, postCondition)
		}
	}

	return resp, nil
}
