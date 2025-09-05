package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"synnergy/core"
)

// TestBankNodeTypes lists available bank node types via the CLI and verifies output.
func TestBankNodeTypes(t *testing.T) {
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"banknodes", "types"})
		_ = rootCmd.Execute()
	})

	if !strings.Contains(output, core.BankInstitutionalNodeType) {
		t.Fatalf("expected bank node types in output, got %s", output)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
