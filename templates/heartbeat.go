package templates

import "github.com/drduh/gone/config"

// Server status response
type Heartbeat struct {

	// Application version
	Version map[string]string `json:"version"`

	// Server hostname
	Hostname string `json:"hostname"`

	// Time since start ("3m45s")
	Uptime string `json:"uptime"`

	// TCP port server is listening on
	Port int `json:"port"`

	// Number of Files in storage
	FileCount int `json:"files"`

	// Number of Messages in storage
	MessageCount int `json:"messages"`

	// Limits configuration
	config.Limits `json:"limits"`

	// File owner information
	config.Owner `json:"owner"`
}
