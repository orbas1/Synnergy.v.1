package security

import "testing"

func TestDDoSMitigator(t *testing.T) {
	d := NewDDoSMitigator()
	d.Block("1.2.3.4")
	if !d.IsBlocked("1.2.3.4") {
		t.Fatalf("address should be blocked")
	}
}
