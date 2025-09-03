package core

// AddressZero represents the zero-value address (all 20 bytes set to zero).
const AddressZero Address = "0x0000000000000000000000000000000000000000"

// IsZeroAddress parses the supplied string and reports whether it is the zero
// address. Invalid strings return false.
func IsZeroAddress(addr string) bool {
	a, err := StringToAddress(addr)
	if err != nil {
		return false
	}
	return a.IsZero()
}
