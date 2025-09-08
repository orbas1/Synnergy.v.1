package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// AuditEntry records a single audit log event.
type AuditEntry struct {
	Address   string
	Event     string
	Metadata  map[string]string
	Timestamp time.Time
	Signature []byte
}

// AuditManager coordinates persistent audit logs stored in-memory. In a
// full implementation entries would be persisted to the ledger.
type AuditManager struct {
	mu      sync.RWMutex
	records map[string][]AuditEntry
	pub     ed25519.PublicKey
	priv    ed25519.PrivateKey
}

// ErrInvalidAuditEntry is returned when required fields are missing.
var ErrInvalidAuditEntry = errors.New("address and event required")
var ErrSigningFailed = errors.New("unable to sign audit entry")

// NewAuditManager creates a new AuditManager instance.
func NewAuditManager() *AuditManager {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return &AuditManager{records: make(map[string][]AuditEntry), pub: pub, priv: priv}
}

// NewAuditManagerFromKey allows creating an AuditManager with a predetermined key.
func NewAuditManagerFromKey(priv ed25519.PrivateKey) *AuditManager {
	var pub ed25519.PublicKey
	if priv != nil {
		pub = priv.Public().(ed25519.PublicKey)
	}
	return &AuditManager{records: make(map[string][]AuditEntry), pub: pub, priv: priv}
}

// Log records an audit event for the given address.
func (m *AuditManager) Log(address, event string, metadata map[string]string) error {
	if address == "" || event == "" {
		return ErrInvalidAuditEntry
	}
	if metadata == nil {
		metadata = make(map[string]string)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	entry := AuditEntry{
		Address:   address,
		Event:     event,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
	payload, err := json.Marshal(struct {
		Address   string
		Event     string
		Metadata  map[string]string
		Timestamp time.Time
	}{entry.Address, entry.Event, entry.Metadata, entry.Timestamp})
	if err != nil {
		return ErrSigningFailed
	}
	entry.Signature = ed25519.Sign(m.priv, payload)
	m.records[address] = append(m.records[address], entry)
	return nil
}

// List returns all audit events for the given address.
func (m *AuditManager) List(address string) []AuditEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entries := m.records[address]
	// return copy to prevent external mutation
	out := make([]AuditEntry, len(entries))
	copy(out, entries)
	return out
}

// Verify checks the integrity of an audit entry using its signature.
func (m *AuditManager) Verify(entry AuditEntry) bool {
	payload, err := json.Marshal(struct {
		Address   string
		Event     string
		Metadata  map[string]string
		Timestamp time.Time
	}{entry.Address, entry.Event, entry.Metadata, entry.Timestamp})
	if err != nil {
		return false
	}
	return ed25519.Verify(m.pub, payload, entry.Signature)
}

// PublicKey returns the manager's public key for external verification.
func (m *AuditManager) PublicKey() ed25519.PublicKey {
	return m.pub
}
