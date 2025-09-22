package config

import (
	"path/filepath"
	"testing"
)

func TestLoadUsesProdProfileWhenProvided(t *testing.T) {
	t.Parallel()
	cfg, err := Load(filepath.Join("..", "..", "configs", "prod.yaml"))
	if err != nil {
		t.Fatalf("expected prod config to load, got error: %v", err)
	}
	if cfg.Environment != "production" {
		t.Fatalf("expected production environment from prod profile, got %s", cfg.Environment)
	}
}
