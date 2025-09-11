package core

import "errors"

// Role constants for DAO membership.
const (
	RoleMember = "member"
	RoleAdmin  = "admin"
)

var (
	errInvalidRole   = errors.New("invalid role")
	errMemberExists  = errors.New("member exists")
	errMemberMissing = errors.New("member not found")
	errUnauthorized  = errors.New("unauthorized")
)

func validRole(role string) bool {
	return role == RoleMember || role == RoleAdmin
}

// AddMember adds a member with a specified role to the DAO.
func (d *DAO) AddMember(addr, role string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.Members == nil {
		d.Members = make(map[string]string)
	}
	if !validRole(role) {
		return errInvalidRole
	}
	if _, ok := d.Members[addr]; ok {
		return errMemberExists
	}
	d.Members[addr] = role
	return nil
}

// UpdateMemberRole updates a member's role. Only an admin can change roles.
func (d *DAO) UpdateMemberRole(requester, addr, role string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.Members[requester] != RoleAdmin {
		return errUnauthorized
	}
	if !validRole(role) {
		return errInvalidRole
	}
	if _, ok := d.Members[addr]; !ok {
		return errMemberMissing
	}
	d.Members[addr] = role
	return nil
}

// RemoveMember deletes a member from the DAO.
func (d *DAO) RemoveMember(addr string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Members[addr]; !ok {
		return errMemberMissing
	}
	delete(d.Members, addr)
	return nil
}

// MemberRole returns the role for a given address.
func (d *DAO) MemberRole(addr string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	role, ok := d.Members[addr]
	return role, ok
}

// IsMember checks if an address is part of the DAO.
func (d *DAO) IsMember(addr string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, ok := d.Members[addr]
	return ok
}

// IsAdmin checks if an address has admin privileges in the DAO.
func (d *DAO) IsAdmin(addr string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Members[addr] == RoleAdmin
}

// MembersList returns all DAO members with their roles.
func (d *DAO) MembersList() map[string]string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make(map[string]string, len(d.Members))
	for addr, role := range d.Members {
		out[addr] = role
	}
	return out
}
