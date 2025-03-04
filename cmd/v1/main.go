package v1

import (
	"fmt"
	"time"

	"github.com/drduh/gone/config"
)

func Run() {
	app := config.Load()

	fmt.Printf("%s on %s ran for %s\n",
		app.Version, app.Hostname,
		time.Since(app.Start).String())
}
