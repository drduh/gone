package auth

import "time"

var Tarpit = 2 * time.Second

// SetTarpit configures the tarpit delay programmatically.
func SetTarpit(d time.Duration) {
	Tarpit = d
}

func applyTarpit() {
	if Tarpit <= 0 {
		return
	}
	time.Sleep(Tarpit)
}
