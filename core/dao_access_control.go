package core

import "errors"

// AddMember adds a member with a specified role to the DAO.
func (d *DAO) AddMember(addr, role string) error {
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
	delete(d.Members, addr)
}

// MemberRole returns the role for a given address.
func (d *DAO) MemberRole(addr string) (string, bool) {
	role, ok := d.Members[addr]
	return role, ok
}

// MembersList returns all DAO members with their roles.
func (d *DAO) MembersList() map[string]string {
	out := make(map[string]string, len(d.Members))
	for addr, role := range d.Members {
		out[addr] = role
	}
	return out
}
