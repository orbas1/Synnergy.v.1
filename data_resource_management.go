package synnergy

import "sync"

// DataResourceManager provides simple key/value storage with byte usage
// tracking. It is intended for managing content resources referenced by
// other modules.
type DataResourceManager struct {
	mu    sync.RWMutex
	store map[string][]byte
	size  int64
}

// NewDataResourceManager constructs an empty manager instance.
func NewDataResourceManager() *DataResourceManager {
	return &DataResourceManager{store: make(map[string][]byte)}
}

// Put stores data under the given key, replacing any existing entry.
func (m *DataResourceManager) Put(key string, data []byte) {
	m.mu.Lock()
	if old, ok := m.store[key]; ok {
		m.size -= int64(len(old))
	}
	dup := append([]byte(nil), data...)
	m.store[key] = dup
	m.size += int64(len(dup))
	m.mu.Unlock()
}

// Get retrieves a copy of the data for the specified key.
func (m *DataResourceManager) Get(key string) ([]byte, bool) {
	m.mu.RLock()
	val, ok := m.store[key]
	m.mu.RUnlock()
	if !ok {
		return nil, false
	}
	out := make([]byte, len(val))
	copy(out, val)
	return out, true
}

// Delete removes the key and associated data from the manager.
func (m *DataResourceManager) Delete(key string) {
	m.mu.Lock()
	if v, ok := m.store[key]; ok {
		m.size -= int64(len(v))
		delete(m.store, key)
	}
	m.mu.Unlock()
}

// Keys returns a slice of all keys currently stored.
func (m *DataResourceManager) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.store))
	for k := range m.store {
		keys = append(keys, k)
	}
	return keys
}

// Usage returns the total number of bytes currently stored.
func (m *DataResourceManager) Usage() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.size
}
