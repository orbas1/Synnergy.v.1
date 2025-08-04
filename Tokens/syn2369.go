package tokens

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// VirtualItem represents a virtual world asset for SYN2369 tokens.
type VirtualItem struct {
	ItemID      string
	Owner       string
	Name        string
	Description string
	Attributes  map[string]string
	CreatedAt   time.Time
}

// ItemRegistry manages virtual items.
type ItemRegistry struct {
	mu      sync.RWMutex
	items   map[string]*VirtualItem
	counter uint64
}

// NewItemRegistry creates an empty registry.
func NewItemRegistry() *ItemRegistry {
	return &ItemRegistry{items: make(map[string]*VirtualItem)}
}

// CreateItem registers a new virtual item.
func (r *ItemRegistry) CreateItem(owner, name, desc string, attrs map[string]string) *VirtualItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	id := fmt.Sprintf("VI-%d", r.counter)
	item := &VirtualItem{ItemID: id, Owner: owner, Name: name, Description: desc, Attributes: make(map[string]string), CreatedAt: time.Now()}
	for k, v := range attrs {
		item.Attributes[k] = v
	}
	r.items[id] = item
	return item
}

// TransferItem transfers ownership of an item.
func (r *ItemRegistry) TransferItem(itemID, newOwner string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok {
		return errors.New("item not found")
	}
	it.Owner = newOwner
	return nil
}

// UpdateAttributes updates custom attributes for an item.
func (r *ItemRegistry) UpdateAttributes(itemID string, attrs map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok {
		return errors.New("item not found")
	}
	for k, v := range attrs {
		it.Attributes[k] = v
	}
	return nil
}

// GetItem retrieves an item by ID.
func (r *ItemRegistry) GetItem(itemID string) (*VirtualItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.items[itemID]
	if !ok {
		return nil, false
	}
	cp := *it
	cp.Attributes = make(map[string]string)
	for k, v := range it.Attributes {
		cp.Attributes[k] = v
	}
	return &cp, true
}

// ListItems returns all virtual items.
func (r *ItemRegistry) ListItems() []*VirtualItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*VirtualItem, 0, len(r.items))
	for _, it := range r.items {
		cp := *it
		cp.Attributes = make(map[string]string)
		for k, v := range it.Attributes {
			cp.Attributes[k] = v
		}
		res = append(res, &cp)
	}
	return res
}
