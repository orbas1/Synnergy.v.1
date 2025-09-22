package synnergy

import (
	"errors"
	"testing"
	"time"
)

func TestBridgeTransferManagerLifecycle(t *testing.T) {
	mgr := NewBridgeTransferManager()
	events := make(chan BridgeTransferEvent, 32)
	mgr.RegisterListener(func(ev BridgeTransferEvent) { events <- ev })

	id1 := mgr.Deposit("bridge-1", "alice", "bob", 100, "syn", WithTransferMetadata("note", "audit"))
	if id1 == "" {
		t.Fatalf("expected transfer id")
	}
	if err := mgr.Fail(id1, errors.New("fraudulent")); err != nil {
		t.Fatalf("fail: %v", err)
	}

	id2 := mgr.Deposit("bridge-1", "carol", "dave", 50, "syn")
	if err := mgr.Claim(id2, []byte("proof")); err != nil {
		t.Fatalf("claim: %v", err)
	}

	expiredAt := time.Now().Add(-time.Minute)
	id3 := mgr.Deposit("bridge-1", "erin", "frank", 10, "syn", WithTransferExpiry(expiredAt))
	expired := mgr.SweepExpired(time.Now())
	if len(expired) != 1 || expired[0].ID != id3 {
		t.Fatalf("expected transfer %s to expire", id3)
	}

	metrics := mgr.Metrics()
	if metrics.Total != 3 || metrics.Claimed != 1 || metrics.Failed != 1 || metrics.Expired != 1 {
		t.Fatalf("unexpected metrics %+v", metrics)
	}

	close(events)
	seen := map[BridgeTransferEventType]int{}
	for ev := range events {
		seen[ev.Type]++
	}
	required := []BridgeTransferEventType{TransferEventDeposited, TransferEventFailed, TransferEventClaimed, TransferEventExpired}
	for _, typ := range required {
		if seen[typ] == 0 {
			t.Fatalf("expected event %s", typ)
		}
	}
}
