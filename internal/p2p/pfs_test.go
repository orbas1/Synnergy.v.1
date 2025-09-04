package p2p

import "testing"

func TestPFSChannel(t *testing.T) {
	ch := NewPFSChannel()
	enc, err := ch.Encrypt([]byte("hello"))
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	dec, err := ch.Decrypt(enc)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	if string(dec) != "hello" {
		t.Fatalf("unexpected value: %s", string(dec))
	}
}
