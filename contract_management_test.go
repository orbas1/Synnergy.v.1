package synnergy

import "testing"

func TestContractManager(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	reg := NewContractRegistry(vm)
	addr, err := reg.Deploy([]byte{0x01}, "", 5, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	mgr := NewContractManager(reg)
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
}
