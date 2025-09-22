package auth

import (
	"errors"
	"sort"
	"sync"
)

// Permission represents an action that can be performed within the system.
type Permission string

// PermissionCondition optionally constrains a permission based on metadata.
type PermissionCondition func(metadata map[string]any) bool

// PermissionRule stores metadata about a permission assignment.
type PermissionRule struct {
	Condition PermissionCondition
}

// Role groups a set of permissions under a common name.
type Role struct {
	Name        string
	Permissions map[Permission]PermissionRule
}

// RBAC manages roles, permissions and their assignment to users.
type RBAC struct {
	mu        sync.RWMutex
	roles     map[string]*Role
	userRoles map[string]map[string]struct{}
}

// NewRBAC creates an empty RBAC instance.
func NewRBAC() *RBAC {
	return &RBAC{
		roles:     make(map[string]*Role),
		userRoles: make(map[string]map[string]struct{}),
	}
}

// AddRole registers a new role.
func (r *RBAC) AddRole(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.roles[name]; exists {
		return
	}
	r.roles[name] = &Role{Name: name, Permissions: make(map[Permission]PermissionRule)}
}

// RemoveRole removes a role and revokes it from all users.
func (r *RBAC) RemoveRole(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.roles[name]; !ok {
		return false
	}
	delete(r.roles, name)
	for user, roles := range r.userRoles {
		delete(roles, name)
		if len(roles) == 0 {
			delete(r.userRoles, user)
		}
	}
	return true
}

// AddPermissionToRole attaches a permission to an existing role.
func (r *RBAC) AddPermissionToRole(roleName string, perm Permission) error {
	return r.AddConditionalPermissionToRole(roleName, perm, nil)
}

// AddConditionalPermissionToRole attaches a permission guarded by a condition.
func (r *RBAC) AddConditionalPermissionToRole(roleName string, perm Permission, cond PermissionCondition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	role, ok := r.roles[roleName]
	if !ok {
		return errors.New("role does not exist")
	}
	role.Permissions[perm] = PermissionRule{Condition: cond}
	return nil
}

// AssignRole associates a role with a user.
func (r *RBAC) AssignRole(userID, roleName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.roles[roleName]; !ok {
		return errors.New("role does not exist")
	}
	if _, ok := r.userRoles[userID]; !ok {
		r.userRoles[userID] = make(map[string]struct{})
	}
	r.userRoles[userID][roleName] = struct{}{}
	return nil
}

// RevokeRole removes a role assignment from a user.
func (r *RBAC) RevokeRole(userID, roleName string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	roles, ok := r.userRoles[userID]
	if !ok {
		return false
	}
	if _, exists := roles[roleName]; !exists {
		return false
	}
	delete(roles, roleName)
	if len(roles) == 0 {
		delete(r.userRoles, userID)
	}
	return true
}

// ListUserRoles returns a sorted list of roles assigned to the user.
func (r *RBAC) ListUserRoles(userID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	roles := r.userRoles[userID]
	out := make([]string, 0, len(roles))
	for name := range roles {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// HasPermission checks if the user has the specified permission with metadata conditions.
func (r *RBAC) HasPermission(userID string, perm Permission, metadata map[string]any) bool {
	return r.hasPermission(userID, perm, metadata)
}

func (r *RBAC) hasPermission(userID string, perm Permission, metadata map[string]any) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	roles := r.userRoles[userID]
	for roleName := range roles {
		role := r.roles[roleName]
		if role == nil {
			continue
		}
		if rule, ok := role.Permissions[perm]; ok {
			if rule.Condition == nil || rule.Condition(metadata) {
				return true
			}
		}
	}
	return false
}

// ErrUnauthorized is returned when a user does not have the required permission.
var ErrUnauthorized = errors.New("auth: unauthorized")

// PolicyEnforcer checks user permissions and records audit events.
type PolicyEnforcer struct {
	rbac  *RBAC
	audit AuditLogger
}

// NewPolicyEnforcer creates a PolicyEnforcer with the provided RBAC store and audit logger.
func NewPolicyEnforcer(r *RBAC, l AuditLogger) *PolicyEnforcer {
	return &PolicyEnforcer{rbac: r, audit: l}
}

// Authorize verifies that the user has the specified permission and records the attempt.
func (p *PolicyEnforcer) Authorize(userID string, perm Permission, metadata map[string]any) error {
	allowed := p.rbac.hasPermission(userID, perm, metadata)
	if p.audit != nil {
		p.audit.Log(userID, perm, allowed, metadata)
	}
	if !allowed {
		return ErrUnauthorized
	}
	return nil
}
