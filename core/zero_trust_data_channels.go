package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultChannelRetention = 256

// DataChannel represents a secure channel with an encryption key.
type DataChannel struct {
	Key      []byte
	PrivKey  ed25519.PrivateKey
	PubKey   ed25519.PublicKey
	Messages []SignedMessage
	Open     bool

	Owner        string
	Participants map[string]ed25519.PublicKey
	Retention    int
	Metadata     map[string]string
	Sequence     uint64
}

// SignedMessage bundles an encrypted payload with its signature.
type SignedMessage struct {
	Cipher    []byte
	Signature []byte
	Sender    string
	Timestamp time.Time
}

// ChannelInfo exposes metadata about a channel for CLI and web dashboards.
type ChannelInfo struct {
	ID           string
	Owner        string
	Participants []string
	MessageCount int
	Open         bool
	Metadata     map[string]string
}

// ChannelEventType labels emitted zero-trust events.
type ChannelEventType string

const (
	ChannelEventOpened      ChannelEventType = "channel.opened"
	ChannelEventMessage     ChannelEventType = "channel.message"
	ChannelEventClosed      ChannelEventType = "channel.closed"
	ChannelEventKeyRotated  ChannelEventType = "channel.key_rotated"
	ChannelEventParticipant ChannelEventType = "channel.participant"
)

// ChannelEvent captures mutations broadcast to subscribers.
type ChannelEvent struct {
	Sequence  uint64
	Type      ChannelEventType
	Timestamp time.Time
	ChannelID string
	Payload   map[string]string
}

type channelConfig struct {
	owner        string
	participants map[string]ed25519.PublicKey
	retention    int
	metadata     map[string]string
}

// ChannelOption customises channel configuration during creation.
type ChannelOption func(*channelConfig)

// WithOwner assigns an owner to the channel.
func WithOwner(owner string) ChannelOption {
	return func(cfg *channelConfig) { cfg.owner = strings.TrimSpace(owner) }
}

// WithParticipant registers a participant public key that can submit messages.
func WithParticipant(id string, pub ed25519.PublicKey) ChannelOption {
	clone := append(ed25519.PublicKey(nil), pub...)
	return func(cfg *channelConfig) {
		if cfg.participants == nil {
			cfg.participants = make(map[string]ed25519.PublicKey)
		}
		cfg.participants[strings.TrimSpace(id)] = clone
	}
}

// WithRetention overrides the message retention count.
func WithRetention(limit int) ChannelOption {
	return func(cfg *channelConfig) { cfg.retention = limit }
}

// WithChannelMetadata attaches arbitrary metadata for auditing.
func WithChannelMetadata(meta map[string]string) ChannelOption {
	return func(cfg *channelConfig) { cfg.metadata = cloneStringMap(meta) }
}

// ZeroTrustEngine manages encrypted data channels backed by ledger escrows.
type ZeroTrustEngine struct {
	mu        sync.RWMutex
	channels  map[string]*DataChannel
	events    []ChannelEvent
	eventSeq  uint64
	watchers  map[uint64]chan ChannelEvent
	watcher   uint64
	retention int
}

// NewZeroTrustEngine creates a new ZeroTrustEngine instance.
func NewZeroTrustEngine() *ZeroTrustEngine {
	return &ZeroTrustEngine{
		channels:  make(map[string]*DataChannel),
		retention: defaultChannelRetention,
	}
}

// SetDefaultRetention adjusts the default per-channel message retention.
func (e *ZeroTrustEngine) SetDefaultRetention(limit int) {
	if limit <= 0 {
		limit = defaultChannelRetention
	}
	e.mu.Lock()
	e.retention = limit
	e.mu.Unlock()
}

