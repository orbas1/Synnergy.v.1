package synnergy_test

import (
	"sync"
	"testing"

	synnergy "synnergy"
	"synnergy/adapters/coreledger"
	"synnergy/core"
)

type recordingObserver struct {
	mu     sync.Mutex
	events []synnergy.ContractRegistryEventType
}

func (r *recordingObserver) HandleContractRegistryEvent(event synnergy.ContractRegistryEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event.Type)
}

func (r *recordingObserver) Types() []synnergy.ContractRegistryEventType {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]synnergy.ContractRegistryEventType, len(r.events))
	copy(out, r.events)
	return out
}

func TestContractManager(t *testing.T) {
	vm := synnergy.NewSimpleVM()
	_ = vm.Start()
	ledger := core.NewLedger()
	ledger.Credit("owner", 50)
	ledger.Credit("new", 50)
	obs := &recordingObserver{}
	reg := synnergy.NewContractRegistry(vm, coreledger.Wrap(ledger), synnergy.WithContractRegistryObserver(obs))
	addr, err := reg.Deploy([]byte{0x01}, "", 5, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	mgr := synnergy.NewContractManager(reg)
	if err := mgr.Pause(addr); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if _, _, err := reg.Invoke(addr, "", nil, 5); err == nil {
		t.Fatalf("expected error invoking paused contract")
	}
	if err := mgr.Resume(addr); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if err := mgr.Transfer(addr, "new"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if c, err := mgr.Info(addr); err != nil || c.Owner != "new" {
		t.Fatalf("info mismatch")
	}
	if err := mgr.Upgrade(addr, []byte{0x02}, 6); err != nil {
		t.Fatalf("upgrade: %v", err)
	}

	want := []synnergy.ContractRegistryEventType{
		synnergy.ContractRegistryEventDeploy,
		synnergy.ContractRegistryEventPause,
		synnergy.ContractRegistryEventResume,
		synnergy.ContractRegistryEventTransfer,
		synnergy.ContractRegistryEventUpgrade,
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
