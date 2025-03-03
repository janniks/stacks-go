package clarity_value

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"

	"github.com/janniks/stacks-go/lib/address"
)

// Constants
const (
	MaxStringLen          = 128
	MaxValueSize          = 1024 * 1024 // 1MB
	ContractMinNameLength = 1
	ContractMaxNameLength = 40
)

// TypePrefix represents the type prefix for serialized Clarity values
type TypePrefix byte

// Type prefixes for serialized Clarity values
const (
	PrefixInt               TypePrefix = 0
	PrefixUInt              TypePrefix = 1
	PrefixBuffer            TypePrefix = 2
	PrefixBoolTrue          TypePrefix = 3
	PrefixBoolFalse         TypePrefix = 4
	PrefixPrincipalStandard TypePrefix = 5
	PrefixPrincipalContract TypePrefix = 6
	PrefixResponseOk        TypePrefix = 7
	PrefixResponseErr       TypePrefix = 8
	PrefixOptionalNone      TypePrefix = 9
	PrefixOptionalSome      TypePrefix = 10
	PrefixList              TypePrefix = 11
	PrefixTuple             TypePrefix = 12
	PrefixStringASCII       TypePrefix = 13
	PrefixStringUTF8        TypePrefix = 14
)

// ClarityValue represents a Clarity value with optional serialized bytes
type ClarityValue struct {
	SerializedBytes []byte
	Value           Value
}

// NewClarityValueWithBytes creates a new ClarityValue with serialized bytes
func NewClarityValueWithBytes(serializedBytes []byte, value Value) ClarityValue {
	return ClarityValue{
		SerializedBytes: serializedBytes,
		Value:           value,
	}
}

// NewClarityValue creates a new ClarityValue without serialized bytes
func NewClarityValue(value Value) ClarityValue {
	return ClarityValue{
		Value: value,
	}
}

// Value represents a Clarity value
type Value interface {
	TypePrefix() TypePrefix
	ReprString() string
	TypeSignature() string
}

// IntValue represents a Clarity integer value
type IntValue int64

// TypePrefix returns the type prefix for IntValue
func (v IntValue) TypePrefix() TypePrefix {
	return PrefixInt
}

// ReprString returns the string representation of IntValue
func (v IntValue) ReprString() string {
	return fmt.Sprintf("%d", v)
}

// TypeSignature returns the type signature of IntValue
func (v IntValue) TypeSignature() string {
	return "int"
}

// UIntValue represents a Clarity unsigned integer value
type UIntValue uint64

// TypePrefix returns the type prefix for UIntValue
func (v UIntValue) TypePrefix() TypePrefix {
	return PrefixUInt
}

// ReprString returns the string representation of UIntValue
func (v UIntValue) ReprString() string {
	return fmt.Sprintf("u%d", v)
}

// TypeSignature returns the type signature of UIntValue
func (v UIntValue) TypeSignature() string {
	return "uint"
}

// BoolValue represents a Clarity boolean value
type BoolValue bool

// TypePrefix returns the type prefix for BoolValue
func (v BoolValue) TypePrefix() TypePrefix {
	if bool(v) {
		return PrefixBoolTrue
	}
	return PrefixBoolFalse
}

// ReprString returns the string representation of BoolValue
func (v BoolValue) ReprString() string {
	return fmt.Sprintf("%t", v)
}

// TypeSignature returns the type signature of BoolValue
func (v BoolValue) TypeSignature() string {
	return "bool"
}

// BufferValue represents a Clarity buffer value
type BufferValue []byte

// TypePrefix returns the type prefix for BufferValue
func (v BufferValue) TypePrefix() TypePrefix {
	return PrefixBuffer
}

// ReprString returns the string representation of BufferValue
func (v BufferValue) ReprString() string {
	return hex.EncodeToString(v)
}

// TypeSignature returns the type signature of BufferValue
func (v BufferValue) TypeSignature() string {
	return fmt.Sprintf("(buff %d)", len(v))
}

// ListValue represents a Clarity list value
type ListValue []ClarityValue

// TypePrefix returns the type prefix for ListValue
func (v ListValue) TypePrefix() TypePrefix {
	return PrefixList
}

// ReprString returns the string representation of ListValue
func (v ListValue) ReprString() string {
	var buffer bytes.Buffer
	buffer.WriteString("(list")
	for _, val := range v {
		buffer.WriteString(" ")
		buffer.WriteString(val.Value.ReprString())
	}
	buffer.WriteString(")")
	return buffer.String()
}

// TypeSignature returns the type signature of ListValue
func (v ListValue) TypeSignature() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("(list %d ", len(v)))
	if len(v) > 0 {
		buffer.WriteString(v[0].Value.TypeSignature())
	} else {
		buffer.WriteString("UnknownType")
	}
	buffer.WriteString(")")
	return buffer.String()
}

// StringUTF8Value represents a Clarity UTF-8 string value
type StringUTF8Value [][]byte

