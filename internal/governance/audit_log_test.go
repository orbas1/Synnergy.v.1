package governance

import "testing"

func TestAuditLog(t *testing.T) {
	a := NewAuditLog()
	a.Append("entry")
	if len(a.Entries()) != 1 {
		t.Fatalf("expected one entry")
	}
}
