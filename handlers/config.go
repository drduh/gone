package handlers

// Server request to log
type Request struct {

	// Handler path
	Action string `json:"action"`

	// User IP
	Address string `json:"address"`

	// User agent
	Agent string `json:"agent"`
}
