package core

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errNetworkStopped = errors.New("network not running")
	errNilTransaction = errors.New("transaction required")
)

const networkQueueSize = 512

// TransactionTarget consumes broadcast transactions. *Node already satisfies
// this interface but tests may provide lightweight alternatives.
type TransactionTarget interface {
	AddTransaction(*Transaction) error
}

// Network manages communication between nodes and relay nodes. It also queues
// transactions for asynchronous broadcasting and integrates biometric
// authentication for secure submission. Additional helpers expose a minimal
// pub‑sub system and lifecycle controls used by the CLI.
type Network struct {
	mu             sync.RWMutex
	nodes          map[string]TransactionTarget
	relays         map[string]TransactionTarget
	auth           *BiometricService
	queue          chan queueItem
	running        bool
	quit           chan struct{}
	subs           map[string][]chan []byte // topic -> subscriber channels
	wg             sync.WaitGroup
	retryLimit     int
	retryBackoff   time.Duration
	enqueueTimeout time.Duration
	metrics        networkMetrics
}

type queueItem struct {
	tx       *Transaction
	attempts int
}

type networkMetrics struct {
	enqueued  atomic.Uint64
	delivered atomic.Uint64
	dropped   atomic.Uint64
	failed    atomic.Uint64
	retries   atomic.Uint64
	highWater atomic.Uint64
}

// NetworkMetrics summarises runtime statistics for observability tooling.
type NetworkMetrics struct {
	Enqueued       uint64
	Delivered      uint64
	Dropped        uint64
	Failed         uint64
	Retries        uint64
	QueueHighWater uint64
	QueueDepth     int
	Subscribers    int
}

