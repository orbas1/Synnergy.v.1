package main

import (
	"log"
	"os"

	"synnergy/cli"
	"synnergy/internal/config"
)

func main() {
	cfg, err := config.Load(os.Getenv("SYN_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
