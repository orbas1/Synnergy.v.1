package core

import (
    "encoding/hex"
    "fmt"
    "strings"
)

// Address represents a 20-byte hex encoded identifier for accounts and contracts.
// It is stored as a string with a 0x prefix for simplicity.
type Address string

// StringToAddress converts a hex string (with or without the 0x prefix) into an
// Address. An error is returned if the string is not exactly 40 hexadecimal
// characters.
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
// malformed an empty slice is returned.
func (a Address) Bytes() []byte {
    b, err := hex.DecodeString(strings.TrimPrefix(string(a), "0x"))
    if err != nil {
        return []byte{}
    }
    return b
}

// Short returns an abbreviated form of the address suitable for logging.
// The first six and last four characters are preserved.
func (a Address) Short() string {
    h := a.Hex()
    if len(h) <= 10 {
        return h
    }
    return h[:6] + "..." + h[len(h)-4:]
}

// String implements fmt.Stringer.
func (a Address) String() string { return a.Hex() }

// Hash represents a simple 32-byte value used for identifiers.
type Hash [32]byte

