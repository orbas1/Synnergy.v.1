package core

import "errors"

// MusicToken represents metadata and royalty distribution for a SYN1600 token.
type MusicToken struct {
	Title  string
	Artist string
	Album  string

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
	return m.Title, m.Artist, m.Album
}

// Update modifies the music metadata. Empty values are ignored.
func (m *MusicToken) Update(title, artist, album string) {
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

// SetRoyaltyShare sets the royalty share for an address.
func (m *MusicToken) SetRoyaltyShare(addr string, share uint64) {
	m.totalShares -= m.royaltySplits[addr]
	m.royaltySplits[addr] = share
	m.totalShares += share
}

// Distribute calculates payouts for each royalty recipient based on shares.
func (m *MusicToken) Distribute(amount uint64) (map[string]uint64, error) {
	if m.totalShares == 0 {
		return nil, errors.New("no royalty recipients")
	}
	payouts := make(map[string]uint64, len(m.royaltySplits))
	for addr, share := range m.royaltySplits {
		payouts[addr] = amount * share / m.totalShares
	}
	return payouts, nil
}
