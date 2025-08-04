package core

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Address represents a 20-byte hexadecimal account identifier.
type Address string

// StringToAddress converts a hexadecimal string into an Address. It performs a
// light validation on length and hex characters.
func StringToAddress(s string) (Address, error) {
	if len(s) != 42 || !strings.HasPrefix(s, "0x") {
		return "", fmt.Errorf("invalid address format")
	}
	if _, err := hex.DecodeString(s[2:]); err != nil {
		return "", err
	}
	return Address(strings.ToLower(s)), nil
}

// Hex returns the hexadecimal representation of the address.
func (a Address) Hex() string { return string(a) }

// Bytes returns the byte slice representation of the address string.
func (a Address) Bytes() []byte { return []byte(a) }

// Short returns an abbreviated form useful for logging (e.g., 0x1234…abcd).
func (a Address) Short() string {
	s := a.Hex()
	if len(s) <= 10 {
		return s
	}
	return s[:6] + "…" + s[len(s)-4:]
}
