package nodes

// NodeInterface defines minimal lifecycle behaviour for specialised nodes.
type NodeInterface interface {
	// ID returns the unique identifier of the node.
	ID() string
	// Start launches any background processing required by the node.
	Start() error
	// Stop terminates background processing and releases resources.
	Stop() error
}
