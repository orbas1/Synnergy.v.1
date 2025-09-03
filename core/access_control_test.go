package core

import "testing"

func TestAccessController(t *testing.T) {
	GlobalAccess = NewAccessController()
	addr := "0x1111111111111111111111111111111111111111"
	if err := GrantRole("admin", addr); err != nil {
		t.Fatalf("grant: %v", err)
	}
	ok, err := HasRole("admin", addr)
	if err != nil || !ok {
		t.Fatalf("expected role granted")
	}
	roles, err := ListRoles(addr)
	if err != nil || len(roles) != 1 || roles[0] != "admin" {
		t.Fatalf("unexpected roles: %v", roles)
	}
	if err := RevokeRole("admin", addr); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	ok, err = HasRole("admin", addr)
	if err != nil {
		t.Fatalf("has role: %v", err)
	}
	if ok {
		t.Fatalf("role should be revoked")
	}
	if _, err := HasRole("", addr); err == nil {
		t.Fatalf("expected error for empty role")
	}
	if err := GrantRole("admin", "invalid"); err == nil {
		t.Fatalf("expected invalid address error")
	}
}
