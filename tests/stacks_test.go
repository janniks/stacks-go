// Package tests contains tests for the stacks-go library.
package tests

import (
	"testing"

	stacks "github.com/hirosystems/stacks-go"
)

func TestVersion(t *testing.T) {
	v := stacks.Version()
	if v == "" {
		t.Error("Version should not be empty")
	}
}

func TestHello(t *testing.T) {
	greeting := stacks.Hello()
	expected := "Hello, Stacks!"
	if greeting != expected {
		t.Errorf("Expected greeting to be '%s', but got '%s'", expected, greeting)
	}
}
