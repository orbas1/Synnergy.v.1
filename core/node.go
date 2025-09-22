package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Node represents a participant in the network.
type Node struct {
	ID            string
	Addr          string
	Ledger        *Ledger
	Consensus     *SynnergyConsensus
	VM            *SNVM
	Mempool       []*Transaction
	Blockchain    []*Block
	Validators    *ValidatorManager
	MaxTxPerBlock int
	mu            sync.Mutex
}

// NewNode creates a new node instance.
func NewNode(id, addr string, ledger *Ledger) *Node {
	return &Node{
		ID:            id,
		Addr:          addr,
		Ledger:        ledger,
		Consensus:     NewSynnergyConsensus(),
		VM:            NewSNVM(),
		Mempool:       []*Transaction{},
		Blockchain:    []*Block{},
		Validators:    NewValidatorManager(MinStake),
		MaxTxPerBlock: 100,
	}
}

// AddTransaction validates and adds a transaction to the mempool.
func (n *Node) AddTransaction(tx *Transaction) error {
	if err := n.ValidateTransaction(tx); err != nil {
		return err
	}
	n.mu.Lock()
	n.Mempool = append(n.Mempool, tx)
	n.mu.Unlock()
	return nil
}

// ValidateTransaction checks if a transaction is well-formed and the sender has
// sufficient balance.
func (n *Node) ValidateTransaction(tx *Transaction) error {
	// Ensure the fee is considered with the amount using explicit uint64
	// arithmetic. This guards against future changes to transaction field
	// types that might otherwise introduce float arithmetic.
	if n.Ledger.GetBalance(tx.From) < uint64(tx.Amount+tx.Fee) {
		return errors.New("insufficient funds")
	}
	return nil
}

// MineBlock packages the current mempool into a sub-block and mines a block.
func (n *Node) MineBlock() *Block {
	n.mu.Lock()
	defer n.mu.Unlock()
	if len(n.Mempool) == 0 {
		return nil
	}
	prevHash := ""
	if len(n.Blockchain) > 0 {
		prevHash = n.Blockchain[len(n.Blockchain)-1].Hash
	}
	eligible := n.eligibleStakes()
	validator := n.Consensus.SelectValidator(prevHash, eligible)
	if validator == "" {
		return nil
	}
	sb := NewSubBlock(n.Mempool, validator)
	if !n.Consensus.ValidateSubBlock(sb) {
		return nil
	}
	n.Mempool = nil
	block := NewBlock([]*SubBlock{sb}, prevHash)
	n.Consensus.MineBlock(block, 3)
	votes := make(map[string]bool, len(eligible))
	for addr := range eligible {
		votes[addr] = true
	}
	if len(votes) == 0 {
		votes[validator] = true
	}
	n.Consensus.FinalizeBlock(block, votes, n.Validators, 1)
	var totalFees uint64
	for _, tx := range sb.Transactions {
		totalFees += tx.Fee
		_ = n.Ledger.ApplyTransaction(tx)
	}
	n.Blockchain = append(n.Blockchain, block)

	dist := DistributeFees(totalFees)
	pool := AdjustForBlockUtilization(dist.ValidatorsMiners, len(sb.Transactions), n.MaxTxPerBlock)
	weights := map[string]uint64{validator: n.Validators.Stake(validator)}
	weights[n.ID] = 1
	shares := ShareProportional(pool, weights)
	contract := NewFeeDistributionContract(n.Ledger)
	contract.Distribute(shares)
	return block
}

const MinStake uint64 = 1

// SetStake assigns stake to an address for validator selection while enforcing a minimum.
func (n *Node) SetStake(addr string, amount uint64) error {
	if amount < MinStake {
		return fmt.Errorf("stake below minimum: %d", amount)
	}
	return n.Validators.Add(context.Background(), addr, amount)
}

func (n *Node) eligibleStakes() map[string]uint64 {
	return n.Validators.Eligible()
}

// ReportDoubleSign slashes a validator for double signing.
func (n *Node) ReportDoubleSign(addr string) {
	n.Validators.SlashWithEvidence(context.Background(), addr, "double-sign")
}

// ReportDowntime slashes a validator for downtime.
func (n *Node) ReportDowntime(addr string) {
	n.Validators.SlashWithEvidence(context.Background(), addr, "downtime")
}

// Rehabilitate removes slashed status from a validator by re-adding the existing stake.
func (n *Node) Rehabilitate(addr string) {
	stake := n.Validators.Stake(addr)
	if stake > 0 {
		_ = n.Validators.Add(context.Background(), addr, stake)
	}
}
