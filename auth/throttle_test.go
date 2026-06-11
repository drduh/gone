package auth

import (
	"sync"
	"testing"
	"time"
)

// newThrottle sets up Throttle tests.
func newThrottle(times []time.Time) *RequestThrottle {
	return &RequestThrottle{
		Lease:        sync.Mutex{},
		RequestTimes: times,
	}
}

// TestAuthorizeDeny tests Throttle limit denial.
func TestAuthorizeDeny(t *testing.T) {
	now := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)

	r := newThrottle([]time.Time{
		now.Add(-10 * time.Second),
		now.Add(-20 * time.Second),
	})

	limit := 2
	if r.authorizeAt(limit, now) {
		t.Fatalf("expect deny with request limit %d", limit)
	}
}

// TestAuthorizeExpired tests expired Throttle times.
func TestAuthorizeExpired(t *testing.T) {
	now := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)

	r := newThrottle([]time.Time{
		now.Add(-1 * time.Minute),
		now.Add(-10 * time.Second),
	})

	limit := 2
	if !r.authorizeAt(limit, now) {
		t.Fatal("expect allow with only one recent request")
	}

	if got := len(r.RequestTimes); got != limit {
		t.Fatalf("len(RequestTimes) = %d, want %d",
			got, limit)
	}
}
