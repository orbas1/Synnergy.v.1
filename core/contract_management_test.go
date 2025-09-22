package core

import (
	"context"
	"sync"
	"testing"
)

type recordingObserver struct {
	mu     sync.Mutex
	events []ContractRegistryEventType
}

func (r *recordingObserver) HandleContractRegistryEvent(event ContractRegistryEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event.Type)
}

func (r *recordingObserver) Types() []ContractRegistryEventType {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]ContractRegistryEventType, len(r.events))
	copy(out, r.events)
	return out
}

func TestContractManager(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	ledger := NewLedger()
	ledger.Credit("owner", 1_000)
	obs := &recordingObserver{}
	reg := NewContractRegistry(vm, ledger, WithContractRegistryObserver(obs))
	addr, err := reg.Deploy([]byte{0x01}, "", 5, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	mgr := NewContractManager(reg)
	if err := mgr.Pause(context.Background(), addr); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if _, _, err := reg.Invoke(addr, "", nil, 5); err == nil {
		t.Fatalf("expected error invoking paused contract")
	}
	if err := mgr.Resume(context.Background(), addr); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if err := mgr.Transfer(context.Background(), addr, "new"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if c, err := mgr.Info(context.Background(), addr); err != nil || c.Owner != "new" {
		t.Fatalf("info mismatch")
	}
	if err := mgr.Upgrade(context.Background(), addr, []byte{0x02}, 6); err != nil {
		t.Fatalf("upgrade: %v", err)
	}

	want := []ContractRegistryEventType{
		ContractRegistryEventDeploy,
		ContractRegistryEventPause,
		ContractRegistryEventResume,
		ContractRegistryEventTransfer,
		ContractRegistryEventUpgrade,
	}
	got := obs.Types()
	if len(got) != len(want) {
		t.Fatalf("unexpected event count %d want %d", len(got), len(want))
	}
	for i, typ := range want {
		if got[i] != typ {
			t.Fatalf("event[%d] = %s, want %s", i, got[i], typ)
		}
	}
}
