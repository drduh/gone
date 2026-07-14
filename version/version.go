package version

// Visibility determines build detail redaction.
type Visibility int

const (

	// Full provides complete build detail.
	Full Visibility = iota

	// Redacted provides limited build detail.
	Redacted
)

// Get returns application version and build detail.
func Get(v Visibility) map[string]string {
	full := map[string]string{
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

	if v == Full {
		return full
	}

	return map[string]string{
		"appID": full["appID"],
	}
}
