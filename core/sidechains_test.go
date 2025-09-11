package core

import "testing"

func TestSidechainRegistry(t *testing.T) {
	reg := NewSidechainRegistry()
	if _, err := reg.Register("chain1", "meta", []string{"val1"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.SubmitHeader("chain1", "h1"); err != nil {
		t.Fatalf("submit header: %v", err)
	}
	if h, ok := reg.GetHeader("chain1"); !ok || h != "h1" {
		t.Fatalf("get header failed")
	}
	if err := reg.Pause("chain1"); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if err := reg.Resume("chain1"); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if err := reg.UpdateValidators("chain1", []string{"val2"}); err != nil {
		t.Fatalf("update validators: %v", err)
	}
	if reg.chains["chain1"].Validators[0] != "val2" {
		t.Fatalf("validators not updated")
	}
	if len(reg.List()) != 1 {
		t.Fatalf("list length")
	}
	if err := reg.Remove("chain1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if len(reg.List()) != 0 {
		t.Fatalf("remove failed")
	}
}

func TestSidechainEscrow(t *testing.T) {
	reg := NewSidechainRegistry()
	if _, err := reg.Register("s", "m", nil); err != nil {
		t.Fatal(err)
	}
	if err := reg.Deposit("s", "alice", 10); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	if bal, err := reg.Balance("s", "alice"); err != nil || bal != 10 {
		t.Fatalf("balance mismatch: %v %d", err, bal)
	}
	if err := reg.Withdraw("s", "alice", 5); err != nil {
		t.Fatalf("withdraw: %v", err)
	}
	if bal, _ := reg.Balance("s", "alice"); bal != 5 {
		t.Fatalf("unexpected balance %d", bal)
	}
	// paused chains reject deposits
	if err := reg.Pause("s"); err != nil {
		t.Fatal(err)
	}
	if err := reg.Deposit("s", "alice", 1); err == nil {
		t.Fatalf("expected deposit error on paused chain")
	}
}

// Test concurrent deposits for race-safety
func TestSidechainConcurrentDeposits(t *testing.T) {
	reg := NewSidechainRegistry()
	if _, err := reg.Register("s", "m", nil); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		go func() {
			if err := reg.Deposit("s", "bob", 1); err != nil {
				t.Errorf("deposit: %v", err)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if bal, _ := reg.Balance("s", "bob"); bal != 10 {
		t.Fatalf("expected 10, got %d", bal)
	}
}
