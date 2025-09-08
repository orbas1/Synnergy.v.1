package synnergy

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

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

// NewContentMetaFromData builds a ContentMeta from raw data, computing the
// size and SHA-256 hash automatically. It mirrors the behaviour expected by
// CLI content upload operations to avoid manual bookkeeping by callers.
func NewContentMetaFromData(id, name string, data []byte) ContentMeta {
	digest := sha256.Sum256(data)
	return ContentMeta{
		ID:      id,
		Name:    name,
		Size:    int64(len(data)),
		Hash:    hex.EncodeToString(digest[:]),
		Created: time.Now().UTC(),
	}
}

// Validate performs basic sanity checks on the metadata ensuring fields are
// populated and structurally valid. It is intended to be run prior to
// persisting metadata on-chain or advertising it over the network.
func (cm ContentMeta) Validate() error {
	switch {
	case cm.ID == "":
		return errors.New("id cannot be empty")
	case cm.Name == "":
		return errors.New("name cannot be empty")
	case cm.Size < 0:
		return errors.New("size cannot be negative")
	case len(cm.Hash) != 64:
		return errors.New("hash must be a 64-character hex string")
	}
	return nil
}
