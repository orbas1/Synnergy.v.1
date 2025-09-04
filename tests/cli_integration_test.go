package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	synn "synnergy"
	"synnergy/cli"
	"synnergy/core"
)

// execCommand executes the Synnergy CLI command with the given arguments and
// returns the combined output from stdout and stderr.
func execCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := cli.RootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	cmd.SetArgs([]string{})
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	r.Close()
	return strings.TrimSpace(buf.String() + string(out)), err
}

// TestCLIIntegration exercises a minimal cross-section of CLI commands to
// verify that the root command and core modules are wired together properly.
func TestCLIIntegration(t *testing.T) {
	synn.LoadGasTable()

	t.Run("help", func(t *testing.T) {
		out, err := execCommand(t, "--help")
		if err != nil {
			t.Fatalf("help failed: %v", err)
		}
		if !strings.Contains(out, "Synnergy blockchain CLI") {
			t.Fatalf("unexpected help output: %s", out)
		}
	})

	t.Run("address parse", func(t *testing.T) {
		out, err := execCommand(t, "address", "parse", core.AddressZero.Hex())
		if err != nil {
			t.Fatalf("address parse failed: %v", err)
		}
		if out != core.AddressZero.Hex() {
			t.Fatalf("unexpected address output: %s", out)
		}
	})

	t.Run("gas snapshot", func(t *testing.T) {
		out, err := execCommand(t, "gas", "snapshot", "--json")
		if err != nil {
			t.Fatalf("gas snapshot failed: %v", err)
		}
		var m map[string]uint64
		if err := json.Unmarshal([]byte(out), &m); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if len(m) == 0 {
			t.Fatalf("expected non-empty gas snapshot")
		}
	})
}
