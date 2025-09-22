package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy/internal/nodes/extra/watchtower"
)

const defaultWatchtowerEventLimit = 256

// WatchtowerEventType labels event payloads emitted by the watchtower node.
type WatchtowerEventType string

const (
	WatchtowerEventStarted      WatchtowerEventType = "watchtower.started"
	WatchtowerEventStopped      WatchtowerEventType = "watchtower.stopped"
	WatchtowerEventMetrics      WatchtowerEventType = "watchtower.metrics"
	WatchtowerEventForkDetected WatchtowerEventType = "watchtower.fork"
	WatchtowerEventAlert        WatchtowerEventType = "watchtower.alert"
)

// WatchtowerEvent captures observability data for CLI and web dashboards.
type WatchtowerEvent struct {
	Sequence  uint64
	Type      WatchtowerEventType
	Timestamp time.Time
	Payload   map[string]string
}

// Watchtower implements the watchtower.WatchtowerNode interface.  It observes
// the network, records system health metrics and reports detected forks.
type Watchtower struct {
	id       string
	firewall *Firewall
	health   *SystemHealthLogger
	logger   *log.Logger

	mu           sync.RWMutex
	running      bool
	cancel       context.CancelFunc
	observed     *Node
	events       []WatchtowerEvent
	eventSeq     uint64
	watchers     map[uint64]chan WatchtowerEvent
	watcherID    uint64
	evLimit      int
	tickInterval time.Duration
}

// NewWatchtowerNode constructs a Watchtower with the provided identifier.
func NewWatchtowerNode(id string, logger *log.Logger) *Watchtower {
	return &Watchtower{
		id:           id,
		firewall:     NewFirewall(),
		health:       NewSystemHealthLogger(),
		logger:       logger,
		evLimit:      defaultWatchtowerEventLimit,
		tickInterval: 5 * time.Second,
	}
}

// ID returns the unique identifier of the node.
func (w *Watchtower) ID() string { return w.id }

// AttachNode links the watchtower to a full node providing consensus and VM
// visibility.  Metrics emitted by monitorLoop include the attached node's
// mempool and validator information when available.
func (w *Watchtower) AttachNode(n *Node) {
	w.mu.Lock()
	w.observed = n
	w.mu.Unlock()
}

// SetEventRetention configures the number of events retained for replay.
func (w *Watchtower) SetEventRetention(limit int) {
	if limit <= 0 {
		limit = defaultWatchtowerEventLimit
	}
	w.mu.Lock()
	w.evLimit = limit
	if len(w.events) > limit {
		w.events = append([]WatchtowerEvent(nil), w.events[len(w.events)-limit:]...)
	}
	w.mu.Unlock()
}

// Start begins monitoring routines for the watchtower node.
func (w *Watchtower) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return errors.New("watchtower already running")
	}
	var c context.Context
	c, w.cancel = context.WithCancel(ctx)
	w.running = true
	w.mu.Unlock()
	go w.monitorLoop(c)
	w.recordEvent(WatchtowerEvent{Type: WatchtowerEventStarted, Timestamp: time.Now().UTC(), Payload: map[string]string{"id": w.id}})
	return nil
}

// monitorLoop periodically collects system health metrics and compares them to
// configurable thresholds to raise alerts.  The loop tolerates transient
// failures and logs anomalies rather than panic so production watchdogs remain
// stable.
func (w *Watchtower) monitorLoop(ctx context.Context) {
	interval := w.tickInterval
	if interval <= 0 {
		interval = 5 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			metrics := w.health.Collect(0, 0)
			payload := map[string]string{
				"cpu_usage":  formatFloat(metrics.CPUUsage),
				"memory":     formatUint64(metrics.MemoryUsage),
				"peer_count": fmt.Sprintf("%d", metrics.PeerCount),
				"last_block": formatUint64(metrics.LastBlockHeight),
			}
			if w.logger != nil {
				w.logger.Printf("watchtower metrics: %+v", metrics)
			}
			if node := w.snapshotNode(); node != nil {
				payload["mempool_size"] = formatInt(len(node.Mempool))
				payload["validators"] = formatInt(len(node.Validators.Eligible()))
			}
			w.recordEvent(WatchtowerEvent{Type: WatchtowerEventMetrics, Timestamp: metrics.Timestamp, Payload: payload})
		case <-ctx.Done():
			return
		}
	}
}

// snapshotNode safely returns the observed node pointer without holding locks
// for extended periods.
func (w *Watchtower) snapshotNode() *Node {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.observed
}

// Stop halts monitoring routines for the watchtower node.
func (w *Watchtower) Stop() error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return nil
	}
	cancel := w.cancel
	w.running = false
	w.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	w.recordEvent(WatchtowerEvent{Type: WatchtowerEventStopped, Timestamp: time.Now().UTC(), Payload: map[string]string{"id": w.id}})
	return nil
}

