// Package templates provides dynamic HTML templates to render.
package templates

import "embed"

//go:embed data/*.tmpl
var All embed.FS
