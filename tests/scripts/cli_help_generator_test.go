package scripts_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runHelpGenerator(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command("bash", append([]string{"../../scripts/cli_help_generator.sh"}, args...)...)
	return cmd.CombinedOutput()
}

func TestCliHelpGeneratorDryRunRequiresDefinition(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "cli.md")
	out, err := runHelpGenerator(t, "--output", output, "--dry-run")
	if err == nil {
		t.Fatalf("expected dry-run without definitions to fail")
	}
	if !strings.Contains(string(out), "no commands specified") {
		t.Fatalf("expected missing definition message, got: %s", out)
	}
}

func TestCliHelpGeneratorMarkdown(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "cli.md")
	metrics := filepath.Join(dir, "metrics.json")
	definition := filepath.Join(dir, "defs.json")
	entries := []map[string]string{{
		"command": "status",
		"help":    "Display network status",
	}, {
		"command": "ledger sync",
		"help":    "Synchronize ledger state",
	}}
	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("failed to marshal definitions: %v", err)
	}
	if err := os.WriteFile(definition, data, 0o600); err != nil {
		t.Fatalf("failed to write definition file: %v", err)
	}
	out, err := runHelpGenerator(t,
		"--output", output,
		"--definition-file", definition,
		"--format", "markdown",
		"--metrics-file", metrics,
		"--dry-run",
	)
	if err != nil {
		t.Fatalf("expected markdown generation to succeed: %v (%s)", err, out)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read generated output: %v", err)
	}
	if !strings.Contains(string(content), "## status") || !strings.Contains(string(content), "ledger sync") {
		t.Fatalf("expected command sections in output, got: %s", content)
	}
	metricsData, err := os.ReadFile(metrics)
	if err != nil {
		t.Fatalf("failed to read metrics: %v", err)
	}
	if !strings.Contains(string(metricsData), "\"command_count\": 2") {
		t.Fatalf("expected command count in metrics: %s", metricsData)
	}
}
