package p2p

import "errors"

// PFSChannel simulates a peer-to-peer channel with perfect forward secrecy.
type PFSChannel struct{}

// NewPFSChannel creates a new PFSChannel.
func NewPFSChannel() *PFSChannel { return &PFSChannel{} }

// Encrypt is a placeholder for message encryption.
func (c *PFSChannel) Encrypt(msg []byte) ([]byte, error) {
	if len(msg) == 0 {
		return nil, errors.New("empty message")
	}
	return append([]byte(nil), msg...), nil
}

// Decrypt is a placeholder for message decryption.
func (c *PFSChannel) Decrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	return append([]byte(nil), data...), nil
}
