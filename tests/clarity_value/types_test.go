package clarity_value_test

import (
	"strings"
	"testing"

	"github.com/janniks/stacks-go/lib/clarity_value"
)

func TestClarityName(t *testing.T) {
	validNames := []string{
		"hello",
		"hello-world",
		"hello_world",
		"hello123",
		"hello!?+<>=/*",
		"-",
		"+",
		"=",
		"/",
		"*",
		"<",
		">",
		"<=",
		">=",
	}

	invalidNames := []string{
		"",
		"123hello",
		"hello world",
		"$hello",
		"hello$",
	}

	for _, name := range validNames {
		_, err := clarity_value.ValidateClarityName(name)
		if err != nil {
			t.Errorf("Expected valid ClarityName: %s", name)
		}
	}

	for _, name := range invalidNames {
		_, err := clarity_value.ValidateClarityName(name)
		if err == nil {
			t.Errorf("Expected invalid ClarityName: %s", name)
		}
	}
}

func TestContractName(t *testing.T) {
	validNames := []string{
		"hello",
		"hello-world",
		"hello_world",
		"hello123",
		"__transient",
	}

	invalidNames := []string{
		"",
		"123hello",
		"hello world",
		"hello!",
		"hello?",
	}

	for _, name := range validNames {
		_, err := clarity_value.ValidateContractName(name)
		if err != nil {
			t.Errorf("Expected valid ContractName: %s", name)
		}
	}

	for _, name := range invalidNames {
		_, err := clarity_value.ValidateContractName(name)
		if err == nil {
			t.Errorf("Expected invalid ContractName: %s", name)
		}
	}
}

func TestClarityValueRepresentations(t *testing.T) {
	testCases := []struct {
		name          string
		value         clarity_value.Value
		reprString    string
		typeSignature string
	}{
		{
			"Int",
			clarity_value.IntValue(123),
			"123",
			"int",
		},
		{
			"UInt",
			clarity_value.UIntValue(123),
			"u123",
			"uint",
		},
		{
			"BoolTrue",
			clarity_value.BoolValue(true),
			"true",
			"bool",
		},
		{
			"BoolFalse",
			clarity_value.BoolValue(false),
			"false",
			"bool",
		},
		{
			"Buffer",
			clarity_value.BufferValue([]byte{0x01, 0x02, 0x03}),
			"010203",
			"(buff 3)",
		},
		{
			"StringASCII",
			clarity_value.StringASCIIValue([]byte("hello")),
			"\"hello\"",
			"(string-ascii 5)",
		},
		{
			"OptionalNone",
			clarity_value.OptionalNoneValue{},
			"none",
			"(optional UnknownType)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if repr := tc.value.ReprString(); repr != tc.reprString {
				t.Errorf("Expected ReprString %s, got %s", tc.reprString, repr)
			}
			if typeSig := tc.value.TypeSignature(); typeSig != tc.typeSignature {
				t.Errorf("Expected TypeSignature %s, got %s", tc.typeSignature, typeSig)
			}
		})
	}
}

func TestNestedClarityValues(t *testing.T) {
	// Create a simple tuple
	tuple := clarity_value.TupleValue{
		clarity_value.MustClarityName("a"): clarity_value.NewClarityValue(clarity_value.IntValue(1)),
		clarity_value.MustClarityName("b"): clarity_value.NewClarityValue(clarity_value.BoolValue(true)),
	}

	expectedRepr := "(tuple (a 1) (b true))"
	if repr := tuple.ReprString(); repr != expectedRepr {
		t.Errorf("Expected tuple repr %s, got %s", expectedRepr, repr)
	}

	// Create a list with some values
	list := clarity_value.ListValue{
		clarity_value.NewClarityValue(clarity_value.IntValue(1)),
		clarity_value.NewClarityValue(clarity_value.IntValue(2)),
		clarity_value.NewClarityValue(clarity_value.IntValue(3)),
	}

	expectedListRepr := "(list 1 2 3)"
	if repr := list.ReprString(); repr != expectedListRepr {
		t.Errorf("Expected list repr %s, got %s", expectedListRepr, repr)
	}

	// Test optional some
	optSome := clarity_value.OptionalSomeValue{
		Value: clarity_value.NewClarityValue(clarity_value.IntValue(42)),
	}

	expectedOptRepr := "(some 42)"
	if repr := optSome.ReprString(); repr != expectedOptRepr {
		t.Errorf("Expected optional some repr %s, got %s", expectedOptRepr, repr)
	}

	// Test response ok
	respOk := clarity_value.ResponseOkValue{
		Value: clarity_value.NewClarityValue(clarity_value.IntValue(42)),
	}

	expectedRespOkRepr := "(ok 42)"
	if repr := respOk.ReprString(); repr != expectedRespOkRepr {
		t.Errorf("Expected response ok repr %s, got %s", expectedRespOkRepr, repr)
	}

	// Test response err
	respErr := clarity_value.ResponseErrValue{
		Value: clarity_value.NewClarityValue(clarity_value.IntValue(42)),
	}

	expectedRespErrRepr := "(err 42)"
	if repr := respErr.ReprString(); repr != expectedRespErrRepr {
		t.Errorf("Expected response err repr %s, got %s", expectedRespErrRepr, repr)
	}
}

func TestStringUTF8Value(t *testing.T) {
	// Test a simple string
	simpleStr := clarity_value.NewStringUTF8Value([]byte("hello"))
	expectedRepr := "u\"hello\""
	if repr := simpleStr.ReprString(); repr != expectedRepr {
		t.Errorf("Expected UTF8 string repr %s, got %s", expectedRepr, repr)
	}

	// Test a string with non-ASCII characters
	unicodeStr := clarity_value.NewStringUTF8Value([]byte("hello世界"))
	// The exact representation will depend on how Go represents UTF-8 bytes,
	// but we can at least check that it contains the basic expected components
	if repr := unicodeStr.ReprString(); !strings.Contains(repr, "u\"hello") {
		t.Errorf("UTF8 string repr missing expected prefix: %s", repr)
	}
}
