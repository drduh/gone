package version

// Get returns the application version and build information
func Get() map[string]string {
	return map[string]string{
		"id":   Id,
		"vers": Version,
		"user": User,
		"host": Host,
		"os":   OS,
		"arch": Arch,
		"go":   Go,
		"path": Path,
		"time": Time,
	}
}
