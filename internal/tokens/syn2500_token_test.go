package tokens

import "testing"

func TestSyn2500RegistryAdvanced(t *testing.T) {
	reg := NewSyn2500Registry()
	alice := NewSyn2500Member("alice", "addr1", 10, map[string]string{"kyc": "verified"})
	bob := NewSyn2500Member("bob", "addr2", 5, nil)
	bob.AssignRole("treasurer")

	reg.AddMember(alice)
	reg.AddMember(bob)

	if err := reg.UpdateMetadata("bob", map[string]string{"region": "EU"}); err != nil {
		t.Fatalf("update metadata: %v", err)
	}

	if reg.TotalVotingPower() != 15 {
		t.Fatalf("unexpected voting power: %d", reg.TotalVotingPower())
	}

	treasurers := reg.MembersWithRole("treasurer")
	if len(treasurers) != 1 || treasurers[0].ID != "bob" {
		t.Fatalf("expected bob as treasurer, got %+v", treasurers)
	}

	bob.SetStatus("suspended")
	if reg.TotalVotingPower() != 10 {
		t.Fatalf("voting power should exclude suspended members")
	}

	events := reg.Events(10)
	if len(events) < 3 {
		t.Fatalf("expected events, got %+v", events)
	}

	reg.RemoveMember("bob")
	if len(reg.MembersWithRole("treasurer")) != 0 {
		t.Fatalf("role lookup should be empty after removal")
	}
}
