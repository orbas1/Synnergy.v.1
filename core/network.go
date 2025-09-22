package core

import (
	"encoding/json"
	"errors"
	"sync"
)

// Network manages communication between nodes and relay nodes. It also queues
// transactions for asynchronous broadcasting and integrates biometric
// authentication for secure submission. Additional helpers expose a minimal
// pub‑sub system and lifecycle controls used by the CLI.
type Network struct {
	mu      sync.RWMutex
	nodes   map[string]*Node
	relays  map[string]*Node
	auth    *BiometricService
	queue   chan *Transaction
	running bool
	quit    chan struct{}
	subs    map[string]map[*subscription]struct{}
	events  chan BroadcastEvent
	wg      sync.WaitGroup
}

// NewNetwork creates a network with the provided biometric service. The
// networking loop starts automatically, but can be stopped and restarted via
// Start and Stop.
func NewNetwork(auth *BiometricService) *Network {
	n := &Network{
		nodes:  make(map[string]*Node),
		relays: make(map[string]*Node),
		auth:   auth,
		subs:   make(map[string]map[*subscription]struct{}),
		events: make(chan BroadcastEvent, 256),
	}
	n.Start()
	return n
}

// Start launches background processing if not already running.
func (n *Network) Start() {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.running {
		return
	}
	n.queue = make(chan *Transaction, 100)
	n.quit = make(chan struct{})
	n.running = true
	n.wg.Add(1)
	go n.processQueue()
}

// Stop halts background processing and waits for completion.
func (n *Network) Stop() {
	n.mu.Lock()
	if !n.running {
		n.mu.Unlock()
		return
	}
	close(n.quit)
	n.mu.Unlock()
	n.wg.Wait()
	n.mu.Lock()
	n.running = false
	n.queue = nil
	n.mu.Unlock()
}

// AddNode adds a regular node to the network.
func (n *Network) AddNode(node *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.nodes[node.ID] = node
}

// AddRelay adds a relay node used for extended propagation and redundancy.
func (n *Network) AddRelay(node *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.relays[node.ID] = node
}

// Peers returns the identifiers for all known nodes and relays.
func (n *Network) Peers() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]string, 0, len(n.nodes)+len(n.relays))
	for id := range n.nodes {
		out = append(out, id)
	}
	for id := range n.relays {
		out = append(out, id)
	}
	return out
}

// EnqueueTransaction places a transaction into the broadcast queue.
func (n *Network) EnqueueTransaction(tx *Transaction) error {
	n.mu.RLock()
	running := n.running
	ch := n.queue
	n.mu.RUnlock()
	if !running || ch == nil {
		return errors.New("network not running")
	}
	select {
	case ch <- tx:
		return nil
	default:
		return errors.New("broadcast queue full")
	}
}

// Broadcast verifies biometric data, attaches it to the transaction, and enqueues
// the transaction for network propagation. If biometric verification fails an
// error is returned and the transaction is not broadcast.
func (n *Network) Broadcast(tx *Transaction, userID string, biometric []byte, sig []byte) error {
	if err := tx.AttachBiometric(userID, biometric, sig, n.auth); err != nil {
		return err
	}
	return n.EnqueueTransaction(tx)
}

// Subscribe registers a listener for the given topic and returns a receive-only
// channel. Each call creates an independent buffered channel.
func (n *Network) Subscribe(topic string) (<-chan []byte, func()) {
	sub := newSubscription()
	n.mu.Lock()
	if _, ok := n.subs[topic]; !ok {
		n.subs[topic] = make(map[*subscription]struct{})
	}
	n.subs[topic][sub] = struct{}{}
	n.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			n.mu.Lock()
			if subs, ok := n.subs[topic]; ok {
				delete(subs, sub)
				if len(subs) == 0 {
					delete(n.subs, topic)
				}
			}
			n.mu.Unlock()
			sub.close()
		})
	}

	return sub.channel(), cancel
}

// Publish broadcasts arbitrary data to all subscribers of the provided topic.
// Messages are delivered on a best-effort basis.
func (n *Network) Publish(topic string, data []byte) {
	n.mu.RLock()
	subs := make([]*subscription, 0, len(n.subs[topic]))
	for sub := range n.subs[topic] {
		subs = append(subs, sub)
	}
	n.mu.RUnlock()
	for _, sub := range subs {
		sub.send(data)
	}
}

// processQueue processes queued transactions and broadcasts them to all peers
// and relay nodes. Transactions are propagated in a simple fan-out manner to all
// known nodes.
func (n *Network) processQueue() {
	defer n.wg.Done()
	for {
		select {
		case tx := <-n.queue:
			if tx != nil {
				n.broadcast(tx)
			}
		case <-n.quit:
			return
		}
	}
}

// broadcast sends a transaction to all nodes and relay nodes.
func (n *Network) broadcast(tx *Transaction) {
	n.mu.RLock()
	nodes := make([]*Node, 0, len(n.nodes))
	for _, node := range n.nodes {
		nodes = append(nodes, node)
	}
	relays := make([]*Node, 0, len(n.relays))
	for _, node := range n.relays {
		relays = append(relays, node)
	}
	n.mu.RUnlock()
	for _, node := range nodes {
		err := node.AddTransaction(tx)
		n.emitBroadcastEvent(tx, node.ID, "node", err)
	}
	for _, node := range relays {
		err := node.AddTransaction(tx)
		n.emitBroadcastEvent(tx, node.ID, "relay", err)
	}
}

// BroadcastEvent captures the outcome of attempting to deliver a transaction to
// a network participant.
type BroadcastEvent struct {
	TransactionID string `json:"transaction_id"`
	Target        string `json:"target"`
	Role          string `json:"role"`
	Success       bool   `json:"success"`
	Error         string `json:"error,omitempty"`
}

// BroadcastEvents returns a read-only channel that streams broadcast telemetry.
func (n *Network) BroadcastEvents() <-chan BroadcastEvent {
	return n.events
}

func (n *Network) emitBroadcastEvent(tx *Transaction, target string, role string, err error) {
	event := BroadcastEvent{
		TransactionID: tx.ID,
		Target:        target,
		Role:          role,
		Success:       err == nil,
	}
	if err != nil {
		event.Error = err.Error()
	}
	select {
	case n.events <- event:
	default:
	}
	payload, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		return
	}
	n.Publish("network:broadcast", payload)
}

type subscription struct {
	ch     chan []byte
	mu     sync.Mutex
	closed bool
}

func newSubscription() *subscription {
	return &subscription{ch: make(chan []byte, 1)}
}

func (s *subscription) channel() <-chan []byte {
	return s.ch
}

func (s *subscription) send(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	select {
	case s.ch <- append([]byte(nil), data...):
	default:
	}
}

func (s *subscription) close() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	close(s.ch)
	s.mu.Unlock()
}
