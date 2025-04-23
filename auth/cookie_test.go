package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCookieExists tests GetCookie reads a cookie when it exists.
func TestCookieExists(t *testing.T) {
	id := "test-cookie"
	expectedValue := "cookie-value"

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: id, Value: expectedValue})
	rr := httptest.NewRecorder()

	result := GetCookie(rr, req, "default-value", id, time.Hour)
	if result != expectedValue {
		t.Errorf("Expected %q, got %q", expectedValue, result)
	}
}

// TestCookieNotExists tests GetCookie sets and reads a cookie
// when it does not exist.
func TestCookieNotExists(t *testing.T) {
	id := "test-cookie"
	defaultValue := "default-value"

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	result := GetCookie(rr, req, defaultValue, id, time.Hour)
	if result != defaultValue {
		t.Errorf("Expected default value %q, got %q", defaultValue, result)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 ||
		cookies[0].Name != id ||
		cookies[0].Value != defaultValue {
		t.Errorf("Cookie was not set correctly")
	}
}

// TestNewCookie tests NewCookie correctly sets a cookie.
func TestNewCookie(t *testing.T) {
	id := "test-cookie"
	value := "cookie-value"
	duration := time.Hour

	cookie := NewCookie(value, id, duration)
	if cookie.Name != id {
		t.Errorf("Expected %q, got %q", id, cookie.Name)
	}
	if cookie.Value != value {
		t.Errorf("Expected %q, got %q", value, cookie.Value)
	}
	if cookie.Path != "/" {
		t.Errorf("Expected %q, got %q", "/", cookie.Path)
	}
	if cookie.Expires.Before(time.Now().Add(duration - time.Second)) {
		t.Errorf("Cookie expiration is invalid")
	}
}
