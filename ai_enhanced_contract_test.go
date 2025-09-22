package synnergy

import (
	"errors"
	"testing"
)

type testLedger struct {
	balances  map[string]uint64
	contracts map[string]LedgerContractRecord
}

func newTestLedger() *testLedger {
	return &testLedger{
		balances:  make(map[string]uint64),
		contracts: make(map[string]LedgerContractRecord),
	}
}

func (l *testLedger) Transfer(from, to string, amount, fee uint64) error {
	total := amount + fee
	if from != "" {
		if l.balances[from] < total {
			return errors.New("insufficient funds")
		}
		l.balances[from] -= total
	}
	l.balances[to] += amount
	return nil
}

func (l *testLedger) StoreContract(rec LedgerContractRecord) {
	stored := make([]byte, len(rec.WASM))
	copy(stored, rec.WASM)
	rec.WASM = stored
	l.contracts[rec.Address] = rec
}

func (l *testLedger) Contracts() []LedgerContractRecord {
	out := make([]LedgerContractRecord, 0, len(l.contracts))
	for _, rec := range l.contracts {
		wasm := make([]byte, len(rec.WASM))
		copy(wasm, rec.WASM)
		out = append(out, LedgerContractRecord{
			Address:  rec.Address,
			Owner:    rec.Owner,
			Manifest: rec.Manifest,
			GasLimit: rec.GasLimit,
			WASM:     wasm,
		})
	}
	return out
}

func TestAIContractRegistry(t *testing.T) {
	vm := NewSimpleVM()
	_ = vm.Start()
	ledger := newTestLedger()
	ledger.balances["owner"] = 1_000
	reg := NewContractRegistry(vm, ledger)
	aiReg := NewAIContractRegistry(reg)
	addr, err := aiReg.DeployAIContract([]byte{0x01}, "model", "", 5, "owner")
	if err != nil {
		t.Fatalf("deploy: %v", err)
	}
	if h, ok := aiReg.ModelHash(addr); !ok || h != "model" {
		t.Fatalf("model hash mismatch")
	}
	if _, _, err := aiReg.InvokeAIContract(addr, []byte("in"), 5); err != nil {
		t.Fatalf("invoke: %v", err)
	}
}