// NewStringUTF8Value creates a new StringUTF8Value from a byte slice
func NewStringUTF8Value(bytes []byte) StringUTF8Value {
	str := string(bytes)
	result := make([][]byte, 0, len(str))
	for _, r := range str {
		charBytes := []byte(string(r))
		result = append(result, charBytes)
	}
	return result
}

// TypePrefix returns the type prefix for StringUTF8Value
func (v StringUTF8Value) TypePrefix() TypePrefix {
	return PrefixStringUTF8
}

// ReprString returns the string representation of StringUTF8Value
func (v StringUTF8Value) ReprString() string {
	var buffer bytes.Buffer
	buffer.WriteString("u\"")
	for _, c := range v {
		if len(c) > 1 {
			// We escape extended charset
			buffer.WriteString(fmt.Sprintf("\\u{%s}", hex.EncodeToString(c)))
		} else {
			// We render an ASCII char, escaped
			buffer.WriteString(fmt.Sprintf("%s", escapeASCII(c[0])))
		}
	}
	buffer.WriteString("\"")
	return buffer.String()
}

// TypeSignature returns the type signature of StringUTF8Value
func (v StringUTF8Value) TypeSignature() string {
	return fmt.Sprintf("(string-utf8 %d)", len(v)*4)
}

// StringASCIIValue represents a Clarity ASCII string value
type StringASCIIValue []byte

// TypePrefix returns the type prefix for StringASCIIValue
func (v StringASCIIValue) TypePrefix() TypePrefix {
	return PrefixStringASCII
}

// ReprString returns the string representation of StringASCIIValue
func (v StringASCIIValue) ReprString() string {
	var buffer bytes.Buffer
	buffer.WriteString("\"")
	for _, c := range v {
		buffer.WriteString(escapeASCII(c))
	}
	buffer.WriteString("\"")
	return buffer.String()
}

// TypeSignature returns the type signature of StringASCIIValue
func (v StringASCIIValue) TypeSignature() string {
	return fmt.Sprintf("(string-ascii %d)", len(v))
}

// StandardPrincipalData represents a Clarity standard principal
type StandardPrincipalData struct {
	Version byte
	Hash    [20]byte
}

// PrincipalStandardValue represents a Clarity standard principal value
type PrincipalStandardValue StandardPrincipalData

// TypePrefix returns the type prefix for PrincipalStandardValue
func (v PrincipalStandardValue) TypePrefix() TypePrefix {
	return PrefixPrincipalStandard
}

// ReprString returns the string representation of PrincipalStandardValue
func (v PrincipalStandardValue) ReprString() string {
	addr, err := address.EncodeC32Address(v.Version, v.Hash[:])
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return fmt.Sprintf("'%s", addr)
}

// TypeSignature returns the type signature of PrincipalStandardValue
func (v PrincipalStandardValue) TypeSignature() string {
	return "principal"
}

// QualifiedContractIdentifier represents a contract identifier
type QualifiedContractIdentifier struct {
	Issuer StandardPrincipalData
	Name   ClarityName
}

// PrincipalContractValue represents a Clarity contract principal value
type PrincipalContractValue QualifiedContractIdentifier

// TypePrefix returns the type prefix for PrincipalContractValue
func (v PrincipalContractValue) TypePrefix() TypePrefix {
	return PrefixPrincipalContract
}

// ReprString returns the string representation of PrincipalContractValue
func (v PrincipalContractValue) ReprString() string {
	addr, err := address.EncodeC32Address(v.Issuer.Version, v.Issuer.Hash[:])
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return fmt.Sprintf("'%s.%s", addr, v.Name)
}

// TypeSignature returns the type signature of PrincipalContractValue
func (v PrincipalContractValue) TypeSignature() string {
	return "principal"
}

// TupleValue represents a Clarity tuple value
type TupleValue map[ClarityName]ClarityValue

// TypePrefix returns the type prefix for TupleValue
func (v TupleValue) TypePrefix() TypePrefix {
	return PrefixTuple
}

// ReprString returns the string representation of TupleValue
func (v TupleValue) ReprString() string {
	var buffer bytes.Buffer
	buffer.WriteString("(tuple")
	// Create a deterministic order of keys
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)

	for _, key := range keys {
		clarityKey := ClarityName(key)
		value := v[clarityKey]
		buffer.WriteString(fmt.Sprintf(" (%s ", clarityKey))
		buffer.WriteString(value.Value.ReprString())
		buffer.WriteString(")")
	}
	buffer.WriteString(")")
	return buffer.String()
}

// TypeSignature returns the type signature of TupleValue
func (v TupleValue) TypeSignature() string {
	var buffer bytes.Buffer
	buffer.WriteString("(tuple")
	// Create a deterministic order of keys
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)

	for _, key := range keys {
		clarityKey := ClarityName(key)
		value := v[clarityKey]
		buffer.WriteString(fmt.Sprintf(" (%s ", clarityKey))
		buffer.WriteString(value.Value.TypeSignature())
		buffer.WriteString(")")
	}
	buffer.WriteString(")")
	return buffer.String()
}

