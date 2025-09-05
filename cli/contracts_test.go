package cli

import "testing"

func TestContractsList(t *testing.T) {
	out, err := execCommand("contracts", "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty list, got %s", out)
	}
}
