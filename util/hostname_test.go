package util

import "testing"

// TestGetHostname tests a hostname is always returned.
func TestGetHostname(t *testing.T) {
	if got := GetHostname(); got == "" {
		t.Fatal("expected non-empty hostname")
	}
}
