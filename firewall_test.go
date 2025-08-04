package synnergy

import "testing"

func TestFirewall(t *testing.T) {
	f := NewFirewall()

	f.BlockAddress("addr1")
	if !f.IsAddressBlocked("addr1") {
		t.Fatalf("address should be blocked")
	}
	f.UnblockAddress("addr1")
	if f.IsAddressBlocked("addr1") {
		t.Fatalf("address should be unblocked")
	}

	f.BlockToken("token1")
	if !f.IsTokenBlocked("token1") {
		t.Fatalf("token should be blocked")
	}

	f.BlockIP("1.2.3.4")
	if !f.IsIPBlocked("1.2.3.4") {
		t.Fatalf("ip should be blocked")
	}

	addrs, tokens, ips := f.Rules()
	if len(addrs) != 0 {
		t.Fatalf("expected no blocked addresses, got %d", len(addrs))
	}
	if len(tokens) != 1 || tokens[0] != "token1" {
		t.Fatalf("unexpected tokens list: %v", tokens)
	}
	if len(ips) != 1 || ips[0] != "1.2.3.4" {
		t.Fatalf("unexpected ips list: %v", ips)
	}
}
