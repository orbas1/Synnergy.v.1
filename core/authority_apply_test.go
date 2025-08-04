package core

import (
	"testing"
	"time"
)

func TestAuthorityApplication(t *testing.T) {
	reg := NewAuthorityNodeRegistry()
	mgr := NewAuthorityApplicationManager(reg, time.Hour)
	id := mgr.Submit("cand1", "validator", "test")
	if err := mgr.Vote("voter1", id, true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := mgr.Finalize(id); err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if !reg.IsAuthorityNode("cand1") {
		t.Fatalf("candidate not registered after finalise")
	}
}
