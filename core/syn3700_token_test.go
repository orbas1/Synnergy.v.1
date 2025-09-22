package core

import (
	"sync"
	"testing"
)

func TestSYN3700GovernanceLifecycle(t *testing.T) {
	controller := Address("ctrl")
	token := NewSYN3700Token("Index", "IDX")
	if err := token.AddController(controller, controller); err != nil {
		t.Fatalf("add controller: %v", err)
	}
	if err := token.AddComponent("AAA", 0.5, 0.1, controller); err != nil {
		t.Fatalf("add component: %v", err)
	}
	if err := token.UpdateComponent("AAA", 0.75, 0.15, controller); err != nil {
		t.Fatalf("update component: %v", err)
	}
	if err := token.RecordRebalance("AAA", 0.75, controller); err != nil {
		t.Fatalf("record rebalance: %v", err)
	}
	if exceeded := token.DriftExceeded(); len(exceeded) != 0 {
		t.Fatalf("unexpected drift exceeded: %v", exceeded)
	}
	tele := token.Telemetry()
	if tele.ControllerCount != 1 || tele.ComponentCount != 1 {
		t.Fatalf("unexpected telemetry: %+v", tele)
	}
	snap := token.Snapshot()
	restored := NewSYN3700Token("", "")
	restored.Restore(snap)
	if len(restored.ListComponents()) != 1 {
		t.Fatalf("expected component restored")
	}
	plan := restored.RebalancePlan()
	if got := plan["AAA"]; got[0] != got[1] {
		t.Fatalf("rebalance mismatch: %v", got)
	}
}

func TestSYN3700ConcurrentMutations(t *testing.T) {
	controller := Address("ctrl")
	token := NewSYN3700Token("Index", "IDX")
	if err := token.AddController(controller, controller); err != nil {
		t.Fatalf("controller: %v", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := string(rune('A' + i))
			if err := token.AddComponent(name, 0.1, 0.2, controller); err != nil && err != ErrComponentExists {
				t.Errorf("add component: %v", err)
			}
		}(i)
	}
	wg.Wait()
	comps := token.ListComponents()
	if len(comps) == 0 {
		t.Fatalf("expected components to be registered")
	}
}

func TestSYN3700SnapshotRoundTrip(t *testing.T) {
	controller := Address("ctrl")
	token := NewSYN3700Token("Index", "IDX")
	if err := token.AddController(controller, controller); err != nil {
		t.Fatalf("controller: %v", err)
	}
	if err := token.AddComponent("AAA", 0.5, 0.1, controller); err != nil {
		t.Fatalf("component: %v", err)
	}
	snap := token.Snapshot()
	clone := NewSYN3700Token("", "")
	clone.Restore(snap)
	if tele := clone.Telemetry(); tele.ComponentCount != 1 || tele.ControllerCount != 1 {
		t.Fatalf("telemetry mismatch: %+v", tele)
	}
	if err := clone.RemoveComponent("AAA", controller); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if len(clone.ListComponents()) != 0 {
		t.Fatalf("component removal failed")
	}
}
