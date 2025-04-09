package version

import (
	"encoding/json"
	"fmt"
)

// Pretty prints full version data
func Print() {
	vers := Get()
	printJSON(vers)
}

// Pretty prints JSON interface
func printJSON(d interface{}) {
	data, err := json.Marshal(d)
	if err != nil {
		return
	}
	fmt.Println(string(data))
}
