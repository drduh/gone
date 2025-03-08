package handlers

var (
	errorDeny         = "not authorized"
	errorFileNotFound = "file not found"
	errorFileTooLarge = "file too large"

	responseErrorDeny = map[string]string{
		"error": errorDeny,
	}
	responseErrorFileNotFound = map[string]string{
		"error": errorFileNotFound,
	}
	responseErrorFileTooLarge = map[string]string{
		"error": errorFileTooLarge,
	}
)