// OpenChannel initialises a new channel with the provided key and options.
func (e *ZeroTrustEngine) OpenChannel(id string, key []byte, opts ...ChannelOption) error {
	if len(key) != 32 {
		return fmt.Errorf("invalid key length: %d", len(key))
	}
	cfg := channelConfig{retention: e.retention}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	e.mu.Lock()
	if _, exists := e.channels[id]; exists {
		e.mu.Unlock()
		return errors.New("channel already exists")
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		e.mu.Unlock()
		return err
	}
	if cfg.participants == nil {
		cfg.participants = make(map[string]ed25519.PublicKey)
	}
	if cfg.owner != "" {
		cfg.participants[cfg.owner] = pub
	}
	ch := &DataChannel{
		Key:          append([]byte(nil), key...),
		PrivKey:      priv,
		PubKey:       pub,
		Open:         true,
		Owner:        cfg.owner,
		Participants: cfg.participants,
		Retention:    cfg.retention,
		Metadata:     cfg.metadata,
	}
	e.channels[id] = ch
	e.mu.Unlock()
	e.recordEvent(ChannelEvent{Type: ChannelEventOpened, ChannelID: id, Timestamp: time.Now().UTC(), Payload: map[string]string{
		"owner": cfg.owner,
	}})
	return nil
}

// Send encrypts and stores a payload on the channel using the owner identity.
func (e *ZeroTrustEngine) Send(id string, payload []byte) ([]byte, error) {
	return e.SendAs(id, "", payload)
}

// SendAs encrypts and stores a payload on the channel for a specific sender.
func (e *ZeroTrustEngine) SendAs(id, sender string, payload []byte) ([]byte, error) {
	e.mu.RLock()
	ch, ok := e.channels[id]
	e.mu.RUnlock()
	if !ok || !ch.Open {
		return nil, errors.New("channel not open")
	}
	if sender == "" {
		sender = ch.Owner
	}
	if sender == "" {
		return nil, errors.New("sender required")
	}
	if _, ok := ch.Participants[sender]; !ok {
		return nil, fmt.Errorf("sender %s not authorised", sender)
	}
	cipherText, err := Encrypt(ch.Key, payload)
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(ch.PrivKey, cipherText)
	msg := SignedMessage{Cipher: cipherText, Signature: sig, Sender: sender, Timestamp: time.Now().UTC()}
	e.mu.Lock()
	ch.Messages = append(ch.Messages, msg)
	if ch.Retention > 0 && len(ch.Messages) > ch.Retention {
		ch.Messages = append([]SignedMessage(nil), ch.Messages[len(ch.Messages)-ch.Retention:]...)
	}
	ch.Sequence++
	seq := ch.Sequence
	e.mu.Unlock()
	e.recordEvent(ChannelEvent{Type: ChannelEventMessage, ChannelID: id, Timestamp: msg.Timestamp, Payload: map[string]string{
		"sender":   sender,
		"sequence": fmt.Sprintf("%d", seq),
	}})
	return cipherText, nil
}

// Messages returns encrypted messages for a channel.
func (e *ZeroTrustEngine) Messages(id string) []SignedMessage {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ch, ok := e.channels[id]
	if !ok {
		return nil
	}
	out := make([]SignedMessage, len(ch.Messages))
	copy(out, ch.Messages)
	return out
}

// Receive verifies and decrypts a stored message by index.
func (e *ZeroTrustEngine) Receive(id string, index int) ([]byte, error) {
	e.mu.RLock()
	ch, ok := e.channels[id]
	e.mu.RUnlock()
	if !ok {
		return nil, errors.New("channel not found")
	}
	if index < 0 || index >= len(ch.Messages) {
		return nil, errors.New("message index out of range")
	}
	msg := ch.Messages[index]
	if !ed25519.Verify(ch.PubKey, msg.Cipher, msg.Signature) {
		return nil, errors.New("signature verification failed")
	}
	pt, err := Decrypt(ch.Key, msg.Cipher)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

// AuthorizePeer adds a new participant to the channel.
func (e *ZeroTrustEngine) AuthorizePeer(id, participant, pubHex string) error {
	key, err := hex.DecodeString(strings.TrimSpace(pubHex))
	if err != nil {
		return fmt.Errorf("decode key: %w", err)
	}
	if len(key) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid public key size: %d", len(key))
	}
	e.mu.Lock()
	ch, ok := e.channels[id]
	if !ok {
		e.mu.Unlock()
		return errors.New("channel not found")
	}
	if ch.Participants == nil {
		ch.Participants = make(map[string]ed25519.PublicKey)
	}
	ch.Participants[participant] = ed25519.PublicKey(key)
	e.mu.Unlock()
	e.recordEvent(ChannelEvent{Type: ChannelEventParticipant, ChannelID: id, Timestamp: time.Now().UTC(), Payload: map[string]string{
		"participant": participant,
		"action":      "authorized",
	}})
	return nil
}

