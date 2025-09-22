package security

import (
	"testing"
	"time"
)

func TestDDoSMitigatorAllowAndBlock(t *testing.T) {
	cfg := MitigationConfig{Window: time.Second, MaxRequests: 3, BurstAllowance: 1, BlockDuration: 10 * time.Millisecond}
	d := NewDDoSMitigator(cfg)
	now := time.Now()
	for i := 0; i < 3; i++ {
		if !d.Allow("1.2.3.4", now.Add(time.Millisecond*time.Duration(i))) {
			t.Fatalf("request %d should be allowed", i)
		}
	}
	d.Allow("1.2.3.4", now.Add(10*time.Millisecond))
	if d.Allow("1.2.3.4", now.Add(11*time.Millisecond)) {
		t.Fatalf("expected block after exceeding limit")
	}
	if !d.IsBlocked("1.2.3.4", now.Add(20*time.Millisecond)) {
		t.Fatalf("expected blocked state")
	}
	time.Sleep(cfg.BlockDuration * 2)
	if d.IsBlocked("1.2.3.4", now.Add(cfg.BlockDuration*3)) {
		t.Fatalf("expected unblock after duration")
	}
}

func TestDDoSMitigatorSnapshotOrdering(t *testing.T) {
	d := NewDDoSMitigator(MitigationConfig{Window: time.Second, MaxRequests: 1, BurstAllowance: 1})
	now := time.Now()
	d.Allow("10.0.0.1", now)
	for i := 0; i < 3; i++ {
		d.Allow("10.0.0.2", now.Add(time.Millisecond*time.Duration(i)))
	}
	snap := d.Snapshot(time.Now())
	if len(snap) != 2 {
		t.Fatalf("unexpected snapshot length: %d", len(snap))
	}
	if snap[0].IP != "10.0.0.2" {
		t.Fatalf("expected highest score first, got %s", snap[0].IP)
	}
}