// OptionalSomeValue represents a Clarity optional some value
type OptionalSomeValue struct {
	Value ClarityValue
}

// TypePrefix returns the type prefix for OptionalSomeValue
func (v OptionalSomeValue) TypePrefix() TypePrefix {
	return PrefixOptionalSome
}

// ReprString returns the string representation of OptionalSomeValue
func (v OptionalSomeValue) ReprString() string {
	return fmt.Sprintf("(some %s)", v.Value.Value.ReprString())
}

// TypeSignature returns the type signature of OptionalSomeValue
func (v OptionalSomeValue) TypeSignature() string {
	return fmt.Sprintf("(optional %s)", v.Value.Value.TypeSignature())
}

// OptionalNoneValue represents a Clarity optional none value
type OptionalNoneValue struct{}

// TypePrefix returns the type prefix for OptionalNoneValue
func (v OptionalNoneValue) TypePrefix() TypePrefix {
	return PrefixOptionalNone
}

// ReprString returns the string representation of OptionalNoneValue
func (v OptionalNoneValue) ReprString() string {
	return "none"
}

// TypeSignature returns the type signature of OptionalNoneValue
func (v OptionalNoneValue) TypeSignature() string {
	return "(optional UnknownType)"
}

// ResponseOkValue represents a Clarity response ok value
type ResponseOkValue struct {
	Value ClarityValue
}

// TypePrefix returns the type prefix for ResponseOkValue
func (v ResponseOkValue) TypePrefix() TypePrefix {
	return PrefixResponseOk
}

// ReprString returns the string representation of ResponseOkValue
func (v ResponseOkValue) ReprString() string {
	return fmt.Sprintf("(ok %s)", v.Value.Value.ReprString())
}

// TypeSignature returns the type signature of ResponseOkValue
func (v ResponseOkValue) TypeSignature() string {
	return fmt.Sprintf("(response %s UnknownType)", v.Value.Value.TypeSignature())
}

// ResponseErrValue represents a Clarity response error value
type ResponseErrValue struct {
	Value ClarityValue
}

// TypePrefix returns the type prefix for ResponseErrValue
func (v ResponseErrValue) TypePrefix() TypePrefix {
	return PrefixResponseErr
}

// ReprString returns the string representation of ResponseErrValue
func (v ResponseErrValue) ReprString() string {
	return fmt.Sprintf("(err %s)", v.Value.Value.ReprString())
}

// TypeSignature returns the type signature of ResponseErrValue
func (v ResponseErrValue) TypeSignature() string {
	return fmt.Sprintf("(response UnknownType %s)", v.Value.Value.TypeSignature())
}

// Helper function to escape ASCII characters
func escapeASCII(c byte) string {
	switch c {
	case '\a':
		return "\\a"
	case '\b':
		return "\\b"
	case '\t':
		return "\\t"
	case '\n':
		return "\\n"
	case '\v':
		return "\\v"
	case '\f':
		return "\\f"
	case '\r':
		return "\\r"
	case '"':
		return "\\\""
	case '\\':
		return "\\\\"
	default:
		if c < 32 || c >= 127 {
			return fmt.Sprintf("\\x%02x", c)
		}
		return string(c)
	}
}

// GuardedString types

// ClarityName represents a validated Clarity name
type ClarityName string

var clarityNameRegex = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9]|[-_!?+<>=/*])*$|^[-+=/*]$|^[<>]=?$`)

// ValidateClarityName validates a string as a Clarity name
func ValidateClarityName(s string) (ClarityName, error) {
	if len(s) > MaxStringLen {
		return "", fmt.Errorf("bad name value ClarityName, %s", s)
	}
	if clarityNameRegex.MatchString(s) {
		return ClarityName(s), nil
	}
	return "", fmt.Errorf("bad name value ClarityName, %s", s)
}

// MustClarityName creates a ClarityName, panicking if invalid
func MustClarityName(s string) ClarityName {
	name, err := ValidateClarityName(s)
	if err != nil {
		panic(err)
	}
	return name
}

// String returns the string value of ClarityName
func (c ClarityName) String() string {
	return string(c)
}

// ContractName represents a validated contract name
type ContractName string

var contractNameRegex = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9]|[-_])*$|^__transient$`)

// ValidateContractName validates a string as a contract name
func ValidateContractName(s string) (ContractName, error) {
	if len(s) > MaxStringLen {
		return "", fmt.Errorf("bad name value ContractName, %s", s)
	}
	if contractNameRegex.MatchString(s) {
		return ContractName(s), nil
	}
	return "", fmt.Errorf("bad name value ContractName, %s", s)
}

// MustContractName creates a ContractName, panicking if invalid
func MustContractName(s string) ContractName {
	name, err := ValidateContractName(s)
	if err != nil {
		panic(err)
	}
	return name
}

// String returns the string value of ContractName
func (c ContractName) String() string {
	return string(c)
}
