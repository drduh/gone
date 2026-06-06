// Package version provides application version
// and build details.
package version

var (

	// ID identifies the application ("gone")
	ID = "gone"

	// Version identifies the application version,
	// by default using the build date ("2026.12.31")
	Version = "2026"

	// Arch identifies the architecture of the Go
	// toolchain binaries ("arm64", "x86_64", etc.)
	Arch = "unknown"

	// Go identifies the version of Go used to
	// build the application ("go1.26.3")
	Go = "unknown"

	// Commit identifies the git commit hash the
	// application was built at
	// ("d36145acf6fda3e1821511cd4d4b73d21dbc34a7")
	Commit = "unknown"

	// Path identifies the path the application
	// was built from ("/home/user/git/gone")
	Path = "unknown"

	// Time identifies the time the application
	// was built at ("2026-12-31T12:00:00")
	Time = "unknown"

	// Host identifies the hostname of the host
	// which built the application ("mac.local")
	Host = "unknown"

	// System identifies the operating system of
	// the Go toolchain ("darwin", "linux", etc.)
	System = "unknown"

	// User identifies the system user which
	// built the application ("user")
	User = "unknown"
)
