package security

// Encryptor provides trivial XOR-based encryption for example purposes.
type Encryptor struct {
	key byte
}

// NewEncryptor creates an Encryptor with the given single-byte key.
func NewEncryptor(key byte) *Encryptor { return &Encryptor{key: key} }

// Encrypt XORs the data with the key.
func (e *Encryptor) Encrypt(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ e.key
	}
	return out
}

// Decrypt XORs the data with the key (same as Encrypt for XOR cipher).
func (e *Encryptor) Decrypt(data []byte) []byte {
	return e.Encrypt(data)
}
