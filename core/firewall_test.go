package core

import "testing"

func TestFirewall(t *testing.T) {
	fw := NewFirewall()
	if !fw.IsAllowed("1.1.1.1") {
		t.Fatalf("address should be allowed by default")
	}
	fw.Block("1.1.1.1")
	if fw.IsAllowed("1.1.1.1") {
		t.Fatalf("blocked address reported as allowed")
	}
	fw.Allow("1.1.1.1")
	if !fw.IsAllowed("1.1.1.1") {
		t.Fatalf("allowed address reported as blocked")
	}
	fw.Block("2.2.2.2")
	allow, block := fw.Rules()
	if len(allow) != 1 || len(block) != 1 {
		t.Fatalf("unexpected rule counts: %v %v", allow, block)
	}
	fw.Reset()
	allow, block = fw.Rules()
	if len(allow) != 0 || len(block) != 0 {
		t.Fatalf("expected empty rules after reset")
	}
}
