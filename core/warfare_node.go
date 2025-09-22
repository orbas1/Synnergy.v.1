package core

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	militarynodes "synnergy/internal/nodes/military_nodes"
)

const (
	commanderRootID      = "root"
	defaultEventLogLimit = 512
	commandSkewTolerance = 5 * time.Minute
)

// WarfareEventType categorises emitted warfare node events so that CLI and web
// consoles can render tailored views without additional parsing.
type WarfareEventType string

const (
	WarfareEventCommandAccepted   WarfareEventType = "command.accepted"
	WarfareEventCommandRejected   WarfareEventType = "command.rejected"
	WarfareEventLogisticsRecorded WarfareEventType = "logistics.recorded"
	WarfareEventLogisticsRejected WarfareEventType = "logistics.rejected"
	WarfareEventTacticalBroadcast WarfareEventType = "tactical.broadcast"
	WarfareEventVMTrace           WarfareEventType = "vm.trace"
)

// WarfareEvent is broadcast whenever privileged commands are executed, when
// logistics are updated or when VM hooks fire.  Events are replayable so that
// delayed subscribers – including the JavaScript control panel – can build an
// accurate state view.
type WarfareEvent struct {
	Sequence  uint64
	Type      WarfareEventType
	Timestamp time.Time
	Commander string
	Payload   map[string]string
}

// CommandRequest wraps a privileged command along with metadata used for
// signature verification and replay protection.
type CommandRequest struct {
	Commander string
	Command   string
	Timestamp time.Time
	Nonce     uint64
	Signature []byte
	Metadata  map[string]string
}

// CommandRecord captures the result of executing a privileged command.
type CommandRecord struct {
	Commander string
	Command   string
	Timestamp time.Time
	Nonce     uint64
	Metadata  map[string]string
	Accepted  bool
	Error     string
	Consensus ConsensusSnapshot
	Latency   time.Duration
}

// ConsensusSnapshot exposes the consensus weights at the time a command was
// processed so auditors can confirm policy decisions.
type ConsensusSnapshot struct {
	PoW float64
	PoS float64
	PoH float64
}

// CommanderCredential returns the private and public keys associated with an
// authorised commander.  The private key is presented once so operators can
// store it securely for future signatures.
type CommanderCredential struct {
	ID         string
	PublicKey  string
	PrivateKey string
	IssuedAt   time.Time
}

// LogisticsUpdate represents a pending logistics mutation including optional
// metadata used for compliance attestation.
type LogisticsUpdate struct {
	AssetID   string
	Location  string
	Status    string
	Reporter  string
	Metadata  map[string]string
	Timestamp time.Time
}

type commanderEntry struct {
	pub       ed25519.PublicKey
	lastNonce uint64
}

type warfareMetrics struct {
	commands       uint64
	failedCommands uint64
	logistics      uint64
	tactical       uint64
}

type vmHookRegistrar interface {
	RegisterHook(ExecutionHook)
}

type vmCallable interface {
	Call(string) error
}

// WarfareNode provides military focused extensions on top of a base Node.  It
// maintains audit logs, command authorisations and event subscriptions used by
// the CLI and the browser-based control panel.
type WarfareNode struct {
	*Node

	mu         sync.RWMutex
	logistics  []militarynodes.LogisticsRecord
	commandLog []CommandRecord
	events     []WarfareEvent

	commanders  map[string]commanderEntry
	rootKey     ed25519.PrivateKey
	eventSeq    uint64
	subscriber  uint64
	subscribers map[uint64]chan WarfareEvent

	metrics warfareMetrics
	evLimit int
}

// NewWarfareNode wraps a base node with warfare specific functionality and
// instruments the attached virtual machine so execution traces are streamed to
// subscribers.  A root commander is generated automatically and used by the CLI
// when no external credentials are supplied.
func NewWarfareNode(base *Node) *WarfareNode {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Errorf("generate commander key: %w", err))
	}

	wn := &WarfareNode{
		Node:       base,
		commanders: map[string]commanderEntry{commanderRootID: {pub: pub}},
		rootKey:    priv,
		evLimit:    defaultEventLogLimit,
	}
	if base != nil && base.VM != nil {
		if hooker, ok := interface{}(base.VM).(vmHookRegistrar); ok {
			hooker.RegisterHook(func(trace ExecutionTrace) {
				payload := map[string]string{
					"opcode":   fmt.Sprintf("0x%06X", trace.Opcode),
					"name":     trace.Name,
					"duration": trace.Duration.String(),
					"gas":      fmt.Sprintf("%d", trace.GasCost),
				}
				if trace.Err != nil {
					payload["error"] = trace.Err.Error()
				}
				wn.recordEvent(WarfareEvent{
					Type:      WarfareEventVMTrace,
					Timestamp: time.Now().UTC(),
					Payload:   payload,
				})
			})
		}
	}
	return wn
}

