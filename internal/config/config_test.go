package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTempFile creates a temporary file with provided contents and returns its path.
func writeTempFile(t *testing.T, dir, pattern, content string) string {
	t.Helper()
	file := filepath.Join(dir, pattern)
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return file
}

func TestLoadYAML(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "config.yaml", `
log_level: debug
environment: production
server:
  host: "127.0.0.1"
  port: 9000
database:
  url: "https://db.example.com"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load YAML config: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Fatalf("expected server port 9000 got %d", cfg.Server.Port)
	}
	if cfg.Environment != "production" {
		t.Fatalf("expected environment production got %s", cfg.Environment)
	}
}

func TestEnvOverrideAndJSON(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "config.json", `{
  "log_level": "info",
  "environment": "staging",
  "server": {"host": "0.0.0.0", "port": 8080},
  "database": {"url": "https://db.example.com"}
}`)

	t.Setenv("SYN_SERVER_PORT", "8081")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load JSON config: %v", err)
	}
	if cfg.Server.Port != 8081 {
		t.Fatalf("expected env override to 8081 got %d", cfg.Server.Port)
	}
}

func TestValidation(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.yaml", `log_level: info`)

	if _, err := Load(path); err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}
