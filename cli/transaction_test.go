package cli

import (
	"strings"
	"testing"
)

// TestTransactionVariableFee ensures the variable fee command emits gas and result output.
func TestTransactionVariableFee(t *testing.T) {
	out, err := execCommand("tx", "variablefee", "10", "2", "--json")
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "20") {
		t.Fatalf("unexpected output: %s", out)
	}
}
