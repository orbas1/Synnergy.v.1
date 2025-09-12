package core

import (
	"sync"
	"testing"
)

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

func TestValidatorNodeConcurrentAdds(t *testing.T) {
	ledger := NewLedger()
	vn := NewValidatorNode("v1", "addr", ledger, 1, 3)
	ids := []string{"a", "b", "c"}
	var wg sync.WaitGroup
	for _, id := range ids {
		id := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = vn.AddValidator(id, 5)
		}()
	}
	wg.Wait()
	if !vn.HasQuorum() {
		t.Fatalf("quorum not reached with concurrent adds")
	}
}
