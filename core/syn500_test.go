package core

import (
	"sync"
	"testing"
)

func TestSYN500GrantAndUse(t *testing.T) {
	token := NewSYN500Token("Utility", "UTL", "owner", 2, 1_000)
	if err := token.Grant("alice", 1, 10); err != nil {
		t.Fatalf("grant failed: %v", err)
	}

	usage, ok := token.Usage("alice")
	if !ok {
		t.Fatal("expected usage for alice")
	}
	if usage.Max != 10 || usage.Tier != 1 {
		t.Fatalf("unexpected tier %+v", usage)
	}

	record, err := token.Use("alice", 4, "api call")
	if err != nil {
		t.Fatalf("use failed: %v", err)
	}
	if record.Remaining != 6 {
		t.Fatalf("expected remaining 6, got %d", record.Remaining)
	}
	if record.Digest == "" {
		t.Fatal("expected audit digest")
	}

	if _, err = token.Use("alice", 7, "overflow"); err == nil {
		t.Fatal("expected usage overflow error")
	}
}

func TestSYN500Concurrency(t *testing.T) {
	token := NewSYN500Token("Utility", "UTL", "owner", 2, 1_000)
	if err := token.Grant("bob", 1, 100); err != nil {
		t.Fatalf("grant failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := token.Use("bob", 5, "batch"); err != nil {
				t.Errorf("use error: %v", err)
			}
		}()
	}
	wg.Wait()

	usage, ok := token.Usage("bob")
	if !ok {
		t.Fatal("expected usage for bob")
	}
	if usage.Used != 50 {
		t.Fatalf("expected used 50, got %d", usage.Used)
	}

	snapshot := token.Snapshot()
	if snapshot["audits"].(int) != 10 {
		t.Fatalf("expected 10 audit entries, got %v", snapshot["audits"])
	}

	audits := token.AuditTrail(5)
	if len(audits) != 5 {
		t.Fatalf("expected 5 audit entries, got %d", len(audits))
	}
}

func TestSYN500GrantValidation(t *testing.T) {
	token := NewSYN500Token("Utility", "UTL", "owner", 2, 1_000)
	if err := token.Grant("", 1, 1); err == nil {
		t.Fatal("expected invalid grant error")
	}
	if err := token.Grant("alice", 0, 1); err == nil {
		t.Fatal("expected invalid tier error")
	}
	if err := token.Grant("alice", 1, 0); err == nil {
		t.Fatal("expected invalid max error")
	}
	if err := token.Grant("alice", 1, 2); err != nil {
		t.Fatalf("grant failed: %v", err)
	}
	if err := token.Grant("alice", 2, 1); err != nil {
		t.Fatalf("update grant failed: %v", err)
	}
	usage, _ := token.Usage("alice")
	if usage.Tier != 2 || usage.Max != 1 {
		t.Fatalf("unexpected updated grant: %+v", usage)
	}
	token.Revoke("alice")
	if _, ok := token.Usage("alice"); ok {
		t.Fatal("expected grant revoked")
	}
}
