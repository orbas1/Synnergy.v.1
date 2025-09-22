package config

import (
	"path/filepath"
	"testing"
)

func TestLoadUsesDevDefaultWhenNoPathProvided(t *testing.T) {
	t.Parallel()
	cfg, err := Load("", WithDefaultPaths(filepath.Join("..", "..", "configs", "dev.yaml")))
	if err != nil {
		t.Fatalf("expected to load dev defaults, got error: %v", err)
	}
	if cfg.Environment != "development" {
		t.Fatalf("expected development environment for dev defaults, got %s", cfg.Environment)
	}
}
