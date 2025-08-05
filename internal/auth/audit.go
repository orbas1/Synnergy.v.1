package auth

import (
	"encoding/json"
	"io"
	"log"
	"time"
)

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
	entry := map[string]any{
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
		"user":       userID,
		"permission": string(perm),
		"allowed":    allowed,
		"metadata":   metadata,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf("audit marshal error: %v", err)
		return
	}
	l.logger.Print(string(data))
}
