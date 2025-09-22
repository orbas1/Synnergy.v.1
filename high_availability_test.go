package synnergy

import (
	"testing"
	"time"
)

func TestFailoverManager(t *testing.T) {
	timeout := 10 * time.Millisecond
	m := NewFailoverManager("primary", timeout)
	m.RegisterBackup("backup")

	// Simulate primary missing heartbeats.
	m.mu.Lock()
	m.nodes["primary"] = time.Now().Add(-2 * timeout)
	m.mu.Unlock()

	active := m.Active()
	if active != "backup" {
		t.Fatalf("expected failover to backup, got %s", active)
	}
}

func TestFailoverSnapshot(t *testing.T) {
	timeout := 25 * time.Millisecond
	m := NewFailoverManager("primary", timeout)
	m.RegisterBackup("backup")
	m.Heartbeat("backup")

	m.mu.Lock()
	m.nodes["primary"] = time.Now().Add(-2 * timeout)
	m.mu.Unlock()

	if active := m.Active(); active != "backup" {
		t.Fatalf("expected failover to promote backup, got %s", active)
	}

	// Refresh the backup heartbeat to ensure the snapshot reports a healthy
	// primary after the promotion.
	m.Heartbeat("backup")

	snap := m.Snapshot()
	if snap.Active != "backup" {
		t.Fatalf("expected snapshot to report backup as active, got %s", snap.Active)
	}
	if !snap.Healthy {
		t.Fatalf("expected snapshot to report active node healthy")
	}
	if snap.Failovers != 1 {
		t.Fatalf("expected exactly one failover, got %d", snap.Failovers)
	}
	if snap.LastSwitch.IsZero() {
		t.Fatalf("expected last switch timestamp to be recorded")
	}
	if snap.Timeout != timeout {
		t.Fatalf("expected snapshot timeout %s, got %s", timeout, snap.Timeout)
	}
}
