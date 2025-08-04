package nodes

import (
	"testing"
	"time"
)

func TestElectedAuthorityNode(t *testing.T) {
	n := NewElectedAuthorityNode("addr", "role", time.Hour)
	if !n.IsActive(time.Now()) {
		t.Fatalf("node should be active")
	}
	if !n.IsActive(time.Now().Add(time.Minute)) {
		t.Fatalf("node should still be active")
	}
	if n.IsActive(time.Now().Add(2 * time.Hour)) {
		t.Fatalf("node should be inactive after term")
	}
}
