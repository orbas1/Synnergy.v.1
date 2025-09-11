package core

import (
	"errors"
	"sync"
)

var (
	// ErrTicketSupplyExhausted is returned when issuing beyond available supply.
	ErrTicketSupplyExhausted = errors.New("ticket supply exhausted")
	// ErrTicketNotOwned is returned when transfer requester is not current owner.
	ErrTicketNotOwned = errors.New("ticket not owned by sender")
)

// EventMetadata holds event information and issued tickets for SYN1700 tokens.
type EventMetadata struct {
	Name        string
	Description string
	Location    string
	Start       int64
	End         int64
	Supply      uint64

	mu           sync.RWMutex
	nextTicketID uint64
	Tickets      map[uint64]*Ticket
}

// Ticket represents an issued event ticket.
type Ticket struct {
	ID    uint64
	Owner string
	Class string
	Type  string
	Price uint64
}

// NewEvent initialises a new event metadata record.
func NewEvent(name, desc, location string, start, end int64, supply uint64) *EventMetadata {
	return &EventMetadata{
		Name:        name,
		Description: desc,
		Location:    location,
		Start:       start,
		End:         end,
		Supply:      supply,
		Tickets:     make(map[uint64]*Ticket),
	}
}

// IssueTicket issues a ticket to an owner if supply allows and returns its ID.
func (e *EventMetadata) IssueTicket(owner, class, ticketType string, price uint64) (uint64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if uint64(len(e.Tickets)) >= e.Supply {
		return 0, ErrTicketSupplyExhausted
	}
	e.nextTicketID++
	id := e.nextTicketID
	e.Tickets[id] = &Ticket{ID: id, Owner: owner, Class: class, Type: ticketType, Price: price}
	return id, nil
}

// TransferTicket transfers ownership of a ticket.
func (e *EventMetadata) TransferTicket(id uint64, from, to string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	t, ok := e.Tickets[id]
	if !ok || t.Owner != from {
		return ErrTicketNotOwned
	}
	t.Owner = to
	return nil
}

// VerifyTicket checks if a holder owns the ticket.
func (e *EventMetadata) VerifyTicket(id uint64, holder string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	t, ok := e.Tickets[id]
	return ok && t.Owner == holder
}
