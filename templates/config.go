// Package templates provides HTML source to render.
package templates

import "embed"

//go:embed data/*.tmpl
var All embed.FS
