package core

import (
	"encoding/hex"
	"errors"
	"strings"
)

// Address represents a 20-byte hex encoded account identifier.
// It is stored as a string with a leading 0x prefix to keep things simple
// and compatible with other parts of the codebase that treat addresses as
// strings.
type Address string

// StringToAddress converts a hex string into an Address. The string must be
// 40 hex characters prefixed by 0x.
func StringToAddress(s string) (Address, error) {
	if !strings.HasPrefix(s, "0x") || len(s) != 42 {
		return "", errors.New("invalid address")
	}
	if _, err := hex.DecodeString(s[2:]); err != nil {
		return "", err
	}
	return Address(strings.ToLower(s)), nil
}

// Hex returns the string form of the address.
func (a Address) Hex() string { return string(a) }

// Bytes returns the raw 20 byte representation of the address.
func (a Address) Bytes() []byte {
	b, _ := hex.DecodeString(strings.TrimPrefix(string(a), "0x"))
	return b
}



// Address represents a 20-byte hex encoded identifier for accounts and contracts.
// It is stored as a string with 0x prefix for simplicity. The methods provided
// here offer convenient conversions and formatting helpers used throughout the
// codebase.
type Address string

// StringToAddress converts a hex string (with 0x prefix) into an Address. It
// validates the input ensuring it is 40 hex characters long. An error is
// returned if the input is malformed.
func StringToAddress(s string) (Address, error) {
	s = strings.ToLower(strings.TrimPrefix(s, "0x"))
	if len(s) != 40 {
		return "", fmt.Errorf("invalid address length")
	}
	if _, err := hex.DecodeString(s); err != nil {
		return "", fmt.Errorf("invalid address: %w", err)
	}
	return Address("0x" + s), nil
}

// Hex returns the canonical hexadecimal representation of the address.
func (a Address) Hex() string { return string(a) }

// Bytes returns the raw 20-byte slice of the address. If the address is
// malformed, an empty slice is returned.
func (a Address) Bytes() []byte {
	h := strings.TrimPrefix(string(a), "0x")
	b, err := hex.DecodeString(h)
	if err != nil {
		return []byte{}
	}
	return b
}

// Short returns an abbreviated form of the address suitable for logging. It
// shows the first 6 and last 4 characters of the address.
func (a Address) Short() string {
	h := a.Hex()
	if len(h) <= 10 {
		return h
	}
	return h[:10]
}

// Hash represents a simple 32-byte value used for identifiers.
type Hash [32]byte
	return h[:6] + "..." + h[len(h)-4:]
}

// String implements fmt.Stringer and returns the hexadecimal representation.
func (a Address) String() string { return a.Hex() }

