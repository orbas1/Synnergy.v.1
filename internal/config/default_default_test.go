//go:build !dev && !test && !prod

package config

import "testing"

func TestDefaultConfigPath(t *testing.T) {
	if DefaultConfigPath != "configs/dev.yaml" {
		t.Fatalf("expected default config path to be configs/dev.yaml, got %s", DefaultConfigPath)
	}
}
