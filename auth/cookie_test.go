package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	cookieId    = "test-cookie"
	cookieTime  = time.Hour
	cookieValue = "cookie-value"
)

// TestCookieExists tests GetCookie reads a cookie when it exists.
func TestCookieExists(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieId, Value: cookieValue})
	rr := httptest.NewRecorder()

	result := GetCookie(rr, req, "default-value", cookieId, time.Hour)
	if result != cookieValue {
		t.Errorf("Expected %q, got %q", cookieValue, result)
	}
}

// TestCookieNotExists tests GetCookie sets and reads a cookie
// when it does not exist.
func TestCookieNotExists(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	result := GetCookie(rr, req, cookieValue, cookieId, time.Hour)
	if result != cookieValue {
		t.Errorf("Expected %q, got %q", cookieValue, result)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 ||
		cookies[0].Name != cookieId ||
		cookies[0].Value != cookieValue {
		t.Errorf("Cookie was not set correctly")
	}
}

// TestNewCookie tests NewCookie correctly sets a cookie.
func TestNewCookie(t *testing.T) {
	cookie := NewCookie(cookieValue, cookieId, cookieTime)
	if cookie.Path != "/" {
		t.Errorf("Expected %q, got %q", "/", cookie.Path)
	}
	if cookie.Name != cookieId {
		t.Errorf("Expected %q, got %q", cookieId, cookie.Name)
	}
	if cookie.Value != cookieValue {
		t.Errorf("Expected %q, got %q", cookieValue, cookie.Value)
	}
	if cookie.Expires.Before(
		time.Now().Add(cookieTime - time.Second)) {
		t.Errorf("Cookie expiration is invalid")
	}
}
