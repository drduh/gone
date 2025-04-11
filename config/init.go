// Package config provides required application settings,
// and helper functions necessary to enforce them.
package config

import "flag"

var (
	modeDebug   bool
	modeVersion bool
	pathConfig  string
)

func init() {
	flag.BoolVar(&modeDebug, "debug", false, "Debug mode")
	flag.BoolVar(&modeDebug, "d", false, "Shortcut for -debug")

	flag.BoolVar(&modeVersion, "version", false, "Show version")
	flag.BoolVar(&modeVersion, "vers", false, "Shortcut for -version")
	flag.BoolVar(&modeVersion, "v", false, "Shortcut for -version")

	flag.StringVar(&pathConfig, "config", "", "Path to settings JSON")
	flag.StringVar(&pathConfig, "conf", "", "Shortcut for -config")
	flag.StringVar(&pathConfig, "c", "", "Shortcut for -config")

	flag.Parse()
}
