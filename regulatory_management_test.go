package synnergy

import "testing"

func TestRegulatoryManagerEvaluation(t *testing.T) {
	m := NewRegulatoryManager()
	reg := Regulation{
		ID:               "r1",
		MaxAmount:        100,
		MinAmount:        10,
		AllowedEntities:  []string{"alice"},
		RequireWhitelist: true,
	}
	if err := m.AddRegulation(reg); err != nil {
		t.Fatalf("add regulation: %v", err)
	}

	tx := Transaction{From: "alice", To: "bob", Amount: 50}
	if err := m.ValidateTransaction(tx); err != nil {
		t.Fatalf("unexpected violation: %v", err)
	}

	tx.Amount = 150
	res := m.EvaluateTransactionDetailed(tx)
	if len(res.Violations) != 1 || res.Violations[0].RegulationID != "r1" {
		t.Fatalf("expected amount violation")
	}

	tx.Amount = 5
	if err := m.ValidateTransaction(tx); err == nil {
		t.Fatalf("expected minimum amount violation")
	}

	tx.Amount = 50
	tx.From = "mallory"
	if err := m.ValidateTransaction(tx); err == nil {
		t.Fatalf("expected whitelist violation")
	}

	if err := m.UpdateRegulation(Regulation{ID: "r1", MaxAmount: 200}); err != nil {
		t.Fatalf("update regulation: %v", err)
	}
	if got := m.ListRegulations(); len(got) != 1 || got[0].MaxAmount != 200 {
		t.Fatalf("update not applied: %+v", got)
	}
}

func TestRegulatoryManagerDuplicate(t *testing.T) {
	m := NewRegulatoryManager()
	reg := Regulation{ID: "dup"}
	if err := m.AddRegulation(reg); err != nil {
		t.Fatalf("add regulation: %v", err)
	}
	if err := m.AddRegulation(reg); err == nil {
		t.Fatalf("expected duplicate error")
	}
	m.RemoveRegulation("dup")
	if _, ok := m.GetRegulation("dup"); ok {
		t.Fatalf("regulation not removed")
	}
}
