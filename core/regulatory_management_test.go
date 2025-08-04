package core

import "testing"

func TestRegulatoryManager(t *testing.T) {
	m := NewRegulatoryManager()
	reg := Regulation{ID: "r1", MaxAmount: 100}
	if err := m.AddRegulation(reg); err != nil {
		t.Fatalf("add regulation: %v", err)
	}
	tx := Transaction{Amount: 150}
	v := m.EvaluateTransaction(tx)
	if len(v) != 1 || v[0] != "r1" {
		t.Fatalf("expected violation r1")
	}
}