// SetEventLogLimit adjusts the maximum number of events retained for replay.
// The default of 512 aligns with the JavaScript control panel pagination.
func (w *WarfareNode) SetEventLogLimit(limit int) {
	if limit <= 0 {
		limit = defaultEventLogLimit
	}
	w.mu.Lock()
	w.evLimit = limit
	if len(w.events) > limit {
		w.events = append([]WarfareEvent(nil), w.events[len(w.events)-limit:]...)
	}
	w.mu.Unlock()
}

// GetID satisfies militarynodes.BaseNode.
func (w *WarfareNode) GetID() string {
	if w.Node != nil {
		return w.Node.ID
	}
	return ""
}

// IssueCommander generates a new commander key pair, registers the public key
// and returns the credential so the private key can be stored in an HSM or
// secure enclave.
func (w *WarfareNode) IssueCommander(id string) (CommanderCredential, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return CommanderCredential{}, errors.New("commander id required")
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return CommanderCredential{}, fmt.Errorf("generate commander: %w", err)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.commanders == nil {
		w.commanders = make(map[string]commanderEntry)
	}
	if _, exists := w.commanders[id]; exists {
		return CommanderCredential{}, fmt.Errorf("commander %s already authorised", id)
	}
	w.commanders[id] = commanderEntry{pub: pub}
	cred := CommanderCredential{
		ID:         id,
		PublicKey:  hex.EncodeToString(pub),
		PrivateKey: hex.EncodeToString(priv),
		IssuedAt:   time.Now().UTC(),
	}
	return cred, nil
}

// AuthorizeCommander registers an externally generated commander public key.
func (w *WarfareNode) AuthorizeCommander(id, hexKey string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("commander id required")
	}
	keyBytes, err := hex.DecodeString(strings.TrimSpace(hexKey))
	if err != nil {
		return fmt.Errorf("decode commander key: %w", err)
	}
	if len(keyBytes) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid commander key length: %d", len(keyBytes))
	}
	w.mu.Lock()
	if w.commanders == nil {
		w.commanders = make(map[string]commanderEntry)
	}
	w.commanders[id] = commanderEntry{pub: ed25519.PublicKey(keyBytes)}
	w.mu.Unlock()
	return nil
}

// RevokeCommander removes a commander from the authorised set.
func (w *WarfareNode) RevokeCommander(id string) {
	w.mu.Lock()
	delete(w.commanders, strings.TrimSpace(id))
	w.mu.Unlock()
}

// nextNonce increments the nonce for the provided commander returning the next
// value in sequence.
func (w *WarfareNode) nextNonce(commander string) uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	entry := w.commanders[commander]
	entry.lastNonce++
	w.commanders[commander] = entry
	return entry.lastNonce
}

func (r CommandRequest) canonicalPayload() []byte {
	ts := r.Timestamp.UTC().Format(time.RFC3339Nano)
	keys := make([]string, 0, len(r.Metadata))
	for k := range r.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(r.Commander)
	b.WriteByte('|')
	b.WriteString(r.Command)
	b.WriteByte('|')
	b.WriteString(ts)
	b.WriteByte('|')
	b.WriteString(fmt.Sprintf("%d", r.Nonce))
	for _, k := range keys {
		b.WriteByte('|')
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(r.Metadata[k])
	}
	return []byte(b.String())
}

// CanonicalPayload exposes the byte sequence used for signing a command
// request.  External clients such as the CLI and web console rely on this to
// deterministically reproduce the payload the warfare node verifies before
// executing a command.
func (r CommandRequest) CanonicalPayload() []byte {
	return r.canonicalPayload()
}

