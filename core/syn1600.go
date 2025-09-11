package core

import (
	"errors"
	"sync"
)

var (
	ErrNoRoyaltyRecipients = errors.New("no royalty recipients")
	ErrRecipientNotFound   = errors.New("royalty recipient not found")
)

// MusicToken represents metadata and royalty distribution for a SYN1600 token.
type MusicToken struct {
	Title         string
	Artist        string
	Album         string
	mu            sync.RWMutex
	royaltySplits map[string]uint64
	totalShares   uint64
}

// NewMusicToken initialises a music token with basic metadata.
func NewMusicToken(title, artist, album string) *MusicToken {
	return &MusicToken{
		Title:         title,
		Artist:        artist,
		Album:         album,
		royaltySplits: make(map[string]uint64),
	}
}

// Info returns the music metadata.
func (m *MusicToken) Info() (string, string, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Title, m.Artist, m.Album
}

// Update modifies the music metadata. Empty values are ignored.
func (m *MusicToken) Update(title, artist, album string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if title != "" {
		m.Title = title
	}
	if artist != "" {
		m.Artist = artist
	}
	if album != "" {
		m.Album = album
	}
}

// SetRoyaltyShare sets or updates the royalty share for an address. A share of
// zero removes the recipient.
func (m *MusicToken) SetRoyaltyShare(addr string, share uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalShares -= m.royaltySplits[addr]
	if share == 0 {
		delete(m.royaltySplits, addr)
		return
	}
	m.royaltySplits[addr] = share
	m.totalShares += share
}

// RoyaltyShare retrieves the share for the given address.
func (m *MusicToken) RoyaltyShare(addr string) (uint64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	share, ok := m.royaltySplits[addr]
	return share, ok
}

// RemoveRoyaltyRecipient removes a royalty recipient entirely.
func (m *MusicToken) RemoveRoyaltyRecipient(addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	share, ok := m.royaltySplits[addr]
	if !ok {
		return ErrRecipientNotFound
	}
	m.totalShares -= share
	delete(m.royaltySplits, addr)
	return nil
}

// Distribute calculates payouts for each royalty recipient based on shares.
func (m *MusicToken) Distribute(amount uint64) (map[string]uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.totalShares == 0 {
		return nil, ErrNoRoyaltyRecipients
	}
	payouts := make(map[string]uint64, len(m.royaltySplits))
	for addr, share := range m.royaltySplits {
		payouts[addr] = amount * share / m.totalShares
	}
	return payouts, nil
}
