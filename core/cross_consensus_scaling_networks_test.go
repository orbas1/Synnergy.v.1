package core

import "testing"

func TestConsensusNetworkManager(t *testing.T) {
	m := NewConsensusNetworkManager()
	id := m.RegisterNetwork("pos", "pow")
	if id == 0 {
		t.Fatalf("expected id")
	}
	n, err := m.GetNetwork(id)
	if err != nil || n.SourceConsensus != "pos" {
		t.Fatalf("get network failed: %#v err=%v", n, err)
	}
	if len(m.ListNetworks()) != 1 {
		t.Fatalf("expected one network")

	}
}
