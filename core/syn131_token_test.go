package core

import (
	"fmt"
	"sync"
	"testing"
)

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
	if _, err := reg.Create("t1", "n", "s", "o", 1); err != ErrTokenExists {
		t.Fatalf("expected ErrTokenExists, got %v", err)
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
	if err := reg.UpdateValuation("missing", 10); err != ErrTokenNotFound {
		t.Fatalf("expected ErrTokenNotFound, got %v", err)
	}
}

func TestSYN131RegistryConcurrentCreate(t *testing.T) {
	reg := NewSYN131Registry()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("t-%d", i)
			if _, err := reg.Create(id, "n", "s", "o", 1); err != nil {
				t.Errorf("create %s: %v", id, err)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < 50; i++ {
		if _, ok := reg.Get(fmt.Sprintf("t-%d", i)); !ok {
			t.Fatalf("missing token %d", i)
		}
	}
}
