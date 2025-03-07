package auth

import (
	"net/http"
	"testing"
)

// Tests basic authentication
func TestBasic(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		header   string
		expected bool
	}{
		{
			name:     "Token not configured: always pass",
			token:    "",
			header:   "any",
			expected: true,
		},
		{
			name:     "Token matches header: pass",
			token:    "valid",
			header:   "valid",
			expected: true,
		},
		{
			name:     "Token does not match header: fail",
			token:    "invalid",
			header:   "valid",
			expected: false,
		},
		{
			name:     "No token header: fail",
			token:    "valid",
			header:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{
					HEADER: {tt.header},
				},
			}
			if got := Basic(tt.token, req); got != tt.expected {
				t.Errorf("Basic: %v; want %v", got, tt.expected)
			}
		})
	}
}
