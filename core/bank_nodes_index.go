package core

import (
	"encoding/json"
	"sync"
)

// Bank node type identifiers.
const (
	BankInstitutionalNodeType = "BANK_INSTITUTIONAL"
	CentralBankingNodeType    = "CENTRAL_BANKING"
	CustodialNodeType         = "CUSTODIAL"
)

// BankNodeTypes lists all bank-related node categories supported by the network.
var BankNodeTypes = []string{
	BankInstitutionalNodeType,
	CentralBankingNodeType,
	CustodialNodeType,
}

// BankNodeRecord holds minimal information about a bank node for indexing.
type BankNodeRecord struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// BankNodeIndex provides thread-safe access to bank node records.
type BankNodeIndex struct {
	mu    sync.RWMutex
	nodes map[string]*BankNodeRecord
}

// NewBankNodeIndex returns an initialised BankNodeIndex.
func NewBankNodeIndex() *BankNodeIndex {
	return &BankNodeIndex{nodes: make(map[string]*BankNodeRecord)}
}

// Add inserts or replaces a bank node record.
func (idx *BankNodeIndex) Add(rec *BankNodeRecord) {
	idx.mu.Lock()
	idx.nodes[rec.ID] = rec
	idx.mu.Unlock()
}

// Get retrieves a record by its identifier.
func (idx *BankNodeIndex) Get(id string) (*BankNodeRecord, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	rec, ok := idx.nodes[id]
	return rec, ok
}

// Remove deletes a record from the index.
func (idx *BankNodeIndex) Remove(id string) {
	idx.mu.Lock()
	delete(idx.nodes, id)
	idx.mu.Unlock()
}

// List returns all indexed bank node records.
func (idx *BankNodeIndex) List() []*BankNodeRecord {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make([]*BankNodeRecord, 0, len(idx.nodes))
	for _, n := range idx.nodes {
		out = append(out, n)
	}
	return out
}

// Snapshot returns a copy of the current index map for deterministic replication.
func (idx *BankNodeIndex) Snapshot() map[string]*BankNodeRecord {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	cp := make(map[string]*BankNodeRecord, len(idx.nodes))
	for id, n := range idx.nodes {
		cp[id] = n
	}
	return cp
}

// MarshalJSON serialises the index deterministically.
func (idx *BankNodeIndex) MarshalJSON() ([]byte, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return json.Marshal(idx.List())
}
