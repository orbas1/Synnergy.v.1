package core

import "testing"

// TestDAOStaking verifies staking is gated by membership and tracks balances
// per DAO.
func TestDAOStaking(t *testing.T) {
	mgr := NewDAOManager()
	mgr.AuthorizeRelayer("admin")
	dao, err := mgr.Create("dao", "admin")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	s := NewDAOStaking(mgr)
	if err := s.Stake(dao.ID, "admin", 10); err != nil {
		t.Fatalf("stake: %v", err)
	}
	// Non-member cannot stake
	if err := s.Stake(dao.ID, "bob", 5); err == nil {
		t.Fatalf("expected error for non-member stake")
	}
	if err := s.Unstake(dao.ID, "admin", 5); err != nil {
		t.Fatalf("unstake: %v", err)
	}
	if bal := s.Balance(dao.ID, "admin"); bal != 5 {
		t.Fatalf("expected 5, got %d", bal)
	}
	if total := s.TotalStaked(dao.ID); total != 5 {
		t.Fatalf("unexpected total %d", total)
	}
}
