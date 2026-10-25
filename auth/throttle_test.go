package auth

import (
	"testing"
	"time"
)

const testAddr = "127.0.0.1"

// newThrottle sets up Throttle tests with times for testAddr.
func newThrottle(times []time.Time) *RequestThrottle {
	return &RequestThrottle{
		Clients: map[string][]time.Time{testAddr: times},
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
	if r.isAuthorizedAt(testAddr, limit, now) {
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
	if !r.isAuthorizedAt(testAddr, limit, now) {
		t.Fatal("expect allow with only one recent request")
	}

	if got := len(r.Clients[testAddr]); got != limit {
		t.Fatalf("len(Clients[addr]) = %d, want %d",
			got, limit)
	}
}

// TestAuthorizeUnlimited tests unlimited Throttle.
func TestAuthorizeUnlimited(t *testing.T) {
	now := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)

	r := newThrottle([]time.Time{
		now.Add(-10 * time.Second),
		now.Add(-20 * time.Second),
		now.Add(-30 * time.Second),
	})

	if !r.isAuthorizedAt(testAddr, 0, now) {
		t.Fatal("expect allow with no limit set")
	}
}

// TestAuthorizePerClient tests one client at its limit
// does not block another.
func TestAuthorizePerClient(t *testing.T) {
	now := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)

	r := newThrottle([]time.Time{
		now.Add(-10 * time.Second),
		now.Add(-20 * time.Second),
	})

	limit := 2
	if r.isAuthorizedAt(testAddr, limit, now) {
		t.Fatal("expect deny for client at limit")
	}
	if !r.isAuthorizedAt("10.0.0.2", limit, now) {
		t.Fatal("expect allow for other client")
	}
}

// TestAuthorizeSweep tests idle clients are pruned.
func TestAuthorizeSweep(t *testing.T) {
	now := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)

	r := newThrottle([]time.Time{now.Add(-2 * time.Minute)})
	r.Clients["10.0.0.2"] = []time.Time{now.Add(-3 * time.Minute)}
	r.swept = now.Add(-2 * time.Minute)

	r.isAuthorizedAt("10.0.0.3", 5, now)

	if _, ok := r.Clients["10.0.0.2"]; ok {
		t.Fatal("expect idle client pruned")
	}
	if got := len(r.Clients); got != 1 {
		t.Fatalf("len(Clients) = %d, want 1", got)
	}
}
