package util

import "os"

// GetHostname returns the OS hostname, or "unknown".
func GetHostname() string {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown"
	}
	return h
}
