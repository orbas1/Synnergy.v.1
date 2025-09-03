package cli

import (
	"encoding/json"
	"fmt"
)

// jsonOutput toggles JSON formatted responses for CLI commands. Stage 25
// requires all node operations to support machine readable output for UI
// integration.
var jsonOutput bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output results in JSON")
}

// printOutput writes v to stdout respecting the global --json flag. Non-string
// types are pretty printed to aid scripts and dashboards.
func printOutput(v interface{}) {
	if jsonOutput {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
		return
	}
	switch val := v.(type) {
	case string:
		fmt.Println(val)
	default:
		fmt.Printf("%+v\n", val)
	}
}
