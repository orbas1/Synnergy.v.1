package tokens

import (
	"sync"
	"testing"
	"time"
)

func TestSYN2700Lifecycle(t *testing.T) {
	tok := NewSYN2700Token()
	if err := tok.AddHolder("alice", 0); err != ErrInvalidAmount {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}

	if err := tok.AddHolder("alice", 40); err != nil {
		t.Fatalf("add holder: %v", err)
	}
	if err := tok.AddHolder("bob", 60); err != nil {
		t.Fatalf("add holder: %v", err)
	}
	if tok.TotalSupply() != 100 {
		t.Fatalf("unexpected supply: %d", tok.TotalSupply())
	}

	if err := tok.UpdateHolder("alice", 25); err != nil {
		t.Fatalf("update holder: %v", err)
	}
	if bal, ok := tok.HolderBalance("alice"); !ok || bal != 25 {
		t.Fatalf("unexpected balance: %d %v", bal, ok)
	}

	if err := tok.RemoveHolder("carol"); err != ErrHolderNotFound {
		t.Fatalf("expected ErrHolderNotFound, got %v", err)
	}
	if err := tok.RemoveHolder("alice"); err != nil {
		t.Fatalf("remove holder: %v", err)
	}
	if _, ok := tok.HolderBalance("alice"); ok {
		t.Fatal("alice should be removed")
	}
}

func TestSYN2700DistributeDetailed(t *testing.T) {
	tok := NewSYN2700Token()
	must := func(err error) {
		if err != nil {
			t.Helper()
			t.Fatalf("unexpected error: %v", err)
		}
	}
	must(tok.AddHolder("alice", 10))
	must(tok.AddHolder("bob", 10))
	must(tok.AddHolder("carol", 10))

	res, err := tok.DistributeDetailed(100)
	if err != nil {
		t.Fatalf("distribute detailed: %v", err)
	}
	if res.Unallocated != 0 || res.Distributed != 100 {
		t.Fatalf("unexpected totals: %+v", res)
	}
	for _, addr := range []string{"alice", "bob", "carol"} {
		if res.Amounts[addr] != 33 && res.Amounts[addr] != 34 {
			t.Fatalf("unexpected allocation for %s: %d", addr, res.Amounts[addr])
		}
	}

	history := tok.History(1)
	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}
	if history[0].Distributed != 100 {
		t.Fatalf("unexpected history distributed: %d", history[0].Distributed)
	}
}

func TestSYN2700ConcurrentSnapshots(t *testing.T) {
	tok := NewSYN2700Token()
	must := func(err error) {
		if err != nil {
			t.Helper()
			t.Fatalf("unexpected error: %v", err)
		}
	}
	for i := 0; i < 10; i++ {
		must(tok.AddHolder(time.Now().Add(time.Duration(i)*time.Millisecond).String(), 1))
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = tok.Snapshot()
			_ = tok.TotalSupply()
		}()
	}
	wg.Wait()
}
