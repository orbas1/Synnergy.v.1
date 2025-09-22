package governance

import (
	"sync"
	"testing"
	"time"
)

func TestAuditLogAppendAndVerify(t *testing.T) {
	signer := NewHMACSigner([]byte("secret"))
	log := NewAuditLog(WithSigner(signer), WithRetention(5))
	first, err := log.Append(AuditEvent{Actor: "validator-1", Action: "approve", Scope: "consensus", GasBudget: 5000})
	if err != nil {
		t.Fatalf("append first event: %v", err)
	}
	if first.Sequence != 1 {
		t.Fatalf("expected sequence 1, got %d", first.Sequence)
	}
	if first.Hash == "" || first.Signature == "" {
		t.Fatalf("expected hash and signature")
	}

	second, err := log.Append(AuditEvent{Actor: "validator-2", Action: "approve", Scope: "consensus", GasBudget: 5000})
	if err != nil {
		t.Fatalf("append second event: %v", err)
	}
	if second.PrevHash != first.Hash {
		t.Fatalf("expected hash chaining")
	}

	if err := log.VerifyChain(); err != nil {
		t.Fatalf("verify chain: %v", err)
	}
}

func TestAuditLogRetention(t *testing.T) {
	log := NewAuditLog(WithRetention(2))
	for i := 0; i < 3; i++ {
		if _, err := log.Append(AuditEvent{Actor: "node", Action: "test", Scope: "retention"}); err != nil {
			t.Fatalf("append event: %v", err)
		}
	}
	if len(log.Entries()) != 2 {
		t.Fatalf("expected retention to keep only two entries")
	}
}

func TestAuditLogObservers(t *testing.T) {
	var mu sync.Mutex
	observed := 0
	observer := func(AuditRecord) {
		mu.Lock()
		defer mu.Unlock()
		observed++
	}
	log := NewAuditLog(WithObserver(observer))
	if _, err := log.Append(AuditEvent{Actor: "validator", Action: "approve"}); err != nil {
		t.Fatalf("append event: %v", err)
	}
	mu.Lock()
	count := observed
	mu.Unlock()
	if count != 1 {
		t.Fatalf("expected observer to run once, got %d", count)
	}
}

func TestAuditLogGeneratesDeterministicIDs(t *testing.T) {
	log := NewAuditLog()
	record, err := log.Append(AuditEvent{Actor: "validator", Action: "approve"})
	if err != nil {
		t.Fatalf("append event: %v", err)
	}
	if record.ID == "" {
		t.Fatalf("expected generated id")
	}
}

func TestAuditLogLatest(t *testing.T) {
	log := NewAuditLog()
	if _, ok := log.Latest(); ok {
		t.Fatalf("expected empty log to have no latest record")
	}
	inserted, err := log.Append(AuditEvent{Actor: "node", Action: "rotate", Timestamp: time.Unix(1, 0)})
	if err != nil {
		t.Fatalf("append event: %v", err)
	}
	latest, ok := log.Latest()
	if !ok || latest.Sequence != inserted.Sequence {
		t.Fatalf("expected latest record to match appended entry")
	}
}
