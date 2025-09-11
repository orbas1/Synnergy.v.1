package core

import "testing"

func TestRegulatoryManager(t *testing.T) {
	m := NewRegulatoryManager()
	reg := Regulation{ID: "r1", MaxAmount: 100}
	if err := m.AddRegulation(reg); err != nil {
		t.Fatalf("add regulation: %v", err)
	}

	tx := Transaction{Amount: 150}
	if err := m.ValidateTransaction(tx); err == nil {
		t.Fatalf("expected validation error")
	}

	if err := m.UpdateRegulation(Regulation{ID: "r1", MaxAmount: 200}); err != nil {
		t.Fatalf("update: %v", err)
	}

	tx.Amount = 150
	if err := m.ValidateTransaction(tx); err != nil {
		t.Fatalf("unexpected violation after update: %v", err)
	}
}
