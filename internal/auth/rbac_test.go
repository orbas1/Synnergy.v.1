package auth

import (
	"testing"
)

func TestPolicyEnforcerAuthorize(t *testing.T) {
	rbac := NewRBAC()
	rbac.AddRole("admin")
	if err := rbac.AddPermissionToRole("admin", Permission("write")); err != nil {
		t.Fatalf("AddPermissionToRole: %v", err)
	}
	if err := rbac.AssignRole("u1", "admin"); err != nil {
		t.Fatalf("AssignRole: %v", err)
	}

	logger := NewMemoryAuditLogger()
	enforcer := NewPolicyEnforcer(rbac, logger)

	if err := enforcer.Authorize("u1", Permission("write"), map[string]any{"resource": "doc"}); err != nil {
		t.Fatalf("Authorize allowed returned error: %v", err)
	}
	if err := enforcer.Authorize("u1", Permission("delete"), nil); err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	events := logger.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 audit events, got %d", len(events))
	}
	if !events[0].Allowed || events[1].Allowed {
		t.Fatalf("unexpected audit allowed flags: %+v", events)
	}
}

func TestConditionalPermissions(t *testing.T) {
	rbac := NewRBAC()
	rbac.AddRole("approver")
	err := rbac.AddConditionalPermissionToRole("approver", Permission("approve"), func(metadata map[string]any) bool {
		return metadata["region"] == "eu"
	})
	if err != nil {
		t.Fatalf("AddConditionalPermissionToRole: %v", err)
	}
	if err := rbac.AssignRole("user", "approver"); err != nil {
		t.Fatalf("AssignRole: %v", err)
	}
	if !rbac.HasPermission("user", Permission("approve"), map[string]any{"region": "eu"}) {
		t.Fatalf("expected approval for eu region")
	}
	if rbac.HasPermission("user", Permission("approve"), map[string]any{"region": "us"}) {
		t.Fatalf("expected us region to be rejected")
	}
}

func TestRoleManagement(t *testing.T) {
	rbac := NewRBAC()
	rbac.AddRole("viewer")
	if err := rbac.AssignRole("user", "viewer"); err != nil {
		t.Fatalf("assign: %v", err)
	}
	roles := rbac.ListUserRoles("user")
	if len(roles) != 1 || roles[0] != "viewer" {
		t.Fatalf("unexpected roles: %v", roles)
	}
	if !rbac.RevokeRole("user", "viewer") {
		t.Fatalf("expected revoke to succeed")
	}
	if len(rbac.ListUserRoles("user")) != 0 {
		t.Fatalf("expected no roles after revoke")
	}
	if !rbac.RemoveRole("viewer") {
		t.Fatalf("expected remove role to succeed")
	}
	if rbac.RemoveRole("viewer") {
		t.Fatalf("expected removing non-existent role to return false")
	}
}
