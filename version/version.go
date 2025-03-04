package version

import (
	"strings"
)

// Returns short version as string
func Short() string {
	return strings.Join([]string{Id, Version}, "-")
}
