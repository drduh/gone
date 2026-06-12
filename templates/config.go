// Package templates provides HTML source to render.
package templates

import "embed"

// TemplatesData contains embedded template files.
//
//go:embed data/*.tmpl
var TemplatesData embed.FS
