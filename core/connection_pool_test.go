package core

import "testing"

func TestConnectionPool(t *testing.T) {
	p := NewConnectionPool(1)
	if _, err := p.Acquire("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := p.Acquire("b"); err == nil {
		t.Fatalf("expected pool exhaustion error")
	}
	if p.Size() != 1 {
		t.Fatalf("expected size 1, got %d", p.Size())
	}
	p.Release("a")
	if p.Size() != 0 {
		t.Fatalf("release failed")
	}
	if stats := p.Stats(); stats.Capacity != 1 || stats.Active != 0 {
		t.Fatalf("unexpected stats %+v", stats)
	}
}
