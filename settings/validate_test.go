package settings

import "testing"

// TestInvalidPort tests loading invalid port value fails.
func TestInvalidPort(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{name: "port zero", in: []byte(`{"port":0}`)},
		{name: "port negative", in: []byte(`{"port":-1}`)},
		{name: "port too large", in: []byte(`{"port":90001}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Settings
			if err := loadSettings(tt.in, &s); err == nil {
				t.Fatalf("expected error for invalid port, got nil")
			}
		})
	}
}

// TestMissingTimeFormat tests loading without timeFormat fails.
func TestMissingTimeFormat(t *testing.T) {
	var s Settings
	in := []byte(`{"audit":{"timeFormat":""}}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for missing audit timeFormat, got nil")
	}
}

// TestMissingAuthToken tests loading empty basic auth token fails.
func TestMissingAuthToken(t *testing.T) {
	var s Settings
	in := []byte(`{"auth":{"basic":{"field":"X-Auth","token":""}}}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatalf("expected error for basic auth without token, got nil")
	}
}
