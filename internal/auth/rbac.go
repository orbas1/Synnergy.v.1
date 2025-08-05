package auth

import (
	"errors"
	"sync"
)

// Permission represents an action that can be performed within the system.
type Permission string

// Role groups a set of permissions under a common name.
type Role struct {
	Name        string
	Permissions map[Permission]struct{}
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
	r.roles[name] = &Role{Name: name, Permissions: make(map[Permission]struct{})}
}

// AddPermissionToRole attaches a permission to an existing role.
func (r *RBAC) AddPermissionToRole(roleName string, perm Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	role, ok := r.roles[roleName]
	if !ok {
		return errors.New("role does not exist")
	}
	role.Permissions[perm] = struct{}{}
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

// hasPermission reports whether the given user has the provided permission.
func (r *RBAC) hasPermission(userID string, perm Permission) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	roles := r.userRoles[userID]
	for roleName := range roles {
		role := r.roles[roleName]
		if role == nil {
			continue
		}
		if _, ok := role.Permissions[perm]; ok {
			return true
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
	allowed := p.rbac.hasPermission(userID, perm)
	if p.audit != nil {
		p.audit.Log(userID, perm, allowed, metadata)
	}
	if !allowed {
		return ErrUnauthorized
	}
	return nil
}
