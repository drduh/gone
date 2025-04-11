package util

import "time"

// IsDaytime returns true if it is approximately daytime.
func IsDaytime() bool {
	now := time.Now().Hour()
	return now >= 7 && now < 19
}
