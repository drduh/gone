package handlers

var (
	errorDeny         = "not authorized"
	errorFileCopyFail = "file copy fail"
	errorFileFormFail = "file form fail"
	errorFileNotFound = "file not found"
	errorFileTooLarge = "file too large"
	errorRateLimit    = "limit exceeded"
	errorTmplExec     = "failed to exec"
	errorTmplParse    = "failed parsing"

	responseErrorDeny = map[string]string{
		"error": errorDeny,
	}
	responseErrorFileCopyFail = map[string]string{
		"error": errorFileCopyFail,
	}
	responseErrorFileFormFail = map[string]string{
		"error": errorFileFormFail,
	}
	responseErrorFileNotFound = map[string]string{
		"error": errorFileNotFound,
	}
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
