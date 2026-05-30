package settings

import "testing"

// TestMalformedJSON tests loading invalid JSON fails.
func TestMalformedJSON(t *testing.T) {
	var s Settings
	in := []byte(`{"port": 8080`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for malformed data, got nil")
	}
}

// TestTrailingData tests loading extra data fails.
func TestTrailingData(t *testing.T) {
	var s Settings
	in := []byte(`{"port": 8080} {"port": 9090}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for trailing data, got nil")
	}
}

// TestWrongTypePort tests loading port as string fails.
func TestWrongTypePort(t *testing.T) {
	var s Settings
	in := []byte(`{"port":"8080"}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for wrong type, got nil")
	}
}

// TestTypoField tests loading invalid fields fails.
func TestTypoField(t *testing.T) {
	var s Settings
	in := []byte(`{"prot": 8080}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for unknown field, got nil")
	}
}

// TestTypoNestedField tests loading invalid nested fields fails.
func TestTypoNestedField(t *testing.T) {
	var s Settings
	in := []byte(`{"audit":{"logFil":"gone.log"}}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for unknown nested field, got nil")
	}
}

// TestMergeSettings tests loading settings to override defaults.
func TestMergeSettings(t *testing.T) {
	var s Settings
	if err := loadSettings(defaultSettings, &s); err != nil {
		t.Fatalf("error loading default settings: %v", err)
	}
	if err := loadSettings([]byte(`{"showBuild": false}`), &s); err != nil {
		t.Fatalf("error loading override settings: %v", err)
	}
	if s.ShowBuild != false {
		t.Fatalf("expected showBuild:false, got %v", s.ShowBuild)
	}
}
