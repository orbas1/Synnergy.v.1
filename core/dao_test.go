package core

import "testing"

func TestDAOManager(t *testing.T) {
	mgr := NewDAOManager()
	dao := mgr.Create("TestDAO", "creator")
	if dao.Name != "TestDAO" {
		t.Fatalf("unexpected name")
	}
	if err := mgr.Join(dao.ID, "user1"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := mgr.Leave(dao.ID, "user1"); err != nil {
		t.Fatalf("leave: %v", err)
	}
	if err := mgr.Leave(dao.ID, "user1"); err == nil {
		t.Fatalf("expected error for unknown member")
	}
}
