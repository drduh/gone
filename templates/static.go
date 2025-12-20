package templates

import _ "embed"

//go:embed data/static.txt
var StaticData string

// Static represents static embedded data.
type Static struct {

	// Content data
	Data string `json:"data"`
}
