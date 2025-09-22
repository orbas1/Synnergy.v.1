package scripts_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runCISetup(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command("bash", append([]string{"../../scripts/ci_setup.sh"}, args...)...)
	return cmd.CombinedOutput()
}

func TestCISetupHelp(t *testing.T) {
	out, err := runCISetup(t, "--help")
	if err != nil {
		t.Fatalf("expected help invocation to succeed: %v", err)
	}
	if !strings.Contains(string(out), "Usage: ci_setup.sh") {
		t.Fatalf("expected usage information, got: %s", out)
	}
}

func TestCISetupDryRunReport(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.json")
	cache := filepath.Join(dir, "cache")
	logFile := filepath.Join(dir, "ci.log")
	out, err := runCISetup(t,
		"--profile", "verify",
		"--project-root", dir,
		"--dry-run",
		"--report", report,
		"--cache-dir", cache,
		"--log-file", logFile,
		"--toolchain", "go",
	)
	if err != nil {
		t.Fatalf("expected dry-run to succeed: %v (%s)", err, out)
	}
	data, err := os.ReadFile(report)
	if err != nil {
		t.Fatalf("failed to read report: %v", err)
	}
	var payload struct {
		Tasks []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"tasks"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("invalid report json: %v", err)
	}
	if len(payload.Tasks) == 0 {
		t.Fatalf("expected at least one task in report")
	}
	for _, task := range payload.Tasks {
		if task.Status != "skipped" {
			t.Fatalf("expected dry-run tasks to be skipped, got: %#v", task)
		}
	}
	if _, err := os.Stat(cache); err == nil {
		t.Fatalf("cache directory should not be created during dry-run")
	}
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("expected log file to be created: %v", err)
	}
}

func TestCISetupUnknownProfile(t *testing.T) {
	out, err := runCISetup(t, "--profile", "mystery", "--dry-run")
	if err == nil {
		t.Fatalf("expected failure for unknown profile")
	}
	if !strings.Contains(string(out), "unknown profile") {
		t.Fatalf("expected unknown profile error, got: %s", out)
	}
}
