package templates

import (
	"net/http"

	"github.com/drduh/gone/config"
)

// Server status response
type Heartbeat struct {

	// Application version
	Version string `json:"version"`

	// Server hostname
	Hostname string `json:"hostname"`

	// Time since start ("3m45s")
	Uptime string `json:"uptime"`

	// TCP port server is listening on
	Port int `json:"port"`

	// Number of files in storage
	FileCount int `json:"files"`

	// Limits configuration
	config.Limits `json:"limits"`

	// Client information
	Client `json:"client"`
}

// Client information
type Client struct {

	// Remote IP address
	Address string `json:"address"`

	// HTTP headers
	Headers http.Header `json:"headers"`
}
