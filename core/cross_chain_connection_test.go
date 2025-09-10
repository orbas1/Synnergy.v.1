package core

import "testing"

func TestChainConnectionManager(t *testing.T) {
	cm := NewChainConnectionManager()
	id := cm.Open("chainA", "chainB", "rel1")
	if id == 0 {
		t.Fatalf("expected connection id")
	}
	if !cm.IsRelayerAuthorized(id, "rel1") {
		t.Fatalf("relayer should be authorized")
	}
	if err := cm.Close(id, "rel2"); err == nil {
		t.Fatalf("expected unauthorized close failure")
	}
	if err := cm.AuthorizeRelayer(id, "rel2"); err != nil {
		t.Fatalf("authorize failed: %v", err)
	}
	if !cm.IsRelayerAuthorized(id, "rel2") {
		t.Fatalf("rel2 should now be authorized")
	}
	if err := cm.Close(id, "rel2"); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	c, err := cm.Get(id)
	if err != nil || c.Open {
		t.Fatalf("expected closed connection")
	}
	if err := cm.Remove(id); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if _, err := cm.Get(id); err == nil {
		t.Fatalf("expected connection to be removed")
	}
}
