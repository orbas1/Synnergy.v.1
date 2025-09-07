package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func resetSyn130() {
	syn130Registry = core.NewTangibleAssetRegistry()
	syn130IDs = nil
}

func TestTokenSyn130RegisterAndList(t *testing.T) {
	resetSyn130()
	if _, err := execCommand("token_syn130", "register", "asset1", "alice", "meta", "100"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	out, err := execCommand("token_syn130", "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "asset1") {
		t.Fatalf("expected asset in list, got %s", out)
	}
}
