package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	cookieID    = "test-cookie"
	cookieMode  = http.SameSiteStrictMode
	cookieTime  = time.Hour
	cookieValue = "cookie-value"
)

// TestCookieExists tests GetCookie reads a cookie when it exists.
func TestCookieExists(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:     cookieID,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Value:    cookieValue},
	)

	rr := httptest.NewRecorder()

	result := GetCookie(rr, req, "default-value", cookieID, time.Hour)
	if result != cookieValue {
		t.Errorf("Expected %q, got %q", cookieValue, result)
	}
}

// TestCookieNotExists tests GetCookie sets and reads a cookie
// when it does not exist.
func TestCookieNotExists(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(),
		http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()

	result := GetCookie(rr, req, cookieValue, cookieID, time.Hour)
	if result != cookieValue {
		t.Errorf("Expected %q, got %q", cookieValue, result)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 ||
		cookies[0].Name != cookieID ||
		cookies[0].Value != cookieValue {
		t.Errorf("Cookie was not set correctly")
	}
	if !cookies[0].HttpOnly {
		t.Errorf("Expected HttpOnly=true, got false")
	}
	if cookies[0].SameSite != cookieMode {
		t.Errorf("Expected SameSite=%d, got %d", cookieMode, cookies[0].SameSite)
	}
	if !cookies[0].Secure {
		t.Errorf("Expected Secure=true, got false")
	}
}

// TestNewCookie tests NewCookie correctly sets a cookie.
func TestNewCookie(t *testing.T) {
	cookie := NewCookie(cookieValue, cookieID, cookieTime)
	if cookie.Path != "/" {
		t.Errorf("Expected Path=%q, got %q", "/",
			cookie.Path)
	}
	if cookie.Name != cookieID {
		t.Errorf("Expected ID=%q, got %q",
			cookieID, cookie.Name)
	}
	if cookie.Value != cookieValue {
		t.Errorf("Expected Value=%q, got %q",
			cookieValue, cookie.Value)
	}
	if !cookie.HttpOnly {
		t.Errorf("Expected HttpOnly=true, got false")
	}
	if cookie.SameSite != cookieMode {
		t.Errorf("Expected SameSite=%d, got %d",
			cookieMode, cookie.SameSite)
	}
	if !cookie.Secure {
		t.Errorf("Expected Secure=true, got false")
	}
	if cookie.Expires.Before(
		time.Now().Add(cookieTime - time.Second)) {
		t.Errorf("Cookie expiration is invalid")
	}
}
