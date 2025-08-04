package core

import "testing"

func TestDAOAccessControl(t *testing.T) {
	dao := &DAO{Members: make(map[string]string)}
	if err := dao.AddMember("addr1", "member"); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if role, ok := dao.MemberRole("addr1"); !ok || role != "member" {
		t.Fatalf("unexpected role %v %v", role, ok)
	}
	dao.RemoveMember("addr1")
	if _, ok := dao.MemberRole("addr1"); ok {
		t.Fatalf("expected removal")
	}
}
