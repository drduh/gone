package auth

import (
	"net/http"
	"net/url"
	"testing"
)

// TestBasic tests authentication with header and form values.
func TestBasic(t *testing.T) {
	orig := Tarpit
	t.Cleanup(func() { SetTarpit(orig) })
	SetTarpit(0)
	headerName := "X-Auth"
	tests := []struct {
		testName    string
		tokenValue  string
		headerValue string
		formValue   string
	}{
		{"noAuth", "", "", ""},
		{"headerPass", "valid-token", "valid-token", ""},
		{"headerFail", "valid-token", "wrong-token", ""},
		{"formPass", "valid-token", "", "valid-token"},
		{"formFail", "valid-token", "", "wrong-token"},
		{"allPass", "valid-token", "valid-token", "valid-token"},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{headerName: {tt.headerValue}},
				Form:   url.Values{headerName: {tt.formValue}},
			}
			req.PostForm = req.Form
			allowed := false
			if tt.tokenValue == "" {
				allowed = true
			} else if tt.headerValue == tt.tokenValue {
				allowed = true
			} else if tt.formValue == tt.tokenValue {
				allowed = true
			}
			ret := Basic(headerName, tt.tokenValue, req)
			if ret != allowed {
				t.Errorf("expected %v, got %v", allowed, ret)
			}
		})
	}
}
