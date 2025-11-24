package auth

import "time"

var Tarpit = 2 * time.Second

// SetTarpit configures the tarpit delay.
func SetTarpit(d time.Duration) {
	Tarpit = d
}

// applyTarpit applies the configured delay.
func applyTarpit() {
	if Tarpit <= 0 {
		return
	}
	time.Sleep(Tarpit)
}
