package version

import (
	"encoding/json"
	"fmt"
)

// Pretty prints full version data
func Print() {
	vers := Full()
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
