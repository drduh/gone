package settings

import "testing"

// TestInvalidAddr tests loading invalid address value.
func TestInvalidAddr(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{name: "empty string",
			in: []byte(`{"serverAddr":"   "}`)},
		{name: "hostname",
			in: []byte(`{"serverAddr":"localhost"}`)},
		{name: "ipv4 out of valid range",
			in: []byte(`{"serverAddr":"123.456.789.1"}`)},
		{name: "ipv4 incomplete",
			in: []byte(`{"serverAddr":"127.0.0"}`)},
		{name: "ipv4 with port",
			in: []byte(`{"serverAddr":"127.0.0.1:8080"}`)},
		{name: "cidr notation",
			in: []byte(`{"serverAddr":"127.0.0.1/24"}`)},
		{name: "ipv6 with brackets",
			in: []byte(`{"serverAddr":"[::1]"}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Settings

			if err := loadSettings(defaultSettings, &s); err != nil {
				t.Fatalf("error loading default settings: %v", err)
			}

			err := loadSettings(tt.in, &s)
			if err == nil {
				t.Fatal("expected server address error")
			}
		})
	}
}

// TestInvalidPort tests loading invalid port value.
func TestInvalidPort(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{name: "port zero",
			in: []byte(`{"serverPort":0}`)},
		{name: "port negative",
			in: []byte(`{"serverPort":-1}`)},
		{name: "port out of valid range",
			in: []byte(`{"serverPort":90001}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Settings

			if err := loadSettings(defaultSettings, &s); err != nil {
				t.Fatalf("error loading default settings: %v", err)
			}

			err := loadSettings(tt.in, &s)
			if err == nil {
				t.Fatal("expected server port error")
			}
		})
	}
}

// TestMissingTimeFormat tests loading without a time format.
func TestMissingTimeFormat(t *testing.T) {
	var s Settings

	if err := loadSettings(defaultSettings, &s); err != nil {
		t.Fatalf("error loading default settings: %v", err)
	}

	in := []byte(`{"audit":{"timeFormat":""}}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatal("expected missing audit time format error")
	}
}

// TestMissingAuthToken tests loading empty auth token.
func TestMissingAuthToken(t *testing.T) {
	var s Settings

	if err := loadSettings(defaultSettings, &s); err != nil {
		t.Fatalf("error loading default settings: %v", err)
	}

	in := []byte(`{"auth":{"basic":{"field":"X-Auth","token":""}}}`)
	if err := loadSettings(in, &s); err == nil {
		t.Fatal("expected missing auth token error")
	}
}

// TestInvalidTarpitDelay tests loading invalid tarpit delay.
func TestInvalidTarpitDelay(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{name: "negative seconds",
			in: []byte(`{"auth":{"tarpitDelay":"-1s"}}`)},
		{name: "negative numeric shorthand",
			in: []byte(`{"auth":{"tarpitDelay":"invalid"}}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Settings

			if err := loadSettings(defaultSettings, &s); err != nil {
				t.Fatalf("error loading default settings: %v", err)
			}

			if err := loadSettings(tt.in, &s); err == nil {
				t.Fatal("expected tarpit delay error")
			}
		})
	}
}

// TestInvalidContentLimits tests loading invalid content limits.
func TestInvalidContentLimits(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
	}{
		{name: "file name length zero",
			in: []byte(`{"limit":{"fileLimits":{"nameLength":0}}}`)},
		{name: "file name length negative",
			in: []byte(`{"limit":{"fileLimits":{"nameLength":-1}}}`)},
		{name: "message length zero",
			in: []byte(`{"limit":{"messageLimits":{"lengthChars":0}}}`)},
		{name: "message length negative",
			in: []byte(`{"limit":{"messageLimits":{"lengthChars":-1}}}`)},
		{name: "wall length zero",
			in: []byte(`{"limit":{"wallLimits":{"lengthChars":0}}}`)},
		{name: "wall length negative",
			in: []byte(`{"limit":{"wallLimits":{"lengthChars":-1}}}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Settings

			if err := loadSettings(defaultSettings, &s); err != nil {
				t.Fatalf("error loading default settings: %v", err)
			}

			err := loadSettings(tt.in, &s)
			if err == nil {
				t.Fatal("expected invalid content limits error")
			}
		})
	}
}
