package core

import (
	"testing"
	"time"
)

func TestElectedAuthorityNode(t *testing.T) {
	en := NewElectedAuthorityNode("addr", "validator", time.Minute)
	if !en.IsActive(time.Now()) {
		t.Fatalf("expected active")
	}
	if en.IsActive(time.Now().Add(time.Hour)) {
		t.Fatalf("expected inactive after term")
	}
}
