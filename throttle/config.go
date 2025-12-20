// Package throttle implements global
// rate-limiting using request times.
package throttle

import (
	"sync"
	"time"
)

// Throttle represents request time records.
type Throttle struct {

	// Throttle lock
	Lease sync.Mutex

	// Recorded request times
	RequestTimes []time.Time
}