// NewNetwork creates a network with the provided biometric service. The
// networking loop starts automatically, but can be stopped and restarted via
// Start and Stop.
func NewNetwork(auth *BiometricService) *Network {
	n := &Network{
		nodes:          make(map[string]TransactionTarget),
		relays:         make(map[string]TransactionTarget),
		auth:           auth,
		subs:           make(map[string][]chan []byte),
		retryLimit:     3,
		retryBackoff:   100 * time.Millisecond,
		enqueueTimeout: 500 * time.Millisecond,
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
	n.queue = make(chan queueItem, networkQueueSize)
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
	quit := n.quit
	n.running = false
	n.quit = nil
	n.mu.Unlock()
	close(quit)
	n.wg.Wait()
	n.mu.Lock()
	n.queue = nil
	n.mu.Unlock()
}

// AddNode adds a regular node to the network.
func (n *Network) AddNode(node *Node) {
	if node == nil {
		return
	}
	n.AddTarget(node.ID, node)
}

// AddTarget registers a custom transaction consumer with the provided
// identifier.
func (n *Network) AddTarget(id string, target TransactionTarget) {
	if target == nil || id == "" {
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.nodes[id] = target
}

// AddRelay adds a relay node used for extended propagation and redundancy.
func (n *Network) AddRelay(node *Node) {
	if node == nil {
		return
	}
	n.AddRelayTarget(node.ID, node)
}

// AddRelayTarget registers a custom relay consumer.
func (n *Network) AddRelayTarget(id string, target TransactionTarget) {
	if target == nil || id == "" {
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.relays[id] = target
}

// RemoveNode removes a node by identifier.
func (n *Network) RemoveNode(id string) {
	n.mu.Lock()
	delete(n.nodes, id)
	n.mu.Unlock()
}

// RemoveRelay removes a relay by identifier.
func (n *Network) RemoveRelay(id string) {
	n.mu.Lock()
	delete(n.relays, id)
	n.mu.Unlock()
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
func (n *Network) EnqueueTransaction(tx *Transaction) {
	_ = n.tryEnqueue(queueItem{tx: tx})
}

// Broadcast verifies biometric data, attaches it to the transaction, and enqueues
// the transaction for network propagation. If biometric verification fails an
// error is returned and the transaction is not broadcast.
func (n *Network) Broadcast(tx *Transaction, userID string, biometric []byte, sig []byte) error {
	if tx == nil {
		return errNilTransaction
	}
	if err := tx.AttachBiometric(userID, biometric, sig, n.auth); err != nil {
		return err
	}
	return n.tryEnqueue(queueItem{tx: tx})
}

// Subscribe registers a listener for the given topic and returns a receive-only
// channel. Each call creates an independent buffered channel.
func (n *Network) Subscribe(topic string) <-chan []byte {
	ch := make(chan []byte, 1)
	n.mu.Lock()
	n.subs[topic] = append(n.subs[topic], ch)
	n.mu.Unlock()
	return ch
}

// Publish broadcasts arbitrary data to all subscribers of the provided topic.
// Messages are delivered on a best-effort basis.
func (n *Network) Publish(topic string, data []byte) {
	n.mu.RLock()
	subs := append([]chan []byte(nil), n.subs[topic]...)
	n.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- append([]byte(nil), data...):
		default:
		}
	}
}

// processQueue processes queued transactions and broadcasts them to all peers
// and relay nodes. Transactions are propagated in a simple fan-out manner to all
// known nodes.
func (n *Network) processQueue() {
	defer n.wg.Done()
	for {
		select {
		case item := <-n.queue:
			if item.tx != nil {
				if n.broadcast(item.tx) {
					n.metrics.delivered.Add(1)
				} else {
					n.handleBroadcastFailure(item)
				}
			}
		case <-n.quit:
			return
		}
	}
}

func (n *Network) tryEnqueue(item queueItem) error {
	if item.tx == nil {
		return errNilTransaction
	}
	n.mu.RLock()
	running := n.running
	ch := n.queue
	quit := n.quit
	timeout := n.enqueueTimeout
	n.mu.RUnlock()
	if !running || ch == nil {
		return errNetworkStopped
	}
	if timeout <= 0 {
		select {
		case ch <- item:
			n.metrics.enqueued.Add(1)
			n.updateHighWater(len(ch))
			return nil
		case <-quit:
			return errNetworkStopped
		}
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- item:
		n.metrics.enqueued.Add(1)
		n.updateHighWater(len(ch))
		return nil
	case <-timer.C:
		n.metrics.dropped.Add(1)
		return errors.New("enqueue timeout")
	case <-quit:
		return errNetworkStopped
	}
}

// broadcast sends a transaction to all nodes and relay nodes.
func (n *Network) broadcast(tx *Transaction) {
	nodes, relays := n.snapshotTargets()
	success := true
	for _, node := range nodes {
		if err := node.AddTransaction(tx); err != nil {
			success = false
		}
	}
	for _, relay := range relays {
		if err := relay.AddTransaction(tx); err != nil {
			success = false
		}
	}
	if !success {
		n.metrics.failed.Add(1)
	}
}

func (n *Network) handleBroadcastFailure(item queueItem) {
	if item.attempts >= n.retryLimit {
		return
	}
	retry := queueItem{tx: item.tx, attempts: item.attempts + 1}
	n.metrics.retries.Add(1)
	backoff := n.retryBackoff
	if backoff <= 0 {
		backoff = 50 * time.Millisecond
	}
	backoff = backoff * time.Duration(1<<item.attempts)
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		timer := time.NewTimer(backoff)
		defer timer.Stop()
		select {
		case <-timer.C:
			if err := n.tryEnqueue(retry); err != nil {
				n.metrics.failed.Add(1)
			}
		case <-n.quit:
			return
		}
	}()
}

func (n *Network) snapshotTargets() (nodes []TransactionTarget, relays []TransactionTarget) {
	n.mu.RLock()
	nodes = make([]TransactionTarget, 0, len(n.nodes))
	for _, node := range n.nodes {
		nodes = append(nodes, node)
	}
	relays = make([]TransactionTarget, 0, len(n.relays))
	for _, node := range n.relays {
		relays = append(relays, node)
	}
	n.mu.RUnlock()
	return
}

func (n *Network) updateHighWater(depth int) {
	if depth <= 0 {
		return
	}
	for {
		current := n.metrics.highWater.Load()
		if uint64(depth) <= current {
			return
		}
		if n.metrics.highWater.CompareAndSwap(current, uint64(depth)) {
			return
		}
	}
}

// QueueDepth reports the current number of queued transactions awaiting
// broadcast.
func (n *Network) QueueDepth() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.queue == nil {
		return 0
	}
	return len(n.queue)
}

// SetRetryPolicy adjusts retry behaviour for transient broadcast failures.
func (n *Network) SetRetryPolicy(limit int, backoff time.Duration) {
	n.mu.Lock()
	if limit >= 0 {
		n.retryLimit = limit
	}
	if backoff > 0 {
		n.retryBackoff = backoff
	}
	n.mu.Unlock()
}

// SetEnqueueTimeout configures how long producers wait before enqueueing times
// out.
func (n *Network) SetEnqueueTimeout(timeout time.Duration) {
	n.mu.Lock()
	n.enqueueTimeout = timeout
	n.mu.Unlock()
}

// Metrics returns a snapshot of broadcast statistics.
func (n *Network) Metrics() NetworkMetrics {
	n.mu.RLock()
	subs := 0
	for _, listeners := range n.subs {
		subs += len(listeners)
	}
	depth := 0
	if n.queue != nil {
		depth = len(n.queue)
	}
	n.mu.RUnlock()
	return NetworkMetrics{
		Enqueued:       n.metrics.enqueued.Load(),
		Delivered:      n.metrics.delivered.Load(),
		Dropped:        n.metrics.dropped.Load(),
		Failed:         n.metrics.failed.Load(),
		Retries:        n.metrics.retries.Load(),
		QueueHighWater: n.metrics.highWater.Load(),
		QueueDepth:     depth,
		Subscribers:    subs,
	}
}
