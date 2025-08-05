package main

import (
	"log"
	"os"

	"synnergy/cli"
	"synnergy/internal/config"
)

func main() {
	cfgPath := os.Getenv("SYN_CONFIG")
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