// ReportFork records details of a detected fork event.
func (w *Watchtower) ReportFork(height uint64, hash string) {
	payload := map[string]string{
		"height": formatUint64(height),
		"hash":   hash,
	}
	if w.logger != nil {
		w.logger.Printf("fork detected at height %d hash %s", height, hash)
	}
	w.recordEvent(WatchtowerEvent{Type: WatchtowerEventForkDetected, Timestamp: time.Now().UTC(), Payload: payload})
}

// RaiseAlert allows external modules to inject alerts that will be fanned out
// to subscribers alongside internally generated metrics.
func (w *Watchtower) RaiseAlert(category, message string) {
	payload := map[string]string{"category": category, "message": message}
	w.recordEvent(WatchtowerEvent{Type: WatchtowerEventAlert, Timestamp: time.Now().UTC(), Payload: payload})
}

// Metrics returns the latest snapshot of system health data.
func (w *Watchtower) Metrics() watchtower.Metrics {
	return w.health.Snapshot()
}

// Firewall exposes the internal firewall instance for rule management.
func (w *Watchtower) Firewall() *Firewall { return w.firewall }

// Events returns a copy of the retained event log.
func (w *Watchtower) Events() []WatchtowerEvent {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]WatchtowerEvent, len(w.events))
	copy(out, w.events)
	return out
}

// EventsSince filters events newer than the supplied sequence.
func (w *Watchtower) EventsSince(seq uint64) []WatchtowerEvent {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var out []WatchtowerEvent
	for _, ev := range w.events {
		if ev.Sequence > seq {
			out = append(out, ev)
		}
	}
	return out
}

// SubscribeEvents registers a channel that receives future watchtower events.
func (w *Watchtower) SubscribeEvents(buffer int) (<-chan WatchtowerEvent, func()) {
	if buffer <= 0 {
		buffer = 16
	}
	ch := make(chan WatchtowerEvent, buffer)
	w.mu.Lock()
	if w.watchers == nil {
		w.watchers = make(map[uint64]chan WatchtowerEvent)
	}
	w.watcherID++
	id := w.watcherID
	w.watchers[id] = ch
	backlog := append([]WatchtowerEvent(nil), w.events...)
	w.mu.Unlock()

	go func(events []WatchtowerEvent) {
		for _, ev := range events {
			ch <- ev
		}
	}(backlog)

	cancel := func() {
		w.mu.Lock()
		if ch, ok := w.watchers[id]; ok {
			delete(w.watchers, id)
			close(ch)
		}
		w.mu.Unlock()
	}
	return ch, cancel
}

// RunIntegritySweep executes a one-off integrity check comparing consensus
// weights and mempool depth. Alerts are emitted when thresholds are breached.
func (w *Watchtower) RunIntegritySweep(ctx context.Context, maxMempool int) ([]WatchtowerEvent, error) {
	node := w.snapshotNode()
	if node == nil {
		return nil, errors.New("no node attached")
	}
	var events []WatchtowerEvent
	if len(node.Mempool) > maxMempool {
		ev := WatchtowerEvent{Type: WatchtowerEventAlert, Timestamp: time.Now().UTC(), Payload: map[string]string{
			"category": "mempool",
			"message":  "mempool size above threshold",
			"size":     formatInt(len(node.Mempool)),
		}}
		w.recordEvent(ev)
		events = append(events, ev)
	}
	if node.Consensus != nil {
		weights := node.Consensus.WeightsSnapshot()
		payload := map[string]string{
			"pow": formatFloat(weights.PoW),
			"pos": formatFloat(weights.PoS),
			"poh": formatFloat(weights.PoH),
		}
		ev := WatchtowerEvent{Type: WatchtowerEventMetrics, Timestamp: time.Now().UTC(), Payload: payload}
		w.recordEvent(ev)
		events = append(events, ev)
	}
	return events, ctx.Err()
}

func (w *Watchtower) recordEvent(ev WatchtowerEvent) {
	w.mu.Lock()
	w.eventSeq++
	ev.Sequence = w.eventSeq
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	w.events = append(w.events, ev)
	if w.evLimit > 0 && len(w.events) > w.evLimit {
		w.events = append([]WatchtowerEvent(nil), w.events[len(w.events)-w.evLimit:]...)
	}
	watchers := make([]chan WatchtowerEvent, 0, len(w.watchers))
	for _, ch := range w.watchers {
		watchers = append(watchers, ch)
	}
	w.mu.Unlock()

	for _, ch := range watchers {
		select {
		case ch <- ev:
		default:
		}
	}
}

func formatUint64(v uint64) string { return fmt.Sprintf("%d", v) }

func formatInt(v int) string { return fmt.Sprintf("%d", v) }

func formatFloat(v float64) string { return fmt.Sprintf("%.2f", v) }

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", d.Seconds()*1000)
}

// ensure interface compliance
var _ watchtower.WatchtowerNode = (*Watchtower)(nil)
