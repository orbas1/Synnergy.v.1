package cli

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"synnergy/core"
)

func TestZeroTrustDataChannelsFlow(t *testing.T) {
	key := make([]byte, 32)
	if err := ztEngine.OpenChannel("ch", key, core.WithOwner("ops")); err != nil {
		t.Fatalf("open: %v", err)
	}
	cipher, err := ztEngine.SendAs("ch", "ops", []byte("hi"))
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if len(cipher) == 0 {
		t.Fatalf("cipher empty")
	}
	if _, err := ztEngine.Receive("ch", 0); err != nil {
		t.Fatalf("receive: %v", err)
	}
	_, pub, _ := ed25519.GenerateKey(nil)
	if err := ztEngine.AuthorizePeer("ch", "ally", hex.EncodeToString(pub)); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if _, err := ztEngine.SendAs("ch", "ally", []byte("ally")); err != nil {
		t.Fatalf("send ally: %v", err)
	}
	if err := ztEngine.CloseChannel("ch"); err != nil {
		t.Fatalf("close: %v", err)
	}
}
