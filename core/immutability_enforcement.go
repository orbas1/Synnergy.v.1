package core

import "errors"

// ImmutabilityEnforcer ensures the genesis block cannot be altered. It stores
// the expected hash of the genesis block and verifies ledgers against it.
type ImmutabilityEnforcer struct {
	genesisHash string
}

// NewImmutabilityEnforcer creates an enforcer for the supplied genesis block.
func NewImmutabilityEnforcer(genesis *Block) *ImmutabilityEnforcer {
	return &ImmutabilityEnforcer{genesisHash: genesis.Hash}
}

// ErrGenesisChanged is returned when the stored genesis block hash differs from
// the ledger's genesis block.
var (
	ErrNilLedger      = errors.New("ledger is nil")
	ErrGenesisMissing = errors.New("genesis block missing")
	ErrGenesisChanged = errors.New("genesis block hash mismatch")
)

// CheckLedger verifies that the ledger's first block matches the expected
// genesis hash.
func (i *ImmutabilityEnforcer) CheckLedger(l *Ledger) error {
	if l == nil {
		return ErrNilLedger
	}
	b, ok := l.GetBlock(1)
	if !ok {
		return ErrGenesisMissing
	}
	if b.Hash != i.genesisHash {
		return ErrGenesisChanged
	}
	return nil
}
