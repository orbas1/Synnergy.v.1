package synnergy

import "sync"

// AccessController manages role-based access permissions for addresses.
type AccessController struct {
	mu    sync.RWMutex
	roles map[string]map[string]struct{}
}

// NewAccessController constructs a new AccessController instance.
func NewAccessController() *AccessController {
	return &AccessController{roles: make(map[string]map[string]struct{})}
}

// Grant assigns a role to an address.
func (a *AccessController) Grant(role, addr string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.roles[addr]; !ok {
		a.roles[addr] = make(map[string]struct{})
	}
	a.roles[addr][role] = struct{}{}
}

// Revoke removes a role from an address.
func (a *AccessController) Revoke(role, addr string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if roles, ok := a.roles[addr]; ok {
		delete(roles, role)
		if len(roles) == 0 {
			delete(a.roles, addr)
		}
	}
}

// HasRole checks whether an address has a specific role.
func (a *AccessController) HasRole(role, addr string) bool {
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
func (a *AccessController) List(addr string) []string {
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

// Audit snapshots all address-role assignments for inspection without
// mutating internal state. The returned map is a copy and safe for
// concurrent use by callers.
func (a *AccessController) Audit() map[string][]string {
        a.mu.RLock()
        defer a.mu.RUnlock()
        snapshot := make(map[string][]string, len(a.roles))
        for addr, roles := range a.roles {
                list := make([]string, 0, len(roles))
                for r := range roles {
                        list = append(list, r)
                }
                snapshot[addr] = list
        }
        return snapshot
}
