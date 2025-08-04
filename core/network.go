package core

// Network manages communication between nodes and relay nodes. It also queues
// transactions for asynchronous broadcasting and integrates biometric
// authentication for secure submission.

type Network struct {
	nodes  map[string]*Node
	relays map[string]*Node
	auth   *BiometricService
	queue  chan *Transaction
}

// NewNetwork creates a network with the provided biometric service. A
// background goroutine processes the transaction queue to broadcast
// transactions to peers and relay nodes.
func NewNetwork(auth *BiometricService) *Network {
	n := &Network{
		nodes:  make(map[string]*Node),
		relays: make(map[string]*Node),
		auth:   auth,
		queue:  make(chan *Transaction, 100),
	}
	go n.processQueue()
	return n
}

// AddNode adds a regular node to the network.
func (n *Network) AddNode(node *Node) { n.nodes[node.ID] = node }

// AddRelay adds a relay node used for extended propagation and redundancy.
func (n *Network) AddRelay(node *Node) { n.relays[node.ID] = node }

// EnqueueTransaction places a transaction into the broadcast queue.
func (n *Network) EnqueueTransaction(tx *Transaction) { n.queue <- tx }

// Broadcast verifies biometric data, attaches it to the transaction, and enqueues
// the transaction for network propagation. If biometric verification fails an
// error is returned and the transaction is not broadcast.
func (n *Network) Broadcast(tx *Transaction, userID string, biometric []byte) error {
	if err := tx.AttachBiometric(userID, biometric, n.auth); err != nil {
		return err
	}
	n.EnqueueTransaction(tx)
	return nil
}

// processQueue processes queued transactions and broadcasts them to all peers
// and relay nodes. Transactions are propagated in a simple fan-out manner to all
// known nodes.
func (n *Network) processQueue() {
	for tx := range n.queue {
		n.broadcast(tx)
	}
}

// broadcast sends a transaction to all nodes and relay nodes.
func (n *Network) broadcast(tx *Transaction) {
	for _, node := range n.nodes {
		_ = node.AddTransaction(tx)
	}
	for _, node := range n.relays {
		_ = node.AddTransaction(tx)
	}
}
