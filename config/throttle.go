package config

import (
	"sync"
	"time"
)

// Throttle requests by time
type Throttle struct {

	// Record times to rate limit
	Times []time.Time

	// File lock
	Lease sync.Mutex
}
