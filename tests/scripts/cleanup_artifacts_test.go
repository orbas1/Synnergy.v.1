package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func runCleanup(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command("bash", append([]string{"../../scripts/cleanup_artifacts.sh"}, args...)...)
	return cmd.CombinedOutput()
}

func TestCleanupArtifactsDryRun(t *testing.T) {
	workspace := t.TempDir()
	buildDir := filepath.Join(workspace, "build")
	if err := os.Mkdir(buildDir, 0o755); err != nil {
		t.Fatalf("failed to create build dir: %v", err)
	}
	artifact := filepath.Join(buildDir, "artifact.bin")
	if err := os.WriteFile(artifact, []byte("payload"), 0o600); err != nil {
		t.Fatalf("failed to write artifact: %v", err)
	}
	metrics := filepath.Join(workspace, "metrics.json")
	out, err := runCleanup(t,
		"--workspace", workspace,
		"--target", "build",
		"--dry-run",
		"--metrics-file", metrics,
	)
	if err != nil {
		t.Fatalf("expected dry-run to succeed: %v (%s)", err, out)
	}
	if _, err := os.Stat(artifact); err != nil {
		t.Fatalf("artifact should remain after dry-run: %v", err)
	}
	data, err := os.ReadFile(metrics)
	if err != nil {
		t.Fatalf("failed to read metrics: %v", err)
	}
	if !strings.Contains(string(data), "\"dry_run\": true") {
		t.Fatalf("expected dry-run flag in metrics: %s", data)
	}
}

func TestCleanupArtifactsDeletion(t *testing.T) {
	workspace := t.TempDir()
	logsDir := filepath.Join(workspace, "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}
	logFile := filepath.Join(logsDir, "old.log")
	if err := os.WriteFile(logFile, []byte("hello"), 0o600); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}
	// Ensure the file is older than the threshold we will use
	past := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(logFile, past, past); err != nil {
		t.Fatalf("failed to update timestamps: %v", err)
	}

	out, err := runCleanup(t,
		"--workspace", workspace,
		"--target", "logs",
		"--max-age", "1",
	)
	if err != nil {
		t.Fatalf("expected cleanup to succeed: %v (%s)", err, out)
	}
	if _, err := os.Stat(logFile); !os.IsNotExist(err) {
		t.Fatalf("expected log file to be removed, got err: %v", err)
	}
}
