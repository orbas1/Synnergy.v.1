package core

import (
	"testing"
	"time"
)

func TestSYN500GrantLifecycle(t *testing.T) {
	token := NewSYN500Token("Loyalty", "LOY", "owner", 2, 10)
	if err := token.Grant("alice", 1, 2, time.Minute); err != nil {
		t.Fatalf("grant: %v", err)
	}
	if err := token.Use("alice", time.Now()); err != nil {
		t.Fatalf("use1: %v", err)
	}
	if err := token.Use("alice", time.Now()); err != nil {
		t.Fatalf("use2: %v", err)
	}
	if err := token.Use("alice", time.Now()); err == nil {
		t.Fatal("expected limit reached error")
	}
	status, ok := token.Status("alice")
	if !ok || status.Used != 2 {
		t.Fatalf("unexpected status: %#v", status)
	}
	tele := token.Telemetry()
	if tele.Grants != 1 {
		t.Fatalf("unexpected telemetry: %#v", tele)
	}
}
