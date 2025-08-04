package nodes

// NodeInterface defines minimal node behaviour independent from core types.
type NodeInterface interface {
	// ID returns the node identifier.
	ID() Address
	// Start begins node operations such as networking routines.
	Start() error
	// Stop gracefully halts node operations.
	Stop() error
	// Peers returns identifiers for all known peers.
	Peers() []Address
	// DialSeed connects the node to a seed peer by address.
	DialSeed(addr Address) error
}
