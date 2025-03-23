package auth

import (
	"net/http"
	"testing"
)

// Tests basic authentication
func TestBasic(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		allowed string
		token   string
		expect  bool
	}{
		{
			name:    "Token not required",
			header:  "X-Auth",
			allowed: "",
			token:   "",
			expect:  true,
		},
		{
			name:    "Required token matches",
			header:  "X-Auth",
			allowed: "valid-token",
			token:   "valid-token",
			expect:  true,
		},
		{
			name:    "Required token does not match",
			header:  "X-Auth",
			allowed: "valid-token",
			token:   "invalid-token",
			expect:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{
					tt.header: {tt.token},
				},
			}
			if got := Basic(tt.header, tt.allowed, req); got != tt.expect {
				t.Errorf("%s=%v; expect %v", tt.name, got, tt.expect)
			}
		})
	}
}
