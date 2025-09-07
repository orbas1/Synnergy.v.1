package main

import (
	"log"
	"os"

	synn "synnergy"
	"synnergy/cli"
)

// main delegates to the government CLI commands within the shared Synnergy
// root command. Gas table loading ensures consistent cost reporting.
func main() {
	synn.LoadGasTable()
	cmd := cli.RootCmd()
	cmd.SetArgs(append([]string{"government"}, os.Args[1:]...))
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
