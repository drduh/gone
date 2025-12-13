// Package version provides application version and
// build system information.
package version

var (

	// Application identifier ("gone")
	Id = "gone"

	// Date-based string ("2025.12.31")
	Version = "2025"

	// CPU architecture ("arm64", "x86_64", etc.)
	Arch = "unknown"

	// Go version ("go1.25.5")
	Go = "unknown"

	// Current git commit short hash ("a08f572")
	Commit = "unknown"

	// Build path ("/home/user/git/gone")
	Path = "unknown"

	// Build time ("2025-12-31T12:00:00")
	Time = "unknown"

	// Build hostname ("system")
	Host = "unknown"

	// Build system ("darwin")
	OS   = "unknown"

	// Build user ("user")
	User = "unknown"
)
