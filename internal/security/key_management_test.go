package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"
	"testing"
)

func TestKeyManagerLifecycle(t *testing.T) {
	km := NewKeyManager()
	version, sym, err := km.GenerateSymmetricKey(PurposeNoiseStatic, "test")
	if err != nil {
		t.Fatalf("generate symmetric: %v", err)
	}
	if version != 1 || len(sym) != 32 {
		t.Fatalf("unexpected symmetric output")
	}
	if err := km.SetSymmetricKey(PurposeEnvelope, []byte("abcd"), "restore", 3); err != nil {
		t.Fatalf("set symmetric: %v", err)
	}
	got, v, err := km.SymmetricKey(PurposeEnvelope)
	if err != nil || v != 3 || string(got) != "abcd" {
		t.Fatalf("unexpected symmetric retrieval: %v %d %s", err, v, string(got))
	}

	pub, sigVersion, err := km.GenerateSigningKey(PurposeStateSigning, "test")
	if err != nil {
		t.Fatalf("generate signing: %v", err)
	}
	if sigVersion != 1 || len(pub) != ed25519.PublicKeySize {
		t.Fatalf("unexpected signing output")
	}
	sig, _, _, err := km.Sign(PurposeStateSigning, []byte("message"))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := km.Verify(PurposeStateSigning, []byte("message"), sig, nil); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := km.Verify(PurposeStateSigning, []byte("other"), sig, nil); err == nil {
		t.Fatalf("expected verification failure")
	}

	log := km.AuditLog()
	if len(log) < 2 {
		t.Fatalf("expected audit entries")
	}
}

func TestKeyManagerDeterministicEntropy(t *testing.T) {
	km := NewKeyManager()
	buf := make([]byte, 64)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("entropy: %v", err)
	}
	reader := NewDeterministicReader(buf)
	km.WithEntropy(reader)
	version, key, err := km.GenerateSymmetricKey(PurposeNoiseStatic, "det")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if version != 1 || len(key) != 32 {
		t.Fatalf("unexpected key output")
	}
}

// deterministic reader for testing

type deterministicReader struct {
	data []byte
	off  int
}

func NewDeterministicReader(data []byte) *deterministicReader {
	return &deterministicReader{data: append([]byte(nil), data...)}
}

func (r *deterministicReader) Read(p []byte) (int, error) {
	n := copy(p, r.data[r.off:])
	r.off += n
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}
