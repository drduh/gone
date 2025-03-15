package config

// Application settings
type Settings struct {

	// Auditor options
	Audit `json:"audit"`

	// Authentication requirements
	Auth `json:"auth"`

	// Limit routes
	Limits `json:"limits"`

	// Paths to route
	Paths `json:"paths"`

	// TCP port to listen on
	Port int `json:"port"`
}

// Auditor logging preferences
type Audit struct {

	// Optional file destination for logs
	Filename string `json:"logFile"`

	// Format for datetime in logs
	TimeFormat string `json:"timeFormat"`
}

// Authentication requirements
type Auth struct {

	// Header key name ("X-Auth")
	Header string `json:"header"`

	// String-based token
	Token string `json:"token"`

	// Route authentication requirements
	Require struct {

		// Whether to require authentication for download
		Download bool `json:"download"`

		// Whether to require authentication for list
		List bool `json:"list"`

		// Whether to require authentication for upload
		Upload bool `json:"upload"`
	} `json:"require"`
}

// Download and upload limits
type Limits struct {

	// Number of allowed downloads before file expiration
	Downloads int `json:"downloads,omitempty"`

	// Maximum period of time to keep files
	Expiration Duration `json:"duration,omitempty"`

	// Frequency of file expiration check
	Ticker Duration `json:"ticker,omitempty"`

	// Maximum file size (in Megabytes)
	MaxSizeMb int `json:"maxSizeMb,omitempty"`

	// Number of requests per minute to rate limit
	PerMinute int `json:"perMinute,omitempty"`
}

// Paths to route
type Paths struct {

	// File download ("/download")
	Download string `json:"download"`

	// Heartbeat/health check ("/heartbeat")
	Heartbeat string `json:"heartbeat"`

	// File list ("/list")
	List string `json:"list"`

	// Embedded/static file ("/static")
	Static string `json:"static"`

	// File upload ("/upload")
	Upload string `json:"upload"`
}
