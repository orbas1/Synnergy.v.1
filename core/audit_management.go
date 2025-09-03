package core

import (
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
}

// AuditManager coordinates persistent audit logs stored in-memory. In a
// full implementation entries would be persisted to the ledger.
type AuditManager struct {
	mu      sync.RWMutex
	records map[string][]AuditEntry
}

// ErrInvalidAuditEntry is returned when required fields are missing.
var ErrInvalidAuditEntry = errors.New("address and event required")

// NewAuditManager creates a new AuditManager instance.
func NewAuditManager() *AuditManager {
	return &AuditManager{records: make(map[string][]AuditEntry)}
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
