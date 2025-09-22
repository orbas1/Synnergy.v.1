package synnergy

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

// IndexEventType represents the kind of change emitted by the indexing node.
type IndexEventType string

const (
	// IndexEventIndexed indicates a key was added or updated.
	IndexEventIndexed IndexEventType = "indexed"
	// IndexEventRemoved indicates a key was removed from the index.
	IndexEventRemoved IndexEventType = "removed"
)

// IndexMetadata captures metadata about a stored value.
type IndexMetadata struct {
	Version   uint64
	UpdatedAt time.Time
	ExpiresAt time.Time
	Tags      map[string]string
}

// IndexResult represents the data returned from a query.
type IndexResult struct {
	Value    []byte
	Metadata IndexMetadata
}

// IndexEvent represents a change notification.
type IndexEvent struct {
	Type IndexEventType
	Key  string
	Data IndexMetadata
}

// IndexEventHandler is notified when the index changes.
type IndexEventHandler interface {
	HandleIndexEvent(ctx context.Context, event IndexEvent) error
}

// IndexEventHandlerFunc adapts a function into an IndexEventHandler.
type IndexEventHandlerFunc func(context.Context, IndexEvent) error

// HandleIndexEvent implements IndexEventHandler.
func (f IndexEventHandlerFunc) HandleIndexEvent(ctx context.Context, event IndexEvent) error {
	return f(ctx, event)
}

// IndexingNodeOption configures an IndexingNode instance.
type IndexingNodeOption func(*IndexingNode)

// WithIndexWatcher registers an event handler for index changes.
func WithIndexWatcher(handler IndexEventHandler) IndexingNodeOption {
	return func(n *IndexingNode) {
		if handler != nil {
			n.watchers = append(n.watchers, handler)
		}
	}
}

// WithIndexTTL sets a default TTL for entries added to the index.
func WithIndexTTL(ttl time.Duration) IndexingNodeOption {
	return func(n *IndexingNode) {
		n.defaultTTL = ttl
	}
}

// WithIndexMaxEntries limits how many entries may reside in the index simultaneously.
func WithIndexMaxEntries(max int) IndexingNodeOption {
	return func(n *IndexingNode) {
		if max > 0 {
			n.maxEntries = max
		}
	}
}

// WithIndexClock allows tests to override the clock source.
func WithIndexClock(now func() time.Time) IndexingNodeOption {
	return func(n *IndexingNode) {
		if now != nil {
			n.now = now
		}
	}
}

// IndexOption customises an indexing operation.
type IndexOption func(*indexOptions)

type indexOptions struct {
	ttl       time.Duration
	tags      map[string]string
	mergeTags bool
}

// WithEntryTTL sets the TTL for a specific entry.
func WithEntryTTL(ttl time.Duration) IndexOption {
	return func(o *indexOptions) {
		o.ttl = ttl
	}
}

// WithEntryTags attaches metadata tags to the entry.
func WithEntryTags(tags map[string]string) IndexOption {
	return func(o *indexOptions) {
		o.tags = tags
	}
}

// MergeEntryTags merges tags with existing metadata instead of replacing it.
func MergeEntryTags() IndexOption {
	return func(o *indexOptions) {
		o.mergeTags = true
	}
}

var (
	// ErrIndexKeyRequired is returned when an empty key is supplied.
	ErrIndexKeyRequired = errors.New("index: key required")
)

// IndexingNode provides fast query capabilities by indexing ledger data into an in-memory key/value store.
type IndexingNode struct {
	mu         sync.RWMutex
	entries    map[string]indexEntry
	watchers   []IndexEventHandler
	defaultTTL time.Duration
	maxEntries int
	now        func() time.Time
}

type indexEntry struct {
	value    []byte
	metadata IndexMetadata
}

// NewIndexingNode constructs an IndexingNode with optional configuration.
func NewIndexingNode(opts ...IndexingNodeOption) *IndexingNode {
	node := &IndexingNode{
		entries: make(map[string]indexEntry),
		now:     time.Now,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(node)
		}
	}
	return node
}

// Index inserts or updates the value for a given key.
func (n *IndexingNode) Index(key string, value []byte) {
	_, _ = n.IndexWithOptions(context.Background(), key, value)
}

// IndexWithOptions indexes a value applying the supplied options.
func (n *IndexingNode) IndexWithOptions(ctx context.Context, key string, value []byte, opts ...IndexOption) (IndexMetadata, error) {
	if key == "" {
		return IndexMetadata{}, ErrIndexKeyRequired
	}
	options := indexOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	now := n.now().UTC()
	ttl := options.ttl
	if ttl <= 0 {
		ttl = n.defaultTTL
	}
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = now.Add(ttl)
	}
	entry := indexEntry{
		value:    append([]byte(nil), value...),
		metadata: IndexMetadata{UpdatedAt: now, ExpiresAt: expiresAt, Tags: cloneMetadata(options.tags)},
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	n.pruneLocked(now)
	existing, ok := n.entries[key]
	if ok {
		entry.metadata.Version = existing.metadata.Version + 1
		if options.mergeTags && len(existing.metadata.Tags) > 0 {
			if entry.metadata.Tags == nil {
				entry.metadata.Tags = make(map[string]string, len(existing.metadata.Tags))
			}
			for k, v := range existing.metadata.Tags {
				if _, exists := entry.metadata.Tags[k]; !exists {
					entry.metadata.Tags[k] = v
				}
			}
		}
	} else {
		entry.metadata.Version = 1
	}
	if entry.metadata.Tags == nil && len(existing.metadata.Tags) > 0 && options.mergeTags {
		entry.metadata.Tags = cloneMetadata(existing.metadata.Tags)
	}
	if n.maxEntries > 0 && !ok && len(n.entries) >= n.maxEntries {
		n.evictOldestLocked()
	}
	n.entries[key] = entry
	metaCopy := copyIndexMetadata(entry.metadata)
	n.fireEvent(ctx, IndexEvent{Type: IndexEventIndexed, Key: key, Data: metaCopy})
	return metaCopy, nil
}

