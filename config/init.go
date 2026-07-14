// Package config provides application settings and
// helper functions to enforce them.
package config

import "flag"

const defaultMode = false

var (
	modeDebug       bool
	modeShowVersion bool

	pathConfig string
	authToken  string
)

func init() {
	flag.BoolVar(&modeShowVersion, "version", defaultMode,
		"Show version and build information")
	flag.BoolVar(&modeShowVersion, "v", defaultMode,
		"Shortcut for -version")

	flag.BoolVar(&modeDebug, "debug", defaultMode,
		"Debug mode")
	flag.BoolVar(&modeDebug, "d", defaultMode,
		"Shortcut for -debug")

	flag.StringVar(&pathConfig, "config", "",
		"Path to settings file")
	flag.StringVar(&pathConfig, "c", "",
		"Shortcut for -config")

	flag.StringVar(&authToken, "auth", "",
		"Request authentication token")
	flag.StringVar(&authToken, "a", "",
		"Shortcut for -auth")
}
