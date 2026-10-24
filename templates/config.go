// Package templates provides HTML source to render.
package templates

import (
	"embed"
	"html/template"
)

// TemplatesData contains embedded template files.
//
//go:embed data/*.tmpl
var TemplatesData embed.FS

// IndexTemplate is parsed at startup.
var IndexTemplate = template.Must(
	template.New("index").ParseFS(TemplatesData, "data/*.tmpl"))
