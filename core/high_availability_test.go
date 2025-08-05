package core

import (
	"testing"
	"time"
)

// TestFailoverManagerKeepsPrimaryWhenHealthy ensures that the configured
// primary node remains active as long as it has sent a recent heartbeat.
func TestFailoverManagerKeepsPrimaryWhenHealthy(t *testing.T) {
	timeout := 50 * time.Millisecond
	fm := NewFailoverManager("primary", timeout)
	fm.RegisterBackup("backup")

	// The primary heartbeat was recorded at construction time. Without
	// exceeding the timeout, Active should return the primary.
	if active := fm.Active(); active != "primary" {
		t.Fatalf("expected primary to remain active, got %s", active)
	}

	// Heartbeating a backup should not change the active node.
	fm.Heartbeat("backup")
	if active := fm.Active(); active != "primary" {
		t.Fatalf("expected primary after backup heartbeat, got %s", active)
	}
}

// TestFailoverManagerFailoverToLatestBackup verifies that when the primary
// fails to heartbeat within the timeout, the most recently heartbeating backup
// is promoted.
func TestFailoverManagerFailoverToLatestBackup(t *testing.T) {
	timeout := 5 * time.Millisecond
	fm := NewFailoverManager("p1", timeout)
	fm.RegisterBackup("b1")
	fm.RegisterBackup("b2")

	// Make the primary stale and set deterministic heartbeat times for backups.
	fm.mu.Lock()
	fm.nodes["p1"] = time.Now().Add(-2 * timeout)
	fm.nodes["b1"] = time.Now().Add(-50 * time.Millisecond)
	fm.nodes["b2"] = time.Now()
	fm.mu.Unlock()

	if active := fm.Active(); active != "b2" {
		t.Fatalf("expected b2 to become primary, got %s", active)
	}

	// Active should now consistently return the promoted node.
	if active := fm.Active(); active != "b2" {
		t.Fatalf("expected b2 to remain primary after promotion, got %s", active)
	}
}

// TestFailoverManagerHeartbeatAndRegister tests that registering a backup and
// issuing heartbeats correctly updates internal state.
func TestFailoverManagerHeartbeatAndRegister(t *testing.T) {
	timeout := 20 * time.Millisecond
	fm := NewFailoverManager("p1", timeout)

	// Record the original heartbeat timestamp for the primary.
	fm.mu.RLock()
	primaryHB := fm.nodes["p1"]
	fm.mu.RUnlock()

	// Sleep to ensure a measurable difference, then heartbeat the primary.
	time.Sleep(1 * time.Millisecond)
	fm.Heartbeat("p1")

	fm.mu.RLock()
	updatedHB := fm.nodes["p1"]
	fm.mu.RUnlock()
	if !updatedHB.After(primaryHB) {
		t.Fatalf("heartbeat did not update timestamp")
	}

	// Register and heartbeat a backup node.
	fm.RegisterBackup("b1")
	fm.Heartbeat("b1")
	fm.mu.RLock()
	_, ok := fm.nodes["b1"]
	fm.mu.RUnlock()
	if !ok {
		t.Fatalf("backup node not registered")
	}
}
