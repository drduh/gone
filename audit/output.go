package audit

import (
	"fmt"
	"io"
	"os"
)

// Returns file or stdout IO writer
func getDest(filename string) (io.Writer, error) {
	var err error
	dest := os.Stdout
	if filename != "" {
		dest, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s", filename)
		}
	}
	return dest, nil
}
