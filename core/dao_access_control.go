package core

import "errors"

// AddMember adds a member with a specified role to the DAO.
func (d *DAO) AddMember(addr, role string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.Members == nil {
		d.Members = make(map[string]string)
	}
	if role != "member" && role != "admin" {
		return errors.New("invalid role")
	}
	d.Members[addr] = role
	return nil
}

// RemoveMember deletes a member from the DAO.
func (d *DAO) RemoveMember(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.Members, addr)
}

// MemberRole returns the role for a given address.
func (d *DAO) MemberRole(addr string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	role, ok := d.Members[addr]
	return role, ok
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
