package scripts_test

import (
	"os/exec"
	"strings"
	"testing"
)

func testScriptHelp(t *testing.T, script string) {
	cmd := exec.Command("bash", script, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s --help failed: %v\n%s", script, err, out)
	}
	if !strings.Contains(string(out), "Usage") {
		t.Fatalf("expected usage output from %s, got: %s", script, out)
	}
}

func TestPackageReleaseHelp(t *testing.T) {
	testScriptHelp(t, "package_release.sh")
}

func TestGenerateDocsHelp(t *testing.T) {
	testScriptHelp(t, "generate_docs.sh")
}

func TestCIScriptHelp(t *testing.T) {
	testScriptHelp(t, "ci_setup.sh")
}

func TestBackupLedgerHelp(t *testing.T) {
	testScriptHelp(t, "backup_ledger.sh")
}
