// Package throttle implements global
// rate-limiting using request times.
package throttle

import (
	"sync"
	"time"
)

// Requests represents time records of requests.
type Requests struct {

	// Mutex lock
	Lease sync.Mutex

	// Request time records
	RequestTimes []time.Time
}
