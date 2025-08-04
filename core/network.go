package core

// Network manages communication between nodes.
type Network struct {
	nodes map[string]*Node
}

// NewNetwork creates an empty network.
func NewNetwork() *Network {
	return &Network{nodes: make(map[string]*Node)}
}

// AddNode adds a node to the network.
func (n *Network) AddNode(node *Node) {
	n.nodes[node.ID] = node
}

// Broadcast sends a transaction to all nodes.
func (n *Network) Broadcast(tx *Transaction) {
	for _, node := range n.nodes {
		_ = node.AddTransaction(tx)
	}
}