// Query retrieves a copy of the value associated with the key.
func (n *IndexingNode) Query(key string) ([]byte, bool) {
	res, ok := n.QueryWithMetadata(key)
	if !ok {
		return nil, false
	}
	return res.Value, true
}

// QueryWithMetadata retrieves both the value and metadata for a key.
func (n *IndexingNode) QueryWithMetadata(key string) (IndexResult, bool) {
	now := n.now().UTC()
	n.mu.Lock()
	defer n.mu.Unlock()
	n.pruneLocked(now)
	entry, ok := n.entries[key]
	if !ok {
		return IndexResult{}, false
	}
	if !entry.metadata.ExpiresAt.IsZero() && now.After(entry.metadata.ExpiresAt) {
		delete(n.entries, key)
		n.fireEvent(context.Background(), IndexEvent{Type: IndexEventRemoved, Key: key, Data: copyIndexMetadata(entry.metadata)})
		return IndexResult{}, false
	}
	value := append([]byte(nil), entry.value...)
	return IndexResult{Value: value, Metadata: copyIndexMetadata(entry.metadata)}, true
}

// Remove deletes the key from the index if present.
func (n *IndexingNode) Remove(ctx context.Context, key string) {
	n.mu.Lock()
	entry, ok := n.entries[key]
	if ok {
		delete(n.entries, key)
	}
	n.mu.Unlock()
	if ok {
		n.fireEvent(ctx, IndexEvent{Type: IndexEventRemoved, Key: key, Data: copyIndexMetadata(entry.metadata)})
	}
}

// Keys returns a snapshot of all indexed keys.
func (n *IndexingNode) Keys() []string {
	n.mu.Lock()
	n.pruneLocked(n.now().UTC())
	keys := make([]string, 0, len(n.entries))
	for k := range n.entries {
		keys = append(keys, k)
	}
	n.mu.Unlock()
	sort.Strings(keys)
	return keys
}

// Count returns the number of entries currently indexed.
func (n *IndexingNode) Count() int {
	n.mu.Lock()
	n.pruneLocked(n.now().UTC())
	count := len(n.entries)
	n.mu.Unlock()
	return count
}

// Snapshot returns a copy of all values stored in the index.
func (n *IndexingNode) Snapshot() map[string]IndexResult {
	n.mu.Lock()
	n.pruneLocked(n.now().UTC())
	out := make(map[string]IndexResult, len(n.entries))
	for k, entry := range n.entries {
		out[k] = IndexResult{Value: append([]byte(nil), entry.value...), Metadata: copyIndexMetadata(entry.metadata)}
	}
	n.mu.Unlock()
	return out
}

// IndexBatch indexes multiple key/value pairs atomically.
func (n *IndexingNode) IndexBatch(ctx context.Context, values map[string][]byte, opts ...IndexOption) error {
	now := n.now().UTC()
	options := indexOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	ttl := options.ttl
	if ttl <= 0 {
		ttl = n.defaultTTL
	}
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = now.Add(ttl)
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	n.pruneLocked(now)
	for key, val := range values {
		if key == "" {
			continue
		}
		entry := indexEntry{
			value: append([]byte(nil), val...),
			metadata: IndexMetadata{
				Version:   1,
				UpdatedAt: now,
				ExpiresAt: expiresAt,
				Tags:      cloneMetadata(options.tags),
			},
		}
		if existing, ok := n.entries[key]; ok {
			entry.metadata.Version = existing.metadata.Version + 1
		} else if n.maxEntries > 0 && len(n.entries) >= n.maxEntries {
			n.evictOldestLocked()
		}
		n.entries[key] = entry
		n.fireEvent(ctx, IndexEvent{Type: IndexEventIndexed, Key: key, Data: copyIndexMetadata(entry.metadata)})
	}
	return nil
}

func (n *IndexingNode) pruneLocked(now time.Time) {
	for key, entry := range n.entries {
		if entry.metadata.ExpiresAt.IsZero() {
			continue
		}
		if now.After(entry.metadata.ExpiresAt) {
			delete(n.entries, key)
			n.fireEvent(context.Background(), IndexEvent{Type: IndexEventRemoved, Key: key, Data: copyIndexMetadata(entry.metadata)})
		}
	}
}

func (n *IndexingNode) evictOldestLocked() {
	var oldestKey string
	var oldest time.Time
	first := true
	for key, entry := range n.entries {
		if first || entry.metadata.UpdatedAt.Before(oldest) {
			oldest = entry.metadata.UpdatedAt
			oldestKey = key
			first = false
		}
	}
	if oldestKey != "" {
		entry := n.entries[oldestKey]
		delete(n.entries, oldestKey)
		n.fireEvent(context.Background(), IndexEvent{Type: IndexEventRemoved, Key: oldestKey, Data: copyIndexMetadata(entry.metadata)})
	}
}

func (n *IndexingNode) fireEvent(ctx context.Context, event IndexEvent) {
	n.mu.RLock()
	watchers := append([]IndexEventHandler(nil), n.watchers...)
	n.mu.RUnlock()
	for _, watcher := range watchers {
		if watcher == nil {
			continue
		}
		_ = watcher.HandleIndexEvent(ctx, event)
	}
}

func copyIndexMetadata(in IndexMetadata) IndexMetadata {
	return IndexMetadata{
		Version:   in.Version,
		UpdatedAt: in.UpdatedAt,
		ExpiresAt: in.ExpiresAt,
		Tags:      cloneMetadata(in.Tags),
	}
}
