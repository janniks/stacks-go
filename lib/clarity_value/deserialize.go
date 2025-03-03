package clarity_value

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// DeserializeError represents an error during deserialization
type DeserializeError struct {
	message string
}

func (e DeserializeError) Error() string {
	return e.message
}

// NewDeserializeError creates a new DeserializeError
func NewDeserializeError(message string) DeserializeError {
	return DeserializeError{message: message}
}

// DecodeClarityName deserializes a ClarityName from a byte reader
func DecodeClarityName(r *bytes.Reader) (ClarityName, error) {
	lenByte, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	if lenByte > MaxStringLen {
		return "", NewDeserializeError(fmt.Sprintf("Failed to deserialize clarity name: too long: %d", lenByte))
	}

	data := make([]byte, lenByte)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}

	s := string(data)
	return ValidateClarityName(s)
}

// DecodeContractName deserializes a ContractName from a byte reader
func DecodeContractName(r *bytes.Reader) (ContractName, error) {
	lenByte, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	if uint8(lenByte) < ContractMinNameLength || uint8(lenByte) > ContractMaxNameLength {
		return "", NewDeserializeError(fmt.Sprintf("Failed to deserialize contract name: too short or too long: %d", lenByte))
	}

	data := make([]byte, lenByte)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}

	s := string(data)
	return ValidateContractName(s)
}

// DecodeStandardPrincipalData deserializes a StandardPrincipalData from a byte reader
func DecodeStandardPrincipalData(r *bytes.Reader) (StandardPrincipalData, error) {
	version, err := r.ReadByte()
	if err != nil {
		return StandardPrincipalData{}, err
	}

	var hash [20]byte
	if _, err := io.ReadFull(r, hash[:]); err != nil {
		return StandardPrincipalData{}, err
	}

	return StandardPrincipalData{
		Version: version,
		Hash:    hash,
	}, nil
}

// DecodeClarityValue deserializes a ClarityValue from a byte reader
func DecodeClarityValue(r *bytes.Reader, withBytes bool) (ClarityValue, error) {
	return decodeClarityValueInternal(r, 0, withBytes)
}

