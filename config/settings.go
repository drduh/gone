package config

import "fmt"

// Application settings
type Settings struct {

	// Auditor options
	Audit `json:"audit"`

	// Authentication requirements
	Auth `json:"auth"`

	// Errors to serve
	Error `json:"error"`

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

	// String-based check
	Basic struct {

		// Header or form field name ("X-Auth")
		Field string `json:"field"`

		// String-based token
		Token string `json:"token"`
	} `json:"basic"`

	// Route authentication requirements
	Require struct {

		// Whether to require authentication to download files
		Download bool `json:"download"`

		// Whether to require authentication to list files
		List bool `json:"list"`

		// Whether to require authentication to post messages
		Message bool `json:"message"`

		// Whether to require authentication to upload files
		Upload bool `json:"upload"`
	} `json:"require"`
}

// Error response
type Error struct {

	// Failed to copy file
	Copy string `json:"copy"`

	// Deny (not authorized)
	Deny string `json:"deny"`

	// File exceeds size limit
	FileSize string `json:"fileSize"`

	// Upload form error
	Form string `json:"form"`

	// Filename not provided
	NoFilename string `json:"noFilename"`

	// File not found in Storage
	NotFound string `json:"notFound"`

	// Template could not be executed
	TmplExec string `json:"tmplExec"`

	// Template could not be parsed
	TmplParse string `json:"tmplParse"`

	// Too many requests
	RateLimit string `json:"rateLimit"`
}

// Index HTML index page properties
type Index struct {

	// Whether to enable Content Security Policy (CSP)
	CSP bool `json:"csp"`

	// Whether to allow theme selection
	ThemePick bool `json:"themePick"`

	// CSS theme to style with (leave empty for auto selection)
	Theme string `json:"theme"`

	// Page title ("gone")
	Title string `json:"title"`

	// Cookie management
	Cookie struct {

		// Label ("goneTheme")
		Id string `json:"id"`

		// Time cookie is valid for ("192h")
		Time Duration `json:"time"`
	} `json:"cookie"`

	// Index form placeholder text
	Placeholder struct {

		// Authentication field
		Auth string `json:"auth"`

		// File selection field
		Filename string `json:"filename"`

		// Message input field
		Message string `json:"message"`
	} `json:"placeholder"`
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
	MaxSizeMb int64 `json:"maxSizeMb,omitempty"`

	// Message character length
	MsgChars int `json:"msgChars,omitempty"`

	// Number of requests per minute to rate limit
	PerMinute int `json:"perMinute,omitempty"`
}

// Paths to route
type Paths struct {

	// Assets for HTML pages ("/assets/")
	Assets string `json:"assets"`

	// File download ("/download/")
	Download string `json:"download"`

	// Heartbeat/health check ("/heartbeat")
	Heartbeat string `json:"heartbeat"`

	// File list ("/list")
	List string `json:"list"`

	// Message post and read ("/msg")
	Message string `json:"message"`

	// Embedded/static file ("/static")
	Static string `json:"static"`

	// File upload ("/upload")
	Upload string `json:"upload"`
}

// Returns address string based on port
func (s *Settings) GetAddr() string {
	return fmt.Sprintf(":%d", s.Port)
}

// Returns Mb size to bytes
func (l *Limits) GetMaxBytes() int64 {
	return l.MaxSizeMb << 20
}