func (w *WarfareNode) snapshotConsensus() ConsensusSnapshot {
	if w.Node == nil || w.Node.Consensus == nil {
		return ConsensusSnapshot{}
	}
	weights := w.Node.Consensus.Weights
	return ConsensusSnapshot{PoW: weights.PoW, PoS: weights.PoS, PoH: weights.PoH}
}

// ExecuteSecureCommand validates the supplied request, verifies its signature,
// checks consensus availability and executes the VM hook.  The resulting record
// is stored in the audit log and broadcast to subscribers.
func (w *WarfareNode) ExecuteSecureCommand(ctx context.Context, req CommandRequest) (CommandRecord, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req.Command = strings.TrimSpace(req.Command)
	if req.Command == "" {
		return CommandRecord{}, errors.New("command required")
	}
	req.Commander = strings.TrimSpace(req.Commander)
	if req.Commander == "" {
		return CommandRecord{}, errors.New("commander required")
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now().UTC()
	}
	if delta := time.Since(req.Timestamp); delta > commandSkewTolerance || delta < -commandSkewTolerance {
		return CommandRecord{}, fmt.Errorf("timestamp outside tolerance: %s", req.Timestamp)
	}

	payload := req.canonicalPayload()
	w.mu.Lock()
	entry, ok := w.commanders[req.Commander]
	if !ok {
		w.mu.Unlock()
		return CommandRecord{}, fmt.Errorf("unknown commander: %s", req.Commander)
	}
	if req.Nonce == 0 {
		entry.lastNonce++
		req.Nonce = entry.lastNonce
	} else if req.Nonce <= entry.lastNonce {
		w.mu.Unlock()
		return CommandRecord{}, fmt.Errorf("nonce %d not greater than %d", req.Nonce, entry.lastNonce)
	} else {
		entry.lastNonce = req.Nonce
	}
	if len(req.Signature) == 0 || !ed25519.Verify(entry.pub, payload, req.Signature) {
		w.mu.Unlock()
		atomic.AddUint64(&w.metrics.failedCommands, 1)
		record := CommandRecord{
			Commander: req.Commander,
			Command:   req.Command,
			Timestamp: req.Timestamp,
			Nonce:     req.Nonce,
			Metadata:  cloneMetadata(req.Metadata),
			Accepted:  false,
			Error:     "signature verification failed",
			Consensus: w.snapshotConsensus(),
		}
		w.appendCommandRecord(record)
		w.recordEvent(WarfareEvent{
			Type:      WarfareEventCommandRejected,
			Timestamp: time.Now().UTC(),
			Commander: req.Commander,
			Payload: map[string]string{
				"command": req.Command,
				"error":   record.Error,
			},
		})
		return record, errors.New(record.Error)
	}
	w.commanders[req.Commander] = entry
	w.mu.Unlock()

	select {
	case <-ctx.Done():
		atomic.AddUint64(&w.metrics.failedCommands, 1)
		record := CommandRecord{
			Commander: req.Commander,
			Command:   req.Command,
			Timestamp: req.Timestamp,
			Nonce:     req.Nonce,
			Metadata:  cloneMetadata(req.Metadata),
			Accepted:  false,
			Error:     ctx.Err().Error(),
			Consensus: w.snapshotConsensus(),
		}
		w.appendCommandRecord(record)
		w.recordEvent(WarfareEvent{
			Type:      WarfareEventCommandRejected,
			Timestamp: time.Now().UTC(),
			Commander: req.Commander,
			Payload: map[string]string{
				"command": req.Command,
				"error":   record.Error,
			},
		})
		return record, ctx.Err()
	default:
	}

	start := time.Now()
	var vmErr error
	if w.Node != nil && w.Node.VM != nil {
		if caller, ok := interface{}(w.Node.VM).(vmCallable); ok {
			vmErr = caller.Call("SecureCommand")
		}
	}
	latency := time.Since(start)

	record := CommandRecord{
		Commander: req.Commander,
		Command:   req.Command,
		Timestamp: req.Timestamp,
		Nonce:     req.Nonce,
		Metadata:  cloneMetadata(req.Metadata),
		Consensus: w.snapshotConsensus(),
		Latency:   latency,
	}
	if vmErr != nil {
		record.Accepted = false
		record.Error = vmErr.Error()
		atomic.AddUint64(&w.metrics.failedCommands, 1)
		w.recordEvent(WarfareEvent{
			Type:      WarfareEventCommandRejected,
			Timestamp: time.Now().UTC(),
			Commander: req.Commander,
			Payload: map[string]string{
				"command": req.Command,
				"error":   record.Error,
			},
		})
	} else {
		record.Accepted = true
		atomic.AddUint64(&w.metrics.commands, 1)
		w.recordEvent(WarfareEvent{
			Type:      WarfareEventCommandAccepted,
			Timestamp: time.Now().UTC(),
			Commander: req.Commander,
			Payload: map[string]string{
				"command": req.Command,
				"latency": latency.String(),
			},
		})
	}
	w.appendCommandRecord(record)
	if vmErr != nil {
		return record, vmErr
	}
	return record, nil
}

