package scripts_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestDeployContractNoArgs ensures the deploy_contract.sh script prints usage when no arguments are provided.
func TestDeployContractNoArgs(t *testing.T) {
	cmd := exec.Command("bash", "../../scripts/deploy_contract.sh")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error when running without args")
	}
	if !strings.Contains(string(out), "Usage:") {
		t.Fatalf("expected usage output, got: %s", out)
	}
}

// TestDeployContractMissingFile verifies script error for nonexistent contract file.
func TestDeployContractMissingFile(t *testing.T) {
	cmd := exec.Command("bash", "../../scripts/deploy_contract.sh", "nonexistent.wasm")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
	if !strings.Contains(string(out), "contract file not found") {
		t.Fatalf("expected missing file message, got: %s", out)
	}
}

// TestDeployContractMissingBinary verifies error when CLI binary is absent.
func TestDeployContractMissingBinary(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "contract-*.wasm")
	if err != nil {
		t.Fatalf("failed to create temp contract: %v", err)
	}
	defer os.Remove(f.Name())

	cmd := exec.Command("bash", "../../scripts/deploy_contract.sh", f.Name())
	cmd.Env = append(os.Environ(), "BIN_PATH=/nonexistent/synnergy")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for missing binary")
	}
	if !strings.Contains(string(out), "binary not found") {
		t.Fatalf("expected binary not found message, got: %s", out)
	}
}
