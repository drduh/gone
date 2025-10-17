// Package settings provides user-supplied and default
// application settings.
package settings

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

	// Content sharing defaults
	Default `json:"default"`

	// Content limits
	Limit `json:"limit"`

	// Paths to route
	Paths `json:"paths"`

	// TCP port to listen on
	Port int `json:"port"`

	// Whether to replace user IP address
	UserMask bool `json:"userMask"`
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

		// Whether to require authentication to edit shared content
		Wall bool `json:"wall"`
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

	// Upload form not valid
	Form string `json:"form"`

	// Filename not provided
	NoFilename string `json:"noFilename"`

	// No files available
	NoFiles string `json:"noFiles"`

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

	// Whether to display build, uptime and version in footer
	ShowBuild bool `json:"showBuild"`

	// Whether to enable Content Security Policy (CSP)
	CSP bool `json:"csp"`

	// Page title ("gone")
	Title string `json:"title"`

	// Cookie management
	Cookie struct {

		// Label ("goneTheme")
		Id string `json:"id"`

		// Time cookie is valid for ("192h")
		Time Duration `json:"time"`
	} `json:"cookie"`

	// CSS style options
	Style struct {

		// Whether to allow theme selection
		AllowPick bool `json:"allowPick"`

		// List of available themes to choose from, if allowed
		Available []string `json:"available"`

		// Theme to style with ("auto" for time-based option)
		Theme string `json:"theme"`
	} `json:"style"`

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

// Content sharing defaults
type Default struct {

	// Number of allowed downloads before File expiration
	Downloads int `json:"downloads,omitempty"`

	// Period of time before removing Files
	Expiration Duration `json:"duration,omitempty"`
}

// Content sharing limits
type Limit struct {

	// Message character length
	CharsMsg int `json:"charsMsg,omitempty"`

	// Wall character length
	CharsWall int `json:"charsWall,omitempty"`

	// Maximum file size (in Megabytes)
	MaxSizeMb int64 `json:"maxSizeMb,omitempty"`

	// Number of requests per minute to rate limit
	PerMinute int `json:"perMinute,omitempty"`

	// Frequency of File expiration check
	Ticker Duration `json:"ticker,omitempty"`
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

	// Random output ("/random/")
	Random string `json:"random"`

	// Default root path ("/")
	Root string `json:"root"`

	// Embedded/static file ("/static")
	Static string `json:"static"`

	// File upload ("/upload")
	Upload string `json:"upload"`

	// Shared content for edit ("/wall")
	Wall string `json:"wall"`
}
