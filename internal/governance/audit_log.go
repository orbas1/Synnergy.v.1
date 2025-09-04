package governance

import "sync"

// AuditLog stores simple string entries.
type AuditLog struct {
	mu      sync.Mutex
	entries []string
}

// NewAuditLog creates an AuditLog.
func NewAuditLog() *AuditLog { return &AuditLog{} }

// Append adds an entry to the log.
func (a *AuditLog) Append(e string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, e)
}

// Entries returns a copy of log entries.
func (a *AuditLog) Entries() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]string, len(a.entries))
	copy(out, a.entries)
	return out
}
