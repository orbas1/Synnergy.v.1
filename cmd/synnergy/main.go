package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"

	synn "synnergy"
	"synnergy/cli"
	"synnergy/internal/config"
)

func main() {
	// Load variables from .env if present to mirror the setup guides.
	gotenv.Load()

	cfgPath := os.Getenv("SYN_CONFIG")
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}

	// Configure logging based on loaded configuration.
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("invalid log level: %v", err)
	}
	logrus.SetLevel(lvl)

	// Warm up caches for shared resources.
	synn.LoadGasTable()
	logrus.Debug("gas table loaded")

	logrus.Infof("starting Synnergy in %s mode on %s:%d", cfg.Environment, cfg.Server.Host, cfg.Server.Port)

	if err := cli.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
