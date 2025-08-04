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
