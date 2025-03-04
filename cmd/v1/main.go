package v1

import (
	"fmt"

	"github.com/drduh/gone/version"
)

func Run() {
	fmt.Println(version.Id, version.Version)
}
