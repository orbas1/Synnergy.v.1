package core

import "testing"

func TestConsensusNetworkManager(t *testing.T) {
	m := NewConsensusNetworkManager()
	if _, err := m.RegisterNetwork("pos", "pow", "rel1"); err == nil {
		t.Fatalf("expected unauthorized register")
	}
	m.AuthorizeRelayer("rel1")
	id, err := m.RegisterNetwork("pos", "pow", "rel1")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	n, err := m.GetNetwork(id)
	if err != nil || n.SourceConsensus != "pos" {
		t.Fatalf("get network failed: %#v err=%v", n, err)
	}
	if len(m.ListNetworks()) != 1 {
		t.Fatalf("expected one network")
	}
	if err := m.RemoveNetwork(id, "bad"); err == nil {
		t.Fatalf("expected unauthorized removal")
	}
	if err := m.RemoveNetwork(id, "rel1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, err := m.GetNetwork(id); err == nil {
		t.Fatalf("expected error after removal")
	}
}
