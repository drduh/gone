package templates

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
}
