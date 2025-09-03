package errors

import "fmt"

// Code identifies a class of application error.
type Code string

const (
	// NotFound indicates a missing resource.
	NotFound Code = "not_found"
	// Invalid indicates input validation failure.
	Invalid Code = "invalid"
	// Internal marks unexpected internal errors.
	Internal Code = "internal"
)

// Error wraps an underlying error with a code and message.
type Error struct {
	Code    Code
	Message string
	Err     error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error for errors.Is/As.
func (e *Error) Unwrap() error { return e.Err }

// New creates a coded error without an underlying cause.
func New(code Code, msg string) *Error { return &Error{Code: code, Message: msg} }

// Wrap attaches a code and message to an existing error.
func Wrap(code Code, msg string, err error) *Error {
	return &Error{Code: code, Message: msg, Err: err}
}

// IsCode reports whether err is an *Error with the given code.
func IsCode(err error, code Code) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}
