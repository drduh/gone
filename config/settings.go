package config

// Application settings
type Settings struct {

	// Auditor options
	Audit `json:"audit"`

	// Authentication requirements
	Auth `json:"auth"`

	// User interface options
	Index `json:"index"`

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

// Authentication properties
type Auth struct {

	// Header key name ("X-Auth")
	Header string `json:"header"`

	// Form field placeholder ("secret required")
	Holder string `json:"holder"`

	// String-based token
	Token string `json:"token"`

	// Route authentication requirements
	Require struct {

		// Whether to require authentication to dowload files
		Download bool `json:"download"`

		// Whether to require authentication to list files
		List bool `json:"list"`

		// Whether to require authentication to post messages
		Message bool `json:"message"`

		// Whether to require authentication to upload files
		Upload bool `json:"upload"`
	} `json:"require"`
}

// Index HTML index page properties
type Index struct {

	// Whether to use CSS stylesheet
	Style bool

	// Page title
	Title string
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

	// Message post and read ("/")
	Message string `json:"message"`

	// Embedded/static file ("/static")
	Static string `json:"static"`

	// File upload ("/upload")
	Upload string `json:"upload"`
}
