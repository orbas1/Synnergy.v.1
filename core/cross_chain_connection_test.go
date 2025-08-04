package core

import "testing"

func TestChainConnectionManager(t *testing.T) {
    cm := NewChainConnectionManager()
    id := cm.Open("chainA", "chainB")
    if id == 0 {
        t.Fatalf("expected connection id")
    }
    if err := cm.Close(id); err != nil {
        t.Fatalf("close failed: %v", err)
    }
    c, err := cm.Get(id)
    if err != nil || c.Open {
        t.Fatalf("expected closed connection")
    }
    if len(cm.List()) != 1 {
        t.Fatalf("unexpected connection list length")
    }
}

