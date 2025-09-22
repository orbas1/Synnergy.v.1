package governance

import (
	"testing"
	"time"
)

func TestReplayProtectorSeen(t *testing.T) {
	r := NewReplayProtector()
	if r.Seen("abc") {
		t.Fatal("ID should not be seen first time")
	}
	if !r.Seen("abc") {
		t.Fatal("ID should be seen second time")
	}
}

func TestReplayProtectorWindowExpiry(t *testing.T) {
	r := NewReplayProtector(WithWindow(1 * time.Second))
	ts := time.Unix(0, 0)
	if err := r.Check("abc", ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Check("abc", ts.Add(500*time.Millisecond)); !IsReplayError(err) {
		t.Fatalf("expected replay within window")
	}
	if err := r.Check("abc", ts.Add(2*time.Second)); err != nil {
		t.Fatalf("expected identifier to expire, got %v", err)
	}
}

func TestReplayProtectorEvictsOldest(t *testing.T) {
	evicted := make(chan string, 1)
	r := NewReplayProtector(WithMaxEntries(1), WithEvictionCallback(func(id string) { evicted <- id }))
	if err := r.Check("first", time.Now()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Check("second", time.Now()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	select {
	case id := <-evicted:
		if id != "first" {
			t.Fatalf("expected first to be evicted, got %s", id)
		}
	default:
		t.Fatalf("expected eviction callback to fire")
	}
}

func TestReplayProtectorDuplicateCallback(t *testing.T) {
	duplicates := make(chan string, 1)
	r := NewReplayProtector(WithDuplicateCallback(func(id string) { duplicates <- id }))
	_ = r.Check("dup", time.Now())
	if err := r.Check("dup", time.Now()); err == nil {
		t.Fatalf("expected replay error")
	}
	select {
	case id := <-duplicates:
		if id != "dup" {
			t.Fatalf("expected duplicate callback for dup, got %s", id)
		}
	default:
		t.Fatalf("expected duplicate callback to trigger")
	}
}

func IsReplayError(err error) bool {
	return err == ErrReplayDetected
}
