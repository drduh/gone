package version

import (
	"encoding/json"
	"fmt"
)

// Print outputs version data in JSON format.
func Print() {
	vers := Get()
	data, err := json.Marshal(vers)
	if err != nil {
		return
	}
	fmt.Println(string(data))
}
