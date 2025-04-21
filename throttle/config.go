// Package throttle implements request rate-limiting
// for the application.
package throttle

import (
	"sync"
	"time"
)

// Throttle requests by time
type Throttle struct {

	// Record times to rate limit
	Times []time.Time

	// Throttle lock
	Lease sync.Mutex
}
