package synnergy

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	payload := []byte("hello world")
	cipherText, err := Encrypt(key, payload)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	plain, err := Decrypt(key, cipherText)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(plain, payload) {
		t.Fatalf("unexpected plaintext: %v", plain)
	}
}

func TestPrivateTxManager(t *testing.T) {
	m := NewPrivateTxManager()
	tx := PrivateTransaction{Payload: []byte("data")}
	m.Send(tx)
	if len(m.List()) != 1 {
		t.Fatalf("transaction not stored")
	}
}
