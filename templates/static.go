package templates

import _ "embed"

//go:embed data/static.txt
var StaticData string

// Static data response
type Static struct {

	// Content data
	Data string `json:"data"`
}
