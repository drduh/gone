// Package config provides application settings and
// helper functions to enforce them.
package config

import "flag"

const (
	defaultMode = false
	defaultPath = ""
)

var (
	modeDebug   bool
	modeVersion bool

	pathConfig string
)

func init() {
	flag.StringVar(&pathConfig, "config", defaultPath,
		"Path to settings file")
	flag.StringVar(&pathConfig, "conf", defaultPath,
		"Shortcut for -config")
	flag.StringVar(&pathConfig, "c", defaultPath,
		"Shortcut for -config")
	flag.StringVar(&pathConfig, "settings", defaultPath,
		"Shortcut for -config")
	flag.StringVar(&pathConfig, "set", defaultPath,
		"Shortcut for -config")
	flag.StringVar(&pathConfig, "s", defaultPath,
		"Shortcut for -config")

	flag.BoolVar(&modeDebug, "debug", defaultMode,
		"Debug mode")
	flag.BoolVar(&modeDebug, "d", defaultMode,
		"Shortcut for -debug")
	flag.BoolVar(&modeDebug, "verbose", defaultMode,
		"Shortcut for -debug")

	flag.BoolVar(&modeVersion, "version", defaultMode,
		"Show application version")
	flag.BoolVar(&modeVersion, "vers", defaultMode,
		"Shortcut for -version")
	flag.BoolVar(&modeVersion, "v", defaultMode,
		"Shortcut for -version")
}
