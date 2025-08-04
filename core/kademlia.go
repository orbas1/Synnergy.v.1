package core

import (
	"encoding/hex"
	"math/big"
	"sort"
	"sync"
)

// Kademlia provides a minimal key/value store with XOR distance lookups
// inspired by the Kademlia DHT. It is intended for lightweight peer discovery
// and metadata storage.
type Kademlia struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewKademlia creates an empty Kademlia table.
func NewKademlia() *Kademlia {
	return &Kademlia{store: make(map[string][]byte)}
}

// Store saves a key/value pair in the DHT.
func (k *Kademlia) Store(key string, value []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.store[key] = append([]byte(nil), value...)
}

// FindValue retrieves a value for a given key.
func (k *Kademlia) FindValue(key string) ([]byte, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	v, ok := k.store[key]
	if !ok {
		return nil, false
	}
	return append([]byte(nil), v...), true
}

// Distance returns the XOR distance between two hex encoded identifiers.
func Distance(a, b string) *big.Int {
	ab, _ := hex.DecodeString(a)
	bb, _ := hex.DecodeString(b)
	// Pad shorter slice
	if len(ab) > len(bb) {
		bb = append(make([]byte, len(ab)-len(bb)), bb...)
	} else if len(bb) > len(ab) {
		ab = append(make([]byte, len(bb)-len(ab)), ab...)
	}
	dist := new(big.Int)
	for i := 0; i < len(ab); i++ {
		dist.Lsh(dist, 8)
		dist.Or(dist, big.NewInt(int64(ab[i]^bb[i])))
	}
	return dist
}

// Closest returns up to n keys in the store sorted by XOR distance to the
// target identifier. It is a helper for peer lookups in the simulated DHT.
func (k *Kademlia) Closest(target string, n int) []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	type kv struct {
		key  string
		dist *big.Int
	}
	arr := make([]kv, 0, len(k.store))
	for key := range k.store {
		arr = append(arr, kv{key: key, dist: Distance(target, key)})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].dist.Cmp(arr[j].dist) < 0 })
	if n > len(arr) {
		n = len(arr)
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = arr[i].key
	}
	return out
}
