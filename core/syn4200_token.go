package core

import (
	"sync"
)

// CharityCampaign tracks donations and progress for a specific symbol.
type CharityCampaign struct {
	Symbol    string
	Purpose   string
	Goal      uint64
	Raised    uint64
	Donations map[string]uint64
}

// SYN4200Token implements charity token functionality supporting donations and progress checks.
type SYN4200Token struct {
	mu        sync.RWMutex
	campaigns map[string]*CharityCampaign
}

// NewSYN4200Token creates an empty charity token registry.
func NewSYN4200Token() *SYN4200Token {
	return &SYN4200Token{campaigns: make(map[string]*CharityCampaign)}
}

// Donate records a donation to a campaign. Creating the campaign if it does not exist.
func (t *SYN4200Token) Donate(symbol, from string, amount uint64, purpose string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	c, ok := t.campaigns[symbol]
	if !ok {
		c = &CharityCampaign{
			Symbol:    symbol,
			Purpose:   purpose,
			Donations: make(map[string]uint64),
		}
		t.campaigns[symbol] = c
	}
	c.Raised += amount
	c.Donations[from] += amount
}

// CampaignProgress returns the amount raised for a campaign.
func (t *SYN4200Token) CampaignProgress(symbol string) (uint64, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	c, ok := t.campaigns[symbol]
	if !ok {
		return 0, false
	}
	return c.Raised, true
}

// Campaign returns a copy of the campaign data for inspection.
func (t *SYN4200Token) Campaign(symbol string) (*CharityCampaign, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	c, ok := t.campaigns[symbol]
	if !ok {
		return nil, false
	}
	cp := *c
	cp.Donations = make(map[string]uint64, len(c.Donations))
	for k, v := range c.Donations {
		cp.Donations[k] = v
	}
	return &cp, true
}
