package version

import (
	"strings"
)

// Returns short version as string
func Short() string {
	return strings.Join([]string{Id, Version}, "-")
}

// Returns full version as string map
func Full() map[string]string {
	return map[string]string{
		"id":   Id,
		"vers": Version,
		"user": User,
		"os":   OS,
		"arch": Arch,
		"go":   Go,
	}
}
