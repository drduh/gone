package version

// Get returns application version and build details.
func Get() map[string]string {
	return map[string]string{
		"id":      Id,
		"version": Version,
		"arch":    Arch,
		"go":      Go,
		"commit":  Commit,
		"path":    Path,
		"time":    Time,
		"host":    Host,
		"system":  System,
		"user":    User,
	}
}
