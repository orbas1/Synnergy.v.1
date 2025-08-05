package synnergy

import "testing"

func TestAccessController(t *testing.T) {
	ac := NewAccessController()
	ac.Grant("admin", "addr1")
	if !ac.HasRole("admin", "addr1") {
		t.Fatalf("expected role granted")
	}
	roles := ac.List("addr1")
	if len(roles) != 1 || roles[0] != "admin" {
		t.Fatalf("unexpected roles: %v", roles)
	}
	ac.Revoke("admin", "addr1")
	if ac.HasRole("admin", "addr1") {
		t.Fatalf("role should be revoked")
	}
}
