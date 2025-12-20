// Package settings provides user-supplied and default
// application settings.
package settings

// Settings represents all application settings.
type Settings struct {

	// Auditor options
	Audit `json:"audit,omitempty"`

	// Authentication requirements
	Auth `json:"auth,omitempty"`

	// Errors to serve
	Error `json:"error,omitempty"`

	// User interface options
	Index `json:"indexPage,omitempty"`

	// Content sharing defaults
	Default `json:"default,omitempty"`

	// Content limits
	Limit `json:"limit,omitempty"`

	// Paths to route
	Paths `json:"paths,omitempty"`

	// TCP port to listen on
	Port int `json:"port,omitempty"`

	// Show build details in status response
	// and index footer
	ShowBuild bool `json:"showBuild,omitempty"`
}

// Audit represents logging preferences.
type Audit struct {

	// Optional file destination for logs
	Filename string `json:"logFile,omitempty"`

	// Format for datetime in logs
	TimeFormat string `json:"timeFormat,omitempty"`
}

// Auth represents the authentication configuration.
type Auth struct {

	// Basic represents token-based authentication.
	Basic struct {

		// Header or form field name ("X-Auth")
		Field string `json:"field,omitempty"`

		// Authentication token string ("mySecret")
		Token string `json:"token,omitempty"`
	} `json:"basic,omitempty"`

	// Require represents paths which may require auth.
	Require struct {

		// Require authentication to clear uploads
		Clear bool `json:"clear,omitempty"`

		// Require authentication to download files
		Download bool `json:"download,omitempty"`

		// Require authentication to load root (index)
		Root bool `json:"root,omitempty"`

		// Require authentication to list files
		List bool `json:"list,omitempty"`

		// Require authentication to post messages
		Message bool `json:"message,omitempty"`

		// Require authentication to upload files
		Upload bool `json:"upload,omitempty"`

		// Require authentication to edit shared content
		Wall bool `json:"wall,omitempty"`
	} `json:"require,omitempty"`
}

// Error represents responses error message strings.
type Error struct {

	// Failed to copy file
	Copy string `json:"copy,omitempty"`

	// Deny (not authorized)
	Deny string `json:"deny,omitempty"`

	// File exceeds size limit
	FileSize string `json:"fileSize,omitempty"`

	// Upload form not valid
	Form string `json:"form,omitempty"`

	// Filename not provided
	NoFilename string `json:"noFilename,omitempty"`

	// No files available
	NoFiles string `json:"noFiles,omitempty"`

	// File not found in Storage
	NotFound string `json:"notFound,omitempty"`

	// Template could not be executed
	TmplExec string `json:"tmplExec,omitempty"`

	// Template could not be parsed
	TmplParse string `json:"tmplParse,omitempty"`

	// Too many requests
	RateLimit string `json:"rateLimit,omitempty"`
}

// Index represents index HTML page properties.
type Index struct {

	// Enable Content Security Policy (CSP)
	CSP bool `json:"csp,omitempty"`

	// Page title ("gone")
	Title string `json:"title,omitempty"`

	// Cookie represents cookie-based user preferences.
	Cookie struct {

		// Label ("goneTheme")
		Id string `json:"id,omitempty"`

		// Time cookie is valid for ("192h")
		Time Duration `json:"time,omitempty"`
	} `json:"cookie,omitempty"`

	// Style represents theme options.
	Style struct {

		// Allow theme selection
		AllowPick bool `json:"allowPick,omitempty"`

		// List of available themes to choose from, if allowed
		Available []string `json:"available,omitempty"`

		// Theme to style with ("auto" for time-based option)
		Theme string `json:"theme,omitempty"`
	} `json:"style,omitempty"`

	// Placeholder represents placeholder text in forms.
	Placeholder struct {

		// Authentication field
		Auth string `json:"auth,omitempty"`

		// File selection field
		File string `json:"file,omitempty"`

		// Message input field
		Message string `json:"message,omitempty"`
	} `json:"placeholder,omitempty"`
}

// Default represents content sharing default options.
type Default struct {

	// Number of allowed downloads before File expiration
	Downloads int `json:"downloads,omitempty"`

	// Period of time before removing Files
	Expiration Duration `json:"duration,omitempty"`
}

// Content represents limits on content sharing.
type Limit struct {

	// Maximum text message size
	MaxSizeMsg int `json:"maxSizeMsg,omitempty"`

	// Maximum wall content size
	MaxSizeWall int `json:"maxSizeWall,omitempty"`

	// Maximum file size (in Megabytes)
	MaxSizeFileMb int64 `json:"maxSizeFileMb,omitempty"`

	// Global requests per minute to rate limit
	ReqsPerMinute int `json:"reqsPerMinute,omitempty"`

	// Frequency of File expiration check
	Ticker Duration `json:"expiryTicker,omitempty"`
}

// Paths represents paths to handle.
type Paths struct {

	// HTML assets, such as stylesheet ("/assets/")
	Assets string `json:"assets,omitempty"`

	// Remove content ("/clear")
	Clear string `json:"clear,omitempty"`

	// File download ("/download/")
	Download string `json:"download,omitempty"`

	// File list ("/list")
	List string `json:"list,omitempty"`

	// Message read and write ("/msg")
	Message string `json:"message,omitempty"`

	// Random output ("/random/")
	Random string `json:"random,omitempty"`

	// Default root path ("/")
	Root string `json:"root,omitempty"`

	// Embedded static file ("/static")
	Static string `json:"static,omitempty"`

	// Status check ("/status")
	Status string `json:"status,omitempty"`

	// File upload ("/upload")
	Upload string `json:"upload,omitempty"`

	// Shared content read and write ("/wall")
	Wall string `json:"wall,omitempty"`
}
