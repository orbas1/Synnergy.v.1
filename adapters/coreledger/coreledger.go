package coreledger

import (
	"synnergy"
	"synnergy/core"
)

// Wrap adapts a core ledger to the synnergy contract registry interface.
func Wrap(l *core.Ledger) synnergy.Ledger {
	if l == nil {
		return nil
	}
	return &adapter{ledger: l}
}

type adapter struct {
	ledger *core.Ledger
}

func (a *adapter) Transfer(from, to string, amount, fee uint64) error {
	return a.ledger.Transfer(from, to, amount, fee)
}

func (a *adapter) StoreContract(rec synnergy.LedgerContractRecord) {
	stored := make([]byte, len(rec.WASM))
	copy(stored, rec.WASM)
	a.ledger.RegisterContract(core.LedgerContract{
		Address:  rec.Address,
		Owner:    rec.Owner,
		Manifest: rec.Manifest,
		GasLimit: rec.GasLimit,
		WASM:     stored,
	})
}

func (a *adapter) Contracts() []synnergy.LedgerContractRecord {
	records := a.ledger.Contracts()
	out := make([]synnergy.LedgerContractRecord, 0, len(records))
	for _, rec := range records {
		wasm := make([]byte, len(rec.WASM))
		copy(wasm, rec.WASM)
		out = append(out, synnergy.LedgerContractRecord{
			Address:  rec.Address,
			Owner:    rec.Owner,
			Manifest: rec.Manifest,
			GasLimit: rec.GasLimit,
			WASM:     wasm,
		})
	}
	return out
}
