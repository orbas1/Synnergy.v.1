package core

import "testing"

func TestValidatorNode(t *testing.T) {
	ledger := NewLedger()
	vn := NewValidatorNode("v1", "addr", ledger, 1, 1)
	if err := vn.AddValidator("v1", 10); err != nil {
		t.Fatalf("add validator failed: %v", err)
	}
	if !vn.HasQuorum() {
		t.Fatalf("quorum should be satisfied after adding validator")
	}
	vn.SlashValidator("v1")
	if vn.HasQuorum() {
		t.Fatalf("quorum should not hold after slashing")
	}
}
