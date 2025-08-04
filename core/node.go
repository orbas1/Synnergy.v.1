package core

// Node represents a participant in the network.
type Node struct {
	ID     string
	Addr   string
	Ledger *Ledger
}

// NewNode creates a new node instance.
func NewNode(id, addr string, ledger *Ledger) *Node {
	return &Node{ID: id, Addr: addr, Ledger: ledger}
}
