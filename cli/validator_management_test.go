package cli

import (
	"strings"
	"testing"
)

// TestValidatorAddStake verifies adding a validator and querying its stake.
func TestValidatorAddStake(t *testing.T) {
	if _, err := execCommand("validator", "add", "val1", "100", "--json"); err != nil {
		t.Fatalf("add: %v", err)
	}
	out, err := execCommand("validator", "stake", "val1", "--json")
	if err != nil {
		t.Fatalf("stake: %v", err)
	}
	if err := rootCmd.PersistentFlags().Set("json", "false"); err != nil {
		t.Fatalf("reset json: %v", err)
	}
	if !strings.Contains(out, "gas cost") || !strings.Contains(out, "100") {
		t.Fatalf("unexpected output: %s", out)
	}
}
