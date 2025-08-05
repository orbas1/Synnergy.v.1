package auth

import (
	"bytes"
	"encoding/json"
	"strings"
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

	var buf bytes.Buffer
	logger := NewStdAuditLogger(&buf)
	enforcer := NewPolicyEnforcer(rbac, logger)

	// authorized request
	if err := enforcer.Authorize("u1", Permission("write"), map[string]any{"resource": "doc"}); err != nil {
		t.Fatalf("Authorize allowed returned error: %v", err)
	}

	// unauthorized request
	if err := enforcer.Authorize("u1", Permission("delete"), nil); err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}

	// verify audit logs
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 audit lines, got %d", len(lines))
	}
	// first log allowed
	var e1 map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(lines[0][strings.Index(lines[0], "{"):])), &e1); err != nil {
		t.Fatalf("unmarshal first entry: %v", err)
	}
	if allowed, ok := e1["allowed"].(bool); !ok || !allowed {
		t.Fatalf("expected first entry allowed=true, got %v", e1["allowed"])
	}
	// second log denied
	var e2 map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(lines[1][strings.Index(lines[1], "{"):])), &e2); err != nil {
		t.Fatalf("unmarshal second entry: %v", err)
	}
	if allowed, ok := e2["allowed"].(bool); !ok || allowed {
		t.Fatalf("expected second entry allowed=false, got %v", e2["allowed"])
	}
}