// RevokePeer removes a participant from the channel.
func (e *ZeroTrustEngine) RevokePeer(id, participant string) {
	e.mu.Lock()
	_, ok := e.channels[id]
	if ok {
		delete(e.channels[id].Participants, participant)
	}
	e.mu.Unlock()
	if ok {
		e.recordEvent(ChannelEvent{Type: ChannelEventParticipant, ChannelID: id, Timestamp: time.Now().UTC(), Payload: map[string]string{
			"participant": participant,
			"action":      "revoked",
		}})
	}
}

// RotateKey replaces the symmetric key for a channel.
func (e *ZeroTrustEngine) RotateKey(id string, newKey []byte) error {
	if len(newKey) != 32 {
		return fmt.Errorf("invalid key length: %d", len(newKey))
	}
	e.mu.Lock()
	ch, ok := e.channels[id]
	if !ok {
		e.mu.Unlock()
		return errors.New("channel not found")
	}
	ch.Key = append(ch.Key[:0], newKey...)
	e.mu.Unlock()
	e.recordEvent(ChannelEvent{Type: ChannelEventKeyRotated, ChannelID: id, Timestamp: time.Now().UTC()})
	return nil
}

// CloseChannel closes the channel and prevents further messages.
func (e *ZeroTrustEngine) CloseChannel(id string) error {
	e.mu.Lock()
	ch, ok := e.channels[id]
	if !ok {
		e.mu.Unlock()
		return errors.New("channel not found")
	}
	ch.Open = false
	e.mu.Unlock()
	e.recordEvent(ChannelEvent{Type: ChannelEventClosed, ChannelID: id, Timestamp: time.Now().UTC()})
	return nil
}

// ChannelInfo returns high-level details about a channel.
func (e *ZeroTrustEngine) ChannelInfo(id string) (ChannelInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ch, ok := e.channels[id]
	if !ok {
		return ChannelInfo{}, errors.New("channel not found")
	}
	participants := make([]string, 0, len(ch.Participants))
	for p := range ch.Participants {
		participants = append(participants, p)
	}
	sort.Strings(participants)
	return ChannelInfo{
		ID:           id,
		Owner:        ch.Owner,
		Participants: participants,
		MessageCount: len(ch.Messages),
		Open:         ch.Open,
		Metadata:     cloneStringMap(ch.Metadata),
	}, nil
}

// Events returns a copy of the event log.
func (e *ZeroTrustEngine) Events() []ChannelEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]ChannelEvent, len(e.events))
	copy(out, e.events)
	return out
}

// EventsSince filters events newer than the supplied sequence.
func (e *ZeroTrustEngine) EventsSince(seq uint64) []ChannelEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []ChannelEvent
	for _, ev := range e.events {
		if ev.Sequence > seq {
			out = append(out, ev)
		}
	}
	return out
}

// SubscribeEvents registers a channel that receives future events.
func (e *ZeroTrustEngine) SubscribeEvents(buffer int) (<-chan ChannelEvent, func()) {
	if buffer <= 0 {
		buffer = 16
	}
	ch := make(chan ChannelEvent, buffer)
	e.mu.Lock()
	if e.watchers == nil {
		e.watchers = make(map[uint64]chan ChannelEvent)
	}
	e.watcher++
	id := e.watcher
	e.watchers[id] = ch
	backlog := append([]ChannelEvent(nil), e.events...)
	e.mu.Unlock()

	go func(events []ChannelEvent) {
		for _, ev := range events {
			ch <- ev
		}
	}(backlog)

	cancel := func() {
		e.mu.Lock()
		if ch, ok := e.watchers[id]; ok {
			delete(e.watchers, id)
			close(ch)
		}
		e.mu.Unlock()
	}
	return ch, cancel
}

func (e *ZeroTrustEngine) recordEvent(ev ChannelEvent) {
	e.mu.Lock()
	e.eventSeq++
	ev.Sequence = e.eventSeq
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	e.events = append(e.events, ev)
	watchers := make([]chan ChannelEvent, 0, len(e.watchers))
	for _, ch := range e.watchers {
		watchers = append(watchers, ch)
	}
	e.mu.Unlock()

	for _, ch := range watchers {
		select {
		case ch <- ev:
		default:
		}
	}
}

func cloneStringMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
