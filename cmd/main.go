// Package cmd provides the primary command
// entrypoint into the application.
package main

import (
	"os"

	"github.com/drduh/gone/cmd/v1"
)

func main() {
	os.Exit(v1.Run())
}
