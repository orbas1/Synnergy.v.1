//go:build experimental

package nodes

// ExperimentalNode demonstrates a feature that is only compiled when the
// "experimental" build tag is supplied. It implements NodeInterface but
// performs no real networking; the functionality can evolve independently of
// stable releases.
type ExperimentalNode struct {
	id      Address
	running bool
}

// NewExperimentalNode constructs a new experimental node with the provided identifier.
func NewExperimentalNode(id Address) *ExperimentalNode {
	return &ExperimentalNode{id: id}
}

// ID returns the node identifier.
func (n *ExperimentalNode) ID() Address { return n.id }

// Start marks the node as running.
func (n *ExperimentalNode) Start() error {
	n.running = true
	return nil
}

// Stop halts the node.
func (n *ExperimentalNode) Stop() error {
	n.running = false
	return nil
}

// IsRunning reports whether the node is active.
func (n *ExperimentalNode) IsRunning() bool { return n.running }

// Peers returns an empty slice as the node does not connect to peers yet.
func (n *ExperimentalNode) Peers() []Address { return nil }

// DialSeed is a stub for connecting to a seed peer.
func (n *ExperimentalNode) DialSeed(addr Address) error { return nil }
