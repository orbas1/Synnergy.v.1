package cli

import "testing"

// TestAIContractListEmpty verifies that the list subcommand returns an empty JSON array
// when no contracts have been deployed.
func TestAIContractListEmpty(t *testing.T) {
	out, err := execCommand("ai_contract", "list", "--json")
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if out != "[]" && out != "null" {
		t.Fatalf("expected empty list, got %q", out)
	}
}
