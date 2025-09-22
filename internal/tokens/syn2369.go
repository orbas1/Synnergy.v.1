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
	Archived    bool
}

// ItemRegistry manages virtual items.
type ItemRegistry struct {
	mu      sync.RWMutex
	items   map[string]*VirtualItem
	counter uint64
	events  map[string][]ItemEvent
}

// NewItemRegistry creates an empty registry.
func NewItemRegistry() *ItemRegistry {
	return &ItemRegistry{items: make(map[string]*VirtualItem), events: make(map[string][]ItemEvent)}
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
	r.appendEvent(id, ItemEvent{Type: "create", Actor: owner, Timestamp: item.CreatedAt, Metadata: attrs})
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
	r.appendEvent(itemID, ItemEvent{Type: "transfer", Actor: newOwner, Timestamp: time.Now()})
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
	r.appendEvent(itemID, ItemEvent{Type: "update", Actor: it.Owner, Timestamp: time.Now(), Metadata: attrs})
	return nil
}

// ArchiveItem marks an item as inactive without deleting its history.
func (r *ItemRegistry) ArchiveItem(itemID, actor string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	it, ok := r.items[itemID]
	if !ok {
		return errors.New("item not found")
	}
	it.Archived = true
	r.appendEvent(itemID, ItemEvent{Type: "archive", Actor: actor, Timestamp: time.Now()})
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

// ItemEvent captures lifecycle transitions for compliance and analytics.
type ItemEvent struct {
	Type      string
	Actor     string
	Metadata  map[string]string
	Timestamp time.Time
}

func (r *ItemRegistry) appendEvent(itemID string, evt ItemEvent) {
	if evt.Metadata != nil {
		cp := make(map[string]string, len(evt.Metadata))
		for k, v := range evt.Metadata {
			cp[k] = v
		}
		evt.Metadata = cp
	}
	r.events[itemID] = append(r.events[itemID], evt)
}

// Events returns the chronological event log for an item.
func (r *ItemRegistry) Events(itemID string, limit int) ([]ItemEvent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	events, ok := r.events[itemID]
	if !ok {
		return nil, false
	}
	if limit <= 0 || limit >= len(events) {
		out := make([]ItemEvent, len(events))
		copy(out, events)
		return out, true
	}
	out := make([]ItemEvent, limit)
	copy(out, events[len(events)-limit:])
	return out, true
}
