package nodes

// Address mirrors the core address type without creating a dependency.
type Address string

// NodeStatus represents the lifecycle state of a node.
type NodeStatus int

const (
	// NodeStopped indicates the node is not running.
	NodeStopped NodeStatus = iota
	// NodeRunning indicates the node is active.
	NodeRunning
)

// Info captures basic runtime details about a node.  It is useful for
// reporting node lists over the CLI or RPC interfaces.
type Info struct {
	ID     Address
	Status NodeStatus
	Peers  []Address
}
