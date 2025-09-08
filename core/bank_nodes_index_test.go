package core

import (
	"encoding/json"
	"testing"
)

func TestBankNodeIndex(t *testing.T) {
	idx := NewBankNodeIndex()
	rec := &BankNodeRecord{ID: "id1", Type: BankInstitutionalNodeType}
	idx.Add(rec)
	if got, ok := idx.Get("id1"); !ok || got.Type != BankInstitutionalNodeType {
		t.Fatalf("expected record, got %v", got)
	}
	idx.Remove("id1")
	if _, ok := idx.Get("id1"); ok {
		t.Fatalf("record not removed")
	}
}

func TestBankNodeIndexJSON(t *testing.T) {
	idx := NewBankNodeIndex()
	idx.Add(&BankNodeRecord{ID: "id1", Type: BankInstitutionalNodeType})
	if _, err := json.Marshal(idx); err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
}
