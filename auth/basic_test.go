package auth

import (
	"net/http"
	"testing"
)

// Tests basic authentication
func TestBasic(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		allowed  string
		token    string
		expected bool
	}{
		{
			name:     "Token not required",
			header:   "X-Auth",
			allowed:  "",
			token:    "",
			expected: true,
		},
		{
			name:     "Required token matches",
			header:   "X-Auth",
			allowed:  "valid-token",
			token:    "valid-token",
			expected: true,
		},
		{
			name:     "Required token does not match",
			header:   "X-Auth",
			allowed:  "valid-token",
			token:    "invalid-token",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{
					tt.header: {tt.token},
				},
			}
			if got := Basic(tt.header, tt.allowed, req); got != tt.expected {
				t.Errorf("%s: got %v (with '%#v'), expected %v",
					tt.name, got, req.Header, tt.expected)
			}
		})
	}
}
