package synnergy

import "time"

// ContentMeta describes a piece of content maintained by a ContentNode.
// It mirrors metadata pinned on-chain allowing nodes to advertise and
// validate content availability.
type ContentMeta struct {
	// ID is a unique identifier derived from the content payload.
	ID string
	// Name provides a human readable label for the content.
	Name string
	// Size reports the number of bytes in the original data.
	Size int64
	// Hash is a SHA-256 hash of the plaintext content encoded as hex.
	Hash string
	// Created indicates when the content was first stored.
	Created time.Time
}

// NewContentMeta constructs a ContentMeta from basic fields and sets the
// creation timestamp to the current UTC time.
func NewContentMeta(id, name string, size int64, hash string) ContentMeta {
	return ContentMeta{
		ID:      id,
		Name:    name,
		Size:    size,
		Hash:    hash,
		Created: time.Now().UTC(),
	}
}
