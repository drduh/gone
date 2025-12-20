// Package auth provides authentication and authorization
// functionality for the application.
package auth

import (
	"sync"
	"time"
)

// RequestThrottle represents the request rate-limiter.
type RequestThrottle struct {

	// Mutex lock
	Lease sync.Mutex

	// Request time records
	RequestTimes []time.Time
}
