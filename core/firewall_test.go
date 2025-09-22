package core

import (
	"testing"
	"time"
)

func TestFirewall(t *testing.T) {
	fw := NewFirewall()
	if !fw.IsAllowed("1.1.1.1") {
		t.Fatalf("address should be allowed by default")
	}
	fw.BlockWithOptions("1.1.1.1", RuleOptions{Reason: "blacklist"})
	if fw.IsAllowed("1.1.1.1") {
		t.Fatalf("blocked address reported as allowed")
	}
	fw.AllowWithOptions("1.1.1.1", RuleOptions{Reason: "allow", Persistent: true})
	if !fw.IsAllowed("1.1.1.1") {
		t.Fatalf("allowed address reported as blocked")
	}
	fw.BlockWithOptions("2.2.2.2", RuleOptions{TTL: 10 * time.Millisecond, Reason: "temp"})
	allow, block := fw.Rules()
	if len(allow) != 1 || len(block) != 1 {
		t.Fatalf("unexpected rule counts: %v %v", allow, block)
	}
	details := fw.RuleDetails()
	if len(details) != 2 {
		t.Fatalf("expected rule details for both allow and block")
	}
	fw.SetDefaultAllow(false)
	if fw.IsAllowed("3.3.3.3") {
		t.Fatalf("unexpected allow when default disabled")
	}
	time.Sleep(15 * time.Millisecond)
	if removed := fw.PruneExpired(); removed == 0 {
		t.Fatalf("expected expired rule to be pruned")
	}
	fw.SetDefaultAllow(true)
	if !fw.IsAllowed("2.2.2.2") {
		t.Fatalf("expired block rule should allow address")
	}
	fw.Reset()
	allow, block = fw.Rules()
	if len(allow) != 0 || len(block) != 0 {
		t.Fatalf("expected empty rules after reset")
	}
}
