package version

// Returns full version as string map
func Full() map[string]string {
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
