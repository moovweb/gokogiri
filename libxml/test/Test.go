package test

import "testing"

func AssertNil(t *testing.T, value interface{}, what string) {
	if value != nil {
		t.Error(what, "should be nil")
	}
}
func Equal(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Error("Expected: ", expected, "\nBut got: ", value)
	}
}
func Assert(t *testing.T, value interface{}, what string) interface{} {
	if value == nil {
		t.Error("Assertion failed: ", what)
	}
	return value
}
