package security

import "errors"

// SecretsManager provides a simple in-memory key-value store for secrets.
type SecretsManager struct {
	store map[string]string
}

// NewSecretsManager creates a new SecretsManager.
func NewSecretsManager() *SecretsManager {
	return &SecretsManager{store: make(map[string]string)}
}

// Store saves a secret value for the given key and returns an error if the
// key or value is empty.
func (s *SecretsManager) Store(key, value string) error {
	if key == "" {
		return errors.New("key required")
	}
	if value == "" {
		return errors.New("value required")
	}
	s.store[key] = value
	return nil
}

// Retrieve gets a secret value by key. An error is returned if the key is empty
// or not found.
func (s *SecretsManager) Retrieve(key string) (string, error) {
	if key == "" {
		return "", errors.New("key required")
	}
	v, ok := s.store[key]
	if !ok {
		return "", errors.New("secret not found")
	}
	return v, nil
}
