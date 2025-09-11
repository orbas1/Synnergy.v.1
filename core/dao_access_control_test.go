package core

import "testing"

func TestDAOAccessControl(t *testing.T) {
	dao := &DAO{Members: make(map[string]string)}
	if err := dao.AddMember("addr1", RoleMember); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := dao.AddMember("addr1", RoleMember); err == nil {
		t.Fatalf("expected duplicate error")
	}
	if role, ok := dao.MemberRole("addr1"); !ok || role != RoleMember {
		t.Fatalf("unexpected role %v %v", role, ok)
	}
	if err := dao.UpdateMemberRole("addr1", "addr1", RoleAdmin); err == nil {
		t.Fatalf("non-admin should not update role")
	}
	// grant admin rights to requester
	if err := dao.AddMember("admin", RoleAdmin); err != nil {
		t.Fatalf("add admin: %v", err)
	}
	if err := dao.UpdateMemberRole("admin", "addr1", RoleAdmin); err != nil {
		t.Fatalf("update role: %v", err)
	}
	if !dao.IsAdmin("addr1") {
		t.Fatalf("expected addr1 to be admin")
	}
	if err := dao.RemoveMember("addr1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, ok := dao.MemberRole("addr1"); ok {
		t.Fatalf("expected removal")
	}
	if err := dao.RemoveMember("addr1"); err == nil {
		t.Fatalf("expected error on missing member")
	}
}
