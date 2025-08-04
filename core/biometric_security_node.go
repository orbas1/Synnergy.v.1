package core

import "errors"

// BiometricSecurityNode couples a Node with biometric authentication to protect
// administrative operations.
type BiometricSecurityNode struct {
	*Node
	Auth *BiometricsAuth
}

// NewBiometricSecurityNode creates a node secured by biometric authentication.
// If auth is nil a new instance is created.
func NewBiometricSecurityNode(base *Node, auth *BiometricsAuth) *BiometricSecurityNode {
	if auth == nil {
		auth = NewBiometricsAuth()
	}
	return &BiometricSecurityNode{Node: base, Auth: auth}
}

// GetID returns the identifier of the embedded node.
func (b *BiometricSecurityNode) GetID() string {
	if b.Node != nil {
		return b.Node.ID
	}
	return ""
}

// Enroll registers biometric data for the given address.
func (b *BiometricSecurityNode) Enroll(addr string, biometric []byte) {
	b.Auth.Enroll(addr, biometric)
}

// Authenticate verifies biometric data for the address.
func (b *BiometricSecurityNode) Authenticate(addr string, biometric []byte) bool {
	return b.Auth.Verify(addr, biometric)
}

// SecureAddTransaction adds a transaction to the node's mempool only if the
// biometric data matches the enrolled template for the provided address.
func (b *BiometricSecurityNode) SecureAddTransaction(addr string, biometric []byte, tx *Transaction) error {
	if !b.Auth.Verify(addr, biometric) {
		return errors.New("biometric verification failed")
	}
	return b.AddTransaction(tx)
}
