package security

// DDoSMitigator provides simple tracking for IP addresses.
type DDoSMitigator struct {
	blocked map[string]struct{}
}

// NewDDoSMitigator creates an empty DDoSMitigator.
func NewDDoSMitigator() *DDoSMitigator {
	return &DDoSMitigator{blocked: make(map[string]struct{})}
}

// Block records an address as blocked.
func (d *DDoSMitigator) Block(addr string) { d.blocked[addr] = struct{}{} }

// IsBlocked checks if an address is blocked.
func (d *DDoSMitigator) IsBlocked(addr string) bool {
	_, ok := d.blocked[addr]
	return ok
}
