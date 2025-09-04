package security

import "testing"

func TestEncryptor(t *testing.T) {
	e := NewEncryptor(0xAA)
	enc := e.Encrypt([]byte{0x01})
	dec := e.Decrypt(enc)
	if dec[0] != 0x01 {
		t.Fatalf("expected 0x01")
	}
}
