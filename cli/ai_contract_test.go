package cli

import "testing"

// TestAIContractModelHashNotFound ensures querying an unknown contract returns not found.
func TestAIContractModelHashNotFound(t *testing.T) {
	if _, ok := aiRegistry.ModelHash("0xdeadbeef"); ok {
		t.Fatalf("unexpected contract found")
	}
}