// decodeClarityValueInternal handles the recursive deserialization of ClarityValue
func decodeClarityValueInternal(r *bytes.Reader, depth uint8, withBytes bool) (ClarityValue, error) {
	if depth >= 16 {
		return ClarityValue{}, NewDeserializeError(fmt.Sprintf("TypeSignatureTooDeep: %d", depth))
	}

	startPos := r.Size() - int64(r.Len())

	header, err := r.ReadByte()
	if err != nil {
		return ClarityValue{}, err
	}

	prefix := TypePrefix(header)

	var value Value
	var deserializeErr error

	switch prefix {
	case PrefixInt:
		var buf [16]byte
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return ClarityValue{}, err
		}
		value = IntValue(int64(binary.BigEndian.Uint64(buf[8:])))

	case PrefixUInt:
		var buf [16]byte
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return ClarityValue{}, err
		}
		value = UIntValue(binary.BigEndian.Uint64(buf[8:]))

	case PrefixBuffer:
		var bufLen uint32
		if err := binary.Read(r, binary.BigEndian, &bufLen); err != nil {
			return ClarityValue{}, err
		}
		if bufLen > MaxValueSize {
			return ClarityValue{}, NewDeserializeError("Illegal buffer type size")
		}
		data := make([]byte, bufLen)
		if _, err := io.ReadFull(r, data); err != nil {
			return ClarityValue{}, err
		}
		value = BufferValue(data)

	case PrefixBoolTrue:
		value = BoolValue(true)

	case PrefixBoolFalse:
		value = BoolValue(false)

	case PrefixPrincipalStandard:
		principal, err := DecodeStandardPrincipalData(r)
		if err != nil {
			return ClarityValue{}, err
		}
		value = PrincipalStandardValue(principal)

	case PrefixPrincipalContract:
		issuer, err := DecodeStandardPrincipalData(r)
		if err != nil {
			return ClarityValue{}, err
		}
		name, err := DecodeClarityName(r)
		if err != nil {
			return ClarityValue{}, err
		}
		value = PrincipalContractValue(QualifiedContractIdentifier{
			Issuer: issuer,
			Name:   name,
		})

	case PrefixResponseOk:
		innerValue, err := decodeClarityValueInternal(r, depth+1, withBytes)
		if err != nil {
			return ClarityValue{}, err
		}
		value = ResponseOkValue{Value: innerValue}

	case PrefixResponseErr:
		innerValue, err := decodeClarityValueInternal(r, depth+1, withBytes)
		if err != nil {
			return ClarityValue{}, err
		}
		value = ResponseErrValue{Value: innerValue}

	case PrefixOptionalNone:
		value = OptionalNoneValue{}

	case PrefixOptionalSome:
		innerValue, err := decodeClarityValueInternal(r, depth+1, withBytes)
		if err != nil {
			return ClarityValue{}, err
		}
		value = OptionalSomeValue{Value: innerValue}

	case PrefixList:
		var listLen uint32
		if err := binary.Read(r, binary.BigEndian, &listLen); err != nil {
			return ClarityValue{}, err
		}
		if listLen > MaxValueSize {
			return ClarityValue{}, NewDeserializeError("Illegal list type size")
		}
		items := make([]ClarityValue, listLen)
		for i := uint32(0); i < listLen; i++ {
			item, err := decodeClarityValueInternal(r, depth+1, withBytes)
			if err != nil {
				return ClarityValue{}, err
			}
			items[i] = item
		}
		value = ListValue(items)

	case PrefixTuple:
		var tupleLen uint32
		if err := binary.Read(r, binary.BigEndian, &tupleLen); err != nil {
			return ClarityValue{}, err
		}
		if tupleLen > MaxValueSize {
			return ClarityValue{}, NewDeserializeError("Illegal tuple type size")
		}
		data := make(TupleValue)
		for i := uint32(0); i < tupleLen; i++ {
			key, err := DecodeClarityName(r)
			if err != nil {
				return ClarityValue{}, err
			}
			val, err := decodeClarityValueInternal(r, depth+1, withBytes)
			if err != nil {
				return ClarityValue{}, err
			}
			data[key] = val
		}
		value = data

	case PrefixStringASCII:
		var bufLen uint32
		if err := binary.Read(r, binary.BigEndian, &bufLen); err != nil {
			return ClarityValue{}, err
		}
		if bufLen > MaxValueSize {
			return ClarityValue{}, NewDeserializeError("Illegal string-ascii type size")
		}
		data := make([]byte, bufLen)
		if _, err := io.ReadFull(r, data); err != nil {
			return ClarityValue{}, err
		}
		value = StringASCIIValue(data)

	case PrefixStringUTF8:
		var totalLen uint32
		if err := binary.Read(r, binary.BigEndian, &totalLen); err != nil {
			return ClarityValue{}, err
		}
		if totalLen > MaxValueSize {
			return ClarityValue{}, NewDeserializeError("Illegal string-utf8 type size")
		}
		data := make([]byte, totalLen)
		if _, err := io.ReadFull(r, data); err != nil {
			return ClarityValue{}, err
		}
		value = NewStringUTF8Value(data)

	default:
		return ClarityValue{}, errors.New("Bad type prefix")
	}

	if withBytes {
		endPos := r.Size() - int64(r.Len())
		allBytes := make([]byte, endPos-startPos)

		// Remember current position
		currentPos := r.Size() - int64(r.Len())

		// Go back to read all the bytes
		_, err = r.Seek(startPos, io.SeekStart)
		if err != nil {
			return ClarityValue{}, err
		}

		_, err = io.ReadFull(r, allBytes)
		if err != nil {
			return ClarityValue{}, err
		}

		// Restore position
		_, err = r.Seek(currentPos, io.SeekStart)
		if err != nil {
			return ClarityValue{}, err
		}

		return NewClarityValueWithBytes(allBytes, value), deserializeErr
	}

	return NewClarityValue(value), deserializeErr
}
