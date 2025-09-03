package core

import (
	"errors"
	"sync"
)

// ErrEmptyRole is returned when no role name is supplied.
var ErrEmptyRole = errors.New("role cannot be empty")

// AccessController manages role based access permissions using validated
// addresses. All operations are safe for concurrent use.
type AccessController struct {
	mu    sync.RWMutex
	roles map[Address]map[string]struct{}
}

// NewAccessController constructs a new AccessController instance.
func NewAccessController() *AccessController {
	return &AccessController{roles: make(map[Address]map[string]struct{})}
}

// Grant assigns a role to an address after basic validation.
func (a *AccessController) Grant(addr Address, role string) error {
	if role == "" {
		return ErrEmptyRole
	}
	if addr.IsZero() {
		return ErrInvalidAddress
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.roles[addr]; !ok {
		a.roles[addr] = make(map[string]struct{})
	}
	a.roles[addr][role] = struct{}{}
	return nil
}

// Revoke removes a role from an address. It is a no-op for unknown addresses.
func (a *AccessController) Revoke(addr Address, role string) error {
	if role == "" {
		return ErrEmptyRole
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if roles, ok := a.roles[addr]; ok {
		delete(roles, role)
		if len(roles) == 0 {
			delete(a.roles, addr)
		}
	}
	return nil
}

// Has reports whether an address holds a role.
func (a *AccessController) Has(addr Address, role string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	roles, ok := a.roles[addr]
	if !ok {
		return false
	}
	_, ok = roles[role]
	return ok
}

// List returns all roles assigned to an address.
func (a *AccessController) List(addr Address) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	roles, ok := a.roles[addr]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(roles))
	for r := range roles {
		out = append(out, r)
	}
	return out
}

// GlobalAccess is the package level access controller used by default.
var GlobalAccess = NewAccessController()

// GrantRole validates inputs and stores the role assignment.
func GrantRole(role, addr string) error {
	a, err := StringToAddress(addr)
	if err != nil {
		return err
	}
	return GlobalAccess.Grant(a, role)
}

// RevokeRole removes a role from an address if present.
func RevokeRole(role, addr string) error {
	a, err := StringToAddress(addr)
	if err != nil {
		return err
	}
	return GlobalAccess.Revoke(a, role)
}

// HasRole reports whether an address owns a given role.
func HasRole(role, addr string) (bool, error) {
	if role == "" {
		return false, ErrEmptyRole
	}
	a, err := StringToAddress(addr)
	if err != nil {
		return false, err
	}
	return GlobalAccess.Has(a, role), nil
}

// ListRoles returns all roles for an address.
func ListRoles(addr string) ([]string, error) {
	a, err := StringToAddress(addr)
	if err != nil {
		return nil, err
	}
	return GlobalAccess.List(a), nil
}

// ErrInvalidAddress is returned when a provided address fails validation.
var ErrInvalidAddress = errors.New("invalid address")
