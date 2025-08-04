package core

import "errors"

// Node represents a participant in the network.
type Node struct {
	ID         string
	Addr       string
	Ledger     *Ledger
	Consensus  *SynnergyConsensus
	VM         *SNVM
	Mempool    []*Transaction
	Blockchain []*Block
	Stakes     map[string]uint64
}

// NewNode creates a new node instance.
func NewNode(id, addr string, ledger *Ledger) *Node {
	return &Node{ID: id, Addr: addr, Ledger: ledger, Consensus: NewSynnergyConsensus(), VM: NewSNVM(), Stakes: make(map[string]uint64)}
}

// AddTransaction validates and adds a transaction to the mempool.
func (n *Node) AddTransaction(tx *Transaction) error {
	if err := n.ValidateTransaction(tx); err != nil {
		return err
	}
	n.Mempool = append(n.Mempool, tx)
	return nil
}

// ValidateTransaction checks if a transaction is well-formed and the sender has
// sufficient balance.
func (n *Node) ValidateTransaction(tx *Transaction) error {
	if n.Ledger.GetBalance(tx.From) < tx.Amount+tx.Fee {
		return errors.New("insufficient funds")
	}
	return nil
}

// MineBlock packages the current mempool into a sub-block and mines a block.
func (n *Node) MineBlock() *Block {
	if len(n.Mempool) == 0 {
		return nil
	}
	validator := n.Consensus.SelectValidator(n.Stakes)
	sb := NewSubBlock(n.Mempool, validator)
	n.Mempool = nil
	prevHash := ""
	if len(n.Blockchain) > 0 {
		prevHash = n.Blockchain[len(n.Blockchain)-1].Hash
	}
	block := NewBlock([]*SubBlock{sb}, prevHash)
	n.Consensus.MineBlock(block, 3)
	for _, tx := range sb.Transactions {
		_ = n.Ledger.ApplyTransaction(tx)
	}
	n.Blockchain = append(n.Blockchain, block)
	return block
}

// SetStake assigns stake to an address for validator selection.
func (n *Node) SetStake(addr string, amount uint64) {
	n.Stakes[addr] = amount
}
