package core

import "testing"

func TestSYN131RegistryCreateGet(t *testing.T) {
	reg := NewSYN131Registry()
	tok, err := reg.Create("t1", "Name", "SYM", "alice", 100)
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if tok.ID != "t1" || tok.Name != "Name" || tok.Symbol != "SYM" || tok.Owner != "alice" || tok.Valuation != 100 {
		t.Fatalf("unexpected token data: %#v", tok)
	}
	got, ok := reg.Get("t1")
	if !ok || got.ID != "t1" {
		t.Fatalf("get failed")
	}
}

func TestSYN131RegistryDuplicate(t *testing.T) {
	reg := NewSYN131Registry()
	if _, err := reg.Create("t1", "n", "s", "o", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := reg.Create("t1", "n", "s", "o", 1); err == nil {
		t.Fatalf("expected duplicate error")
	}
}

func TestSYN131RegistryUpdateValuation(t *testing.T) {
	reg := NewSYN131Registry()
	if _, err := reg.Create("t1", "n", "s", "o", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := reg.UpdateValuation("t1", 500); err != nil {
		t.Fatalf("update valuation error: %v", err)
	}
	tok, _ := reg.Get("t1")
	if tok.Valuation != 500 {
		t.Fatalf("valuation not updated: %d", tok.Valuation)
	}
}

func TestSYN131RegistryUpdateNonexistent(t *testing.T) {
	reg := NewSYN131Registry()
	if err := reg.UpdateValuation("missing", 10); err == nil {
		t.Fatalf("expected error for missing token")
	}
}
