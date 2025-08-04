package core

import "strings"

// AddressZero represents the zero-value address (all 20 bytes set to zero).
const AddressZero = "0x0000000000000000000000000000000000000000"

// IsZeroAddress returns true if the provided address equals AddressZero.
func IsZeroAddress(addr string) bool {
	return strings.ToLower(addr) == AddressZero
}
