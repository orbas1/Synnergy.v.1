package cli

import "testing"

func TestZeroTrustDataChannelsFlow(t *testing.T) {
	key := make([]byte, 32)
	if err := ztEngine.OpenChannel("ch", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	cipher, err := ztEngine.Send("ch", []byte("hi"))
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if len(cipher) == 0 {
		t.Fatalf("cipher empty")
	}
	if _, err := ztEngine.Receive("ch", 0); err != nil {
		t.Fatalf("receive: %v", err)
	}
	if err := ztEngine.CloseChannel("ch"); err != nil {
		t.Fatalf("close: %v", err)
	}
}
