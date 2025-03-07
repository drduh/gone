package handlers

var (
	errorDeny         = "not authorized"
	errorFileNotFound = "file not found"

	responseErrorDeny = map[string]string{
		"error": errorDeny,
	}
	responseErrorFileNotFound = map[string]string{
		"error": errorFileNotFound,
	}
)
