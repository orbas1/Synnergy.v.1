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

// Short returns a shortened representation useful for logs.
func (a Address) Short() string {
	h := a.Hex()
	if len(h) <= 10 {
		return h
	}
	return h[:10]
}

// Hash represents a simple 32-byte value used for identifiers.
type Hash [32]byte
