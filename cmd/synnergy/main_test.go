package main

import (
	"testing"

	synn "synnergy"
	"synnergy/internal/config"
)

func TestRegisterEnterpriseGasMetadata(t *testing.T) {
	synn.ResetGasTable()
	if err := registerEnterpriseGasMetadata(); err != nil {
		t.Fatalf("registerEnterpriseGasMetadata: %v", err)
	}
	entry, ok := synn.GasMetadataFor("EnterpriseBootstrap")
	if !ok {
		t.Fatalf("expected metadata for EnterpriseBootstrap")
	}
	if entry.Category != "orchestrator" {
		t.Fatalf("unexpected category: %s", entry.Category)
	}
}

func TestConfigureLogging(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{
			Level:   "info",
			Format:  "json",
			Outputs: []string{"stdout"},
		},
	}
	if err := configureLogging(cfg); err != nil {
		t.Fatalf("configureLogging: %v", err)
	}
}