// SecureCommand executes a privileged command on behalf of the root commander.
// The CLI falls back to this helper when no explicit credentials are supplied.
func (w *WarfareNode) SecureCommand(cmd string) error {
	nonce := w.nextNonce(commanderRootID)
	req := CommandRequest{
		Commander: commanderRootID,
		Command:   cmd,
		Timestamp: time.Now().UTC(),
		Nonce:     nonce,
	}
	req.Metadata = map[string]string{"origin": "internal"}
	req.Signature = ed25519.Sign(w.rootKey, req.canonicalPayload())
	_, err := w.ExecuteSecureCommand(context.Background(), req)
	return err
}

// RecordLogistics validates and stores a logistics update returning the
// canonical record.  TrackLogistics retains its historical signature and simply
// delegates to this helper.
func (w *WarfareNode) RecordLogistics(update LogisticsUpdate) (militarynodes.LogisticsRecord, error) {
	update.AssetID = strings.TrimSpace(update.AssetID)
	update.Location = strings.TrimSpace(update.Location)
	update.Status = strings.TrimSpace(update.Status)
	if update.AssetID == "" || update.Location == "" || update.Status == "" {
		atomic.AddUint64(&w.metrics.failedCommands, 1)
		w.recordEvent(WarfareEvent{
			Type:      WarfareEventLogisticsRejected,
			Timestamp: time.Now().UTC(),
			Payload: map[string]string{
				"asset": update.AssetID,
				"error": "asset, location and status required",
			},
		})
		return militarynodes.LogisticsRecord{}, errors.New("asset, location and status required")
	}
	if update.Timestamp.IsZero() {
		update.Timestamp = time.Now().UTC()
	}
	rec := militarynodes.LogisticsRecord{
		AssetID:   update.AssetID,
		Location:  update.Location,
		Status:    update.Status,
		Timestamp: update.Timestamp,
	}
	w.mu.Lock()
	w.logistics = append(w.logistics, rec)
	w.mu.Unlock()
	atomic.AddUint64(&w.metrics.logistics, 1)
	payload := map[string]string{
		"asset":    update.AssetID,
		"location": update.Location,
		"status":   update.Status,
	}
	if update.Reporter != "" {
		payload["reporter"] = update.Reporter
	}
	for k, v := range update.Metadata {
		payload[k] = v
	}
	w.recordEvent(WarfareEvent{
		Type:      WarfareEventLogisticsRecorded,
		Timestamp: update.Timestamp,
		Payload:   payload,
	})
	return rec, nil
}

// TrackLogistics maintains backwards compatibility with earlier stages by
// delegating to RecordLogistics and discarding any validation errors.
func (w *WarfareNode) TrackLogistics(assetID, location, status string) {
	_, _ = w.RecordLogistics(LogisticsUpdate{AssetID: assetID, Location: location, Status: status})
}

// BroadcastTactical shares tactical information while logging the broadcast for
// auditing and replay by user interfaces.
func (w *WarfareNode) BroadcastTactical(info string, metadata map[string]string) error {
	info = strings.TrimSpace(info)
	if info == "" {
		return errors.New("tactical info required")
	}
	payload := map[string]string{"message": info}
	for k, v := range metadata {
		payload[k] = v
	}
	atomic.AddUint64(&w.metrics.tactical, 1)
	w.recordEvent(WarfareEvent{
		Type:      WarfareEventTacticalBroadcast,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	})
	return nil
}

// ShareTactical provides the historical method signature expected by existing
// integrations.
func (w *WarfareNode) ShareTactical(info string) {
	_ = w.BroadcastTactical(info, nil)
}

