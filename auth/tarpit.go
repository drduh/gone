package auth

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const jitterFactor = 0.2

var tarpit atomic.Int64

// SetTarpit sets the base delay.
func SetTarpit(d time.Duration) {
	tarpit.Store(int64(d))
}

// ApplyTarpit sleeps for the configured
// tarpit delay plus a jitter duration.
func ApplyTarpit() {
	base := time.Duration(tarpit.Load())
	if base <= 0 {
		return
	}
	j := rand.Float64() // #nosec G404
	jitter := time.Duration(
		float64(base) * jitterFactor * j)
	time.Sleep(base + jitter)
}
