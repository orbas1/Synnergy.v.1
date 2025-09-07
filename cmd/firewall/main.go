package main

import (
	"log"
	"os"

	synn "synnergy"
	"synnergy/cli"
)

// main executes the firewall subcommands from the shared Synnergy CLI. It
// preloads the gas table so reported costs match runtime expectations.
func main() {
	synn.LoadGasTable()
	cmd := cli.RootCmd()
	cmd.SetArgs(append([]string{"firewall"}, os.Args[1:]...))
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
