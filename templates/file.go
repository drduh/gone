package templates

// File information response
type File struct {

	// Name of the file
	Name string `json:"name"`

	// Size of file
	Size int `json:"size"`

	// File owner/uploader
	Owner `json:"owner"`
}
