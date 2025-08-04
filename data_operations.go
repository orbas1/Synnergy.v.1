package synnergy

import (
	"sync"
	"time"
)

// DataFeed holds structured data referenced on chain. It allows concurrent
// updates and retrieval of key/value pairs representing external datasets.
type DataFeed struct {
	ID string

	mu      sync.RWMutex
	data    map[string]string
	updated time.Time
}

// NewDataFeed creates a new DataFeed with the provided identifier.
func NewDataFeed(id string) *DataFeed {
	return &DataFeed{ID: id, data: make(map[string]string)}
}

// Update sets a key/value pair within the feed and records the update time.
func (f *DataFeed) Update(key, value string) {
	f.mu.Lock()
	f.data[key] = value
	f.updated = time.Now().UTC()
	f.mu.Unlock()
}

// Get retrieves a value by key from the feed.
func (f *DataFeed) Get(key string) (string, bool) {
	f.mu.RLock()
	val, ok := f.data[key]
	f.mu.RUnlock()
	return val, ok
}

// Delete removes a key from the feed and records the update time if the key existed.
func (f *DataFeed) Delete(key string) {
	f.mu.Lock()
	if _, ok := f.data[key]; ok {
		delete(f.data, key)
		f.updated = time.Now().UTC()
	}
	f.mu.Unlock()
}

// Keys returns a snapshot of all keys stored in the feed.
func (f *DataFeed) Keys() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]string, 0, len(f.data))
	for k := range f.data {
		keys = append(keys, k)
	}
	return keys
}

// Snapshot returns a copy of the feed's current data map.
func (f *DataFeed) Snapshot() map[string]string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make(map[string]string, len(f.data))
	for k, v := range f.data {
		out[k] = v
	}
	return out
}

// LastUpdated returns the timestamp of the most recent update.
func (f *DataFeed) LastUpdated() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.updated
}
