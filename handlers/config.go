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

const (
	errorFileTooLarge = "file too large"
	errorRateLimit    = "limit exceeded"
	errorTmplExec     = "failed to exec"
	errorTmplParse    = "failed parsing"
)

var (
	responseErrorFileTooLarge = map[string]string{
		"error": errorFileTooLarge,
	}
	responseErrorRateLimit = map[string]string{
		"error": errorRateLimit,
	}
	responseErrorTmplExec = map[string]string{
		"error": errorTmplExec,
	}
	responseErrorTmplParse = map[string]string{
		"error": errorTmplParse,
	}
)
