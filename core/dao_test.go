package core

import "testing"

func TestDAOManager(t *testing.T) {
	mgr := NewDAOManager()
	if _, err := mgr.Create("TestDAO", "creator"); err == nil {
		t.Fatalf("expected unauthorized creator")
	}
	mgr.AuthorizeRelayer("creator")
	dao, err := mgr.Create("TestDAO", "creator")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if dao.Name != "TestDAO" {
		t.Fatalf("unexpected name")
	}

	if err := mgr.Join(dao.ID, "user1"); err == nil {
		t.Fatalf("expected unauthorized join")
	}
	mgr.AuthorizeRelayer("user1")
	if err := mgr.Join(dao.ID, "user1"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if info, err := mgr.Info(dao.ID); err != nil || len(info.Members) != 2 {
		t.Fatalf("info: %v members=%d", err, len(info.Members))
	}
	if daos := mgr.List(); len(daos) != 1 {
		t.Fatalf("expected one dao, got %d", len(daos))
	}

	mgr.RevokeRelayer("user1")
	if err := mgr.Leave(dao.ID, "user1"); err == nil {
		t.Fatalf("expected unauthorized leave")
	}
	mgr.AuthorizeRelayer("user1")
	if err := mgr.Leave(dao.ID, "user1"); err != nil {
		t.Fatalf("leave: %v", err)
	}
	if err := mgr.Leave(dao.ID, "user1"); err == nil {
		t.Fatalf("expected error for unknown member")
	}
}
