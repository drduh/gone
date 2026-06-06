package version

// Get returns application version and build details.
func Get() map[string]string {
	return map[string]string{
		"appID":         ID,
		"appVersion":    Version,
		"buildArch":     Arch,
		"buildCommit":   Commit,
		"buildGoVers":   Go,
		"buildHostname": Host,
		"buildPath":     Path,
		"buildSystem":   System,
		"buildTime":     Time,
		"buildUser":     User,
	}
}
