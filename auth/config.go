// Package auth provides authentication and authorization
// functionality for the application.
package auth

import (
	"sync"
	"time"
)

// RequestThrottle represents the per-client request rate-limiter.
type RequestThrottle struct {

	// Mutex lock
	Lease sync.Mutex

	// Request time records by client address
	Clients map[string][]time.Time

	// Time of last full prune of idle clients
	swept time.Time
}
