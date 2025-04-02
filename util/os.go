package util

import "os"

// Returns hostname, if known
func GetHostname() string {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown"
	}
	return h
}