// Logistics returns a copy of stored logistics records.
func (w *WarfareNode) Logistics() []militarynodes.LogisticsRecord {
	w.mu.RLock()
	defer w.mu.RUnlock()
	cp := make([]militarynodes.LogisticsRecord, len(w.logistics))
	copy(cp, w.logistics)
	return cp
}

// LogisticsByAsset filters stored logistics records for a specific asset ID.
// The returned slice is a copy and may be safely modified by the caller.
func (w *WarfareNode) LogisticsByAsset(assetID string) []militarynodes.LogisticsRecord {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var res []militarynodes.LogisticsRecord
	for _, r := range w.logistics {
		if r.AssetID == assetID {
			res = append(res, r)
		}
	}
	cp := make([]militarynodes.LogisticsRecord, len(res))
	copy(cp, res)
	return cp
}

// CommandLog returns a copy of the privileged command audit trail.
func (w *WarfareNode) CommandLog() []CommandRecord {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]CommandRecord, len(w.commandLog))
	copy(out, w.commandLog)
	return out
}

// Events returns a copy of the retained event log suitable for REST responses
// or CLI JSON output.
func (w *WarfareNode) Events() []WarfareEvent {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]WarfareEvent, len(w.events))
	copy(out, w.events)
	return out
}

// EventsSince returns events with a sequence number greater than the supplied
// value allowing incremental polling by thin clients.
func (w *WarfareNode) EventsSince(seq uint64) []WarfareEvent {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var out []WarfareEvent
	for _, ev := range w.events {
		if ev.Sequence > seq {
			out = append(out, ev)
		}
	}
	return out
}

// SubscribeEvents registers a buffered channel that receives future events. The
// returned cancel function must be invoked by callers to avoid leaks.
func (w *WarfareNode) SubscribeEvents(buffer int) (<-chan WarfareEvent, func()) {
	if buffer <= 0 {
		buffer = 16
	}
	ch := make(chan WarfareEvent, buffer)
	w.mu.Lock()
	if w.subscribers == nil {
		w.subscribers = make(map[uint64]chan WarfareEvent)
	}
	w.subscriber++
	id := w.subscriber
	w.subscribers[id] = ch
	backlog := append([]WarfareEvent(nil), w.events...)
	w.mu.Unlock()

	go func(events []WarfareEvent) {
		for _, ev := range events {
			ch <- ev
		}
	}(backlog)

	cancel := func() {
		w.mu.Lock()
		if ch, ok := w.subscribers[id]; ok {
			delete(w.subscribers, id)
			close(ch)
		}
		w.mu.Unlock()
	}
	return ch, cancel
}

// WarfareMetrics summarises operational statistics for dashboards and tests.
type WarfareMetrics struct {
	Commands       uint64
	FailedCommands uint64
	Logistics      uint64
	Tactical       uint64
}

// MetricsSnapshot returns aggregated counters for observability dashboards.
func (w *WarfareNode) MetricsSnapshot() WarfareMetrics {
	return WarfareMetrics{
		Commands:       atomic.LoadUint64(&w.metrics.commands),
		FailedCommands: atomic.LoadUint64(&w.metrics.failedCommands),
		Logistics:      atomic.LoadUint64(&w.metrics.logistics),
		Tactical:       atomic.LoadUint64(&w.metrics.tactical),
	}
}

func (w *WarfareNode) appendCommandRecord(record CommandRecord) {
	w.mu.Lock()
	w.commandLog = append(w.commandLog, record)
	w.mu.Unlock()
}

func (w *WarfareNode) recordEvent(ev WarfareEvent) {
	w.mu.Lock()
	w.eventSeq++
	ev.Sequence = w.eventSeq
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	w.events = append(w.events, ev)
	if w.evLimit > 0 && len(w.events) > w.evLimit {
		w.events = append([]WarfareEvent(nil), w.events[len(w.events)-w.evLimit:]...)
	}
	subs := make([]chan WarfareEvent, 0, len(w.subscribers))
	for _, ch := range w.subscribers {
		subs = append(subs, ch)
	}
	w.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

func cloneMetadata(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// ensure interface compliance
var _ militarynodes.WarfareNode = (*WarfareNode)(nil)
