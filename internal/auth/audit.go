package auth

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

// AuditEvent captures details about an authorization decision.
type AuditEvent struct {
	Timestamp  time.Time      `json:"timestamp"`
	UserID     string         `json:"user"`
	Permission Permission     `json:"permission"`
	Allowed    bool           `json:"allowed"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// AuditLogger records authorization attempts for compliance and forensics.
type AuditLogger interface {
	Log(userID string, perm Permission, allowed bool, metadata map[string]any)
}

// StdAuditLogger is a basic implementation writing JSON entries to a logger.
type StdAuditLogger struct {
	logger *log.Logger
}

// NewStdAuditLogger creates an AuditLogger that writes to the provided writer.
func NewStdAuditLogger(w io.Writer) *StdAuditLogger {
	return &StdAuditLogger{logger: log.New(w, "", log.LstdFlags|log.LUTC)}
}

// Log writes a JSON encoded audit entry.
func (l *StdAuditLogger) Log(userID string, perm Permission, allowed bool, metadata map[string]any) {
	event := AuditEvent{
		Timestamp:  time.Now().UTC(),
		UserID:     userID,
		Permission: perm,
		Allowed:    allowed,
		Metadata:   sanitizeMetadata(metadata),
	}
	data, err := json.Marshal(event)
	if err != nil {
		l.logger.Printf("audit marshal error: %v", err)
		return
	}
	l.logger.Print(string(data))
}

// MemoryAuditLogger stores audit events in memory for testing.
type MemoryAuditLogger struct {
	mu     sync.Mutex
	events []AuditEvent
}

// NewMemoryAuditLogger creates an in-memory audit logger.
func NewMemoryAuditLogger() *MemoryAuditLogger {
	return &MemoryAuditLogger{}
}

// Log stores the event in memory.
func (m *MemoryAuditLogger) Log(userID string, perm Permission, allowed bool, metadata map[string]any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, AuditEvent{
		Timestamp:  time.Now().UTC(),
		UserID:     userID,
		Permission: perm,
		Allowed:    allowed,
		Metadata:   sanitizeMetadata(metadata),
	})
}

// Events returns a copy of the recorded audit events.
func (m *MemoryAuditLogger) Events() []AuditEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]AuditEvent, len(m.events))
	copy(out, m.events)
	return out
}

var sensitiveKeys = []string{"secret", "token", "password", "key"}

func sanitizeMetadata(md map[string]any) map[string]any {
	if len(md) == 0 {
		return nil
	}
	out := make(map[string]any, len(md))
	for k, v := range md {
		if isSensitiveKey(k) {
			out[k] = "***"
			continue
		}
		out[k] = v
	}
	return out
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
