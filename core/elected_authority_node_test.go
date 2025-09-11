package core

import (
	"testing"
	"time"
)

func TestElectedAuthorityNode(t *testing.T) {
	dao := &DAO{Members: map[string]string{"admin": RoleAdmin, "member": RoleMember}}
	en := NewElectedAuthorityNode("addr", "validator", time.Minute)
	if !en.IsActive(time.Now()) {
		t.Fatalf("expected active")
	}
	if en.IsActive(time.Now().Add(time.Hour)) {
		t.Fatalf("expected inactive after term")
	}
	// non-admin cannot renew
	if err := en.RenewTerm("member", dao, time.Hour); err == nil {
		t.Fatalf("expected unauthorized renewal")
	}
	// admin renews term
	if err := en.RenewTerm("admin", dao, time.Hour); err != nil {
		t.Fatalf("renew term: %v", err)
	}
	if !en.IsActive(time.Now().Add(59 * time.Minute)) {
		t.Fatalf("expected renewed term")
	}
}
