package synnergy

import "errors"

// BiometricSecurityNode couples a node identifier with biometric authentication
// to protect privileged operations.
type BiometricSecurityNode struct {
	id   string
	Auth *BiometricsAuth
}

// NewBiometricSecurityNode creates a node secured by biometric authentication.
// If auth is nil a new instance is created.
func NewBiometricSecurityNode(id string, auth *BiometricsAuth) *BiometricSecurityNode {
	if auth == nil {
		auth = NewBiometricsAuth()
	}
	return &BiometricSecurityNode{id: id, Auth: auth}
}

// GetID returns the identifier of the node.
func (b *BiometricSecurityNode) GetID() string { return b.id }

// Enroll registers biometric data for the given address.
func (b *BiometricSecurityNode) Enroll(addr string, biometric []byte) {
	b.Auth.Enroll(addr, biometric)
}

// Remove deletes biometric data associated with the address.
func (b *BiometricSecurityNode) Remove(addr string) {
	b.Auth.Remove(addr)
}

// Authenticate verifies biometric data for the address.
func (b *BiometricSecurityNode) Authenticate(addr string, biometric []byte) bool {
	return b.Auth.Verify(addr, biometric)
}

// SecureExecute runs fn only if biometric verification succeeds for the address.
func (b *BiometricSecurityNode) SecureExecute(addr string, biometric []byte, fn func() error) error {
	if !b.Auth.Verify(addr, biometric) {
		return errors.New("biometric verification failed")
	}
	if fn != nil {
		return fn()
	}
	return nil
}
