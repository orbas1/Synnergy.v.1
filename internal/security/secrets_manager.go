package security

// SecretsManager provides a simple in-memory key-value store for secrets.
type SecretsManager struct {
	store map[string]string
}

// NewSecretsManager creates a new SecretsManager.
func NewSecretsManager() *SecretsManager {
	return &SecretsManager{store: make(map[string]string)}
}

// Store saves a secret value for the given key.
func (s *SecretsManager) Store(key, value string) {
	s.store[key] = value
}

// Retrieve gets a secret value by key.
func (s *SecretsManager) Retrieve(key string) (string, bool) {
	v, ok := s.store[key]
	return v, ok
}
