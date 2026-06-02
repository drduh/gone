package auth

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const jitterFactor = 0.2

var tarpit atomic.Int64

func SetTarpit(d time.Duration) {
	tarpit.Store(int64(d))
}

func ApplyTarpit() {
	base := time.Duration(tarpit.Load())
	if base <= 0 {
		return
	}
	jitter := time.Duration(
		float64(base) * jitterFactor * rand.Float64())
	time.Sleep(base + jitter)
}
