package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Ledger maintains account balances and block history. It persists blocks to a
// simple write-ahead log so that a node can recover on restart. The WAL path is
// optional; if empty the ledger operates purely in memory.
// UTXO represents an unspent transaction output owned by an address.
type UTXO struct {
	ID     string `json:"id"`
	Amount uint64 `json:"amount"`
}

// Ledger maintains account balances, a simple UTXO view and block history. It
// persists blocks to a simple write-ahead log so that a node can recover on
// restart. The WAL path is optional; if empty the ledger operates purely in
// memory.
type Ledger struct {
	mu       sync.RWMutex
	balances map[string]uint64
	blocks   []*Block
	walPath  string
	utxos    map[string][]*UTXO
	mempool  []*Transaction
	nextUTXO uint64
}

// NewLedger creates a new ledger. If a path is supplied it will replay any
// existing WAL file to restore previous blocks.
func NewLedger(path ...string) *Ledger {
	l := &Ledger{
		balances: make(map[string]uint64),
		blocks:   []*Block{},
		utxos:    make(map[string][]*UTXO),
		mempool:  []*Transaction{},
	}
	if len(path) > 0 {
		l.walPath = path[0]
		l.replayWAL()
	}
	return l
}

// replayWAL loads blocks from the write-ahead log if configured. Errors are
// ignored which keeps recovery best-effort.
func (l *Ledger) replayWAL() {
	if l.walPath == "" {
		return
	}
	f, err := os.Open(l.walPath)
	if err != nil {
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for {
		var b Block
		if err := dec.Decode(&b); err != nil {
			break
		}
		l.blocks = append(l.blocks, &b)
	}
}

// appendWAL writes a block to the WAL if a path is configured.
func (l *Ledger) appendWAL(b *Block) {
	if l.walPath == "" {
		return
	}
	f, err := os.OpenFile(l.walPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(b)
}

// Head returns the current height and hash of the latest block.
func (l *Ledger) Head() (int, string) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	height := len(l.blocks)
	if height == 0 {
		return 0, ""
	}
	return height, l.blocks[height-1].Hash
}

// GetBlock returns the block at the provided 1-indexed height.
func (l *Ledger) GetBlock(height int) (*Block, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if height <= 0 || height > len(l.blocks) {
		return nil, false
	}
	return l.blocks[height-1], true
}

// AddBlock appends a block to the chain and persists it to the WAL.
func (l *Ledger) AddBlock(b *Block) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.blocks = append(l.blocks, b)
	l.appendWAL(b)
}

// GetBalance returns the balance for a given address.
func (l *Ledger) GetBalance(addr string) uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.balances[addr]
}

// GetUTXOs returns a copy of the unspent outputs for an address.
func (l *Ledger) GetUTXOs(addr string) []UTXO {
	l.mu.RLock()
	defer l.mu.RUnlock()
	outs := l.utxos[addr]
	res := make([]UTXO, len(outs))
	for i, u := range outs {
		res[i] = *u
	}
	return res
}

// Credit adds funds to an address and updates the UTXO view.
func (l *Ledger) Credit(addr string, amount uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.balances[addr] += amount
	l.updateUTXO(addr)
}

// Mint is an exported helper that credits funds to an address. It mirrors the
// CLI `mint` command.
func (l *Ledger) Mint(addr string, amount uint64) {
	l.Credit(addr, amount)
}

// Transfer moves funds from one address to another. A fee is deducted from the
// sender. It returns an error if the sender lacks sufficient balance.
func (l *Ledger) Transfer(from, to string, amount, fee uint64) error {
	tx := NewTransaction(from, to, amount, fee, 0)
	return l.ApplyTransaction(tx)
}

// ApplyTransaction applies a transaction to the ledger, deducting both amount
// and fee from the sender. It returns an error if the sender lacks sufficient
// funds.
func (l *Ledger) ApplyTransaction(tx *Transaction) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	total := uint64(tx.Amount + tx.Fee)
	if l.balances[tx.From] < total {
		return errors.New("insufficient funds")
	}
	l.balances[tx.From] -= total
	l.balances[tx.To] += uint64(tx.Amount)
	l.updateUTXO(tx.From)
	l.updateUTXO(tx.To)
	return nil
}

// AddToPool appends a transaction to the mem-pool.
func (l *Ledger) AddToPool(tx *Transaction) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.mempool = append(l.mempool, tx)
}

// Pool returns a snapshot of the current mem-pool transactions.
func (l *Ledger) Pool() []*Transaction {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]*Transaction, len(l.mempool))
	copy(out, l.mempool)
	return out
}

func (l *Ledger) newUTXO(amount uint64) *UTXO {
	u := &UTXO{ID: fmt.Sprintf("u%d", l.nextUTXO), Amount: amount}
	l.nextUTXO++
	return u
}

func (l *Ledger) updateUTXO(addr string) {
	if l.balances[addr] == 0 {
		delete(l.utxos, addr)
		return
	}
	l.utxos[addr] = []*UTXO{l.newUTXO(l.balances[addr])}
}
