package templates

import _ "embed"

// StaticData contains embedded files data.
//
//go:embed data/static.txt
var StaticData string

// Static represents static embedded data.
type Static struct {

	// Content data
	Data string `json:"data"`
}
