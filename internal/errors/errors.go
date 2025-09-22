package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Code identifies a class of application error.
type Code string

// Severity qualifies how critical an error is for routing and alerting.
type Severity string

const (
	// NotFound indicates a missing resource.
	NotFound Code = "not_found"
	// Invalid indicates input validation failure.
	Invalid Code = "invalid"
	// Internal marks unexpected internal errors.
	Internal Code = "internal"
	// Unauthorized indicates the caller failed authentication/authorization.
	Unauthorized Code = "unauthorized"
	// Conflict indicates a state conflict (optimistic locking, replay, etc.).
	Conflict Code = "conflict"
)

const (
	// SeverityInfo represents a non-critical informational error.
	SeverityInfo Severity = "info"
	// SeverityWarn signals an issue that should be surfaced to operators.
	SeverityWarn Severity = "warn"
	// SeverityError marks high priority failures that require attention.
	SeverityError Severity = "error"
)

// Error wraps an underlying error with structured metadata suitable for the CLI
// renderer, REST handlers and telemetry pipeline.
type Error struct {
	Code          Code           `json:"code"`
	Message       string         `json:"message"`
	Err           error          `json:"-"`
	Severity      Severity       `json:"severity"`
	CorrelationID string         `json:"correlation_id,omitempty"`
	Status        int            `json:"status"`
	Retryable     bool           `json:"retryable"`
	Timestamp     time.Time      `json:"timestamp"`
	Details       map[string]any `json:"details,omitempty"`
}

// Option mutates an error during construction.
type Option func(*Error)

// New creates a coded error without an underlying cause.
func New(code Code, msg string, opts ...Option) *Error {
	return Wrap(code, msg, nil, opts...)
}

// Wrap attaches a code and message to an existing error and applies options.
func Wrap(code Code, msg string, err error, opts ...Option) *Error {
	e := &Error{
		Code:      code,
		Message:   msg,
		Err:       err,
		Severity:  SeverityError,
		Status:    defaultStatus(code),
		Timestamp: time.Now().UTC(),
		Details:   map[string]any{},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error for errors.Is/As.
func (e *Error) Unwrap() error { return e.Err }

// IsCode reports whether err is an *Error with the given code.
func IsCode(err error, code Code) bool {
	var coded *Error
	if errors.As(err, &coded) {
		return coded.Code == code
	}
	return false
}

// ToMap returns a representation suitable for CLI and JSON output.
func (e *Error) ToMap() map[string]any {
	if e == nil {
		return nil
	}
	cloned := map[string]any{
		"code":      e.Code,
		"message":   e.Message,
		"severity":  e.Severity,
		"status":    e.Status,
		"retryable": e.Retryable,
		"timestamp": e.Timestamp.Format(time.RFC3339Nano),
	}
	if e.CorrelationID != "" {
		cloned["correlation_id"] = e.CorrelationID
	}
	if len(e.Details) > 0 {
		cloned["details"] = e.Details
	}
	return cloned
}

// MarshalJSON ensures the wrapped error is omitted while the structured fields
// remain serialisable.
func (e *Error) MarshalJSON() ([]byte, error) {
	type alias Error
	return json.Marshal((*alias)(e))
}

// WithCorrelation assigns a correlation identifier.
func WithCorrelation(id string) Option {
	return func(e *Error) {
		e.CorrelationID = id
	}
}

// WithSeverity sets a custom severity.
func WithSeverity(sev Severity) Option {
	return func(e *Error) {
		e.Severity = sev
	}
}

// WithStatus overrides the HTTP status code mapping.
func WithStatus(status int) Option {
	return func(e *Error) {
		e.Status = status
	}
}

// WithRetryable toggles the retryable flag.
func WithRetryable(retryable bool) Option {
	return func(e *Error) {
		e.Retryable = retryable
	}
}

// WithDetail adds a key/value pair to the error metadata.
func WithDetail(key string, value any) Option {
	return func(e *Error) {
		if e.Details == nil {
			e.Details = map[string]any{}
		}
		e.Details[key] = value
	}
}

// Parse attempts to unwrap an error back into the structured Error type.
func Parse(err error) (*Error, bool) {
	var coded *Error
	if errors.As(err, &coded) {
		return coded, true
	}
	return nil, false
}

func defaultStatus(code Code) int {
	switch code {
	case NotFound:
		return http.StatusNotFound
	case Invalid:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
