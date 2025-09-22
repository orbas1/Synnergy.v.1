package tokens

import (
	"errors"
	"testing"
)

func TestSYN223TransferWithHooks(t *testing.T) {
	tok := NewSYN223Token("SYN", "SYN223", "issuer", 1_000)
	tok.AddToWhitelist("issuer")
	tok.AddToWhitelist("receiver")
	tok.SetMetadata("receiver", map[string]string{"kyc": "level3"})

	var called bool
	tok.RegisterHook("receiver", func(from string, amount uint64, data []byte) error {
		called = true
		if string(data) != "ok" {
			return errors.New("unexpected payload")
		}
		return nil
	})

	if err := tok.TransferWithData("issuer", "receiver", 200, []byte("ok"), "settlement"); err != nil {
		t.Fatalf("transfer with hook: %v", err)
	}
	if !called {
		t.Fatal("expected hook invocation")
	}
	if tok.BalanceOf("receiver") != 200 {
		t.Fatalf("unexpected balance: %d", tok.BalanceOf("receiver"))
	}

	events := tok.Events(5)
	if len(events) != 1 || events[0].Memo != "settlement" {
		t.Fatalf("unexpected events: %+v", events)
	}

	meta := tok.Metadata("receiver")
	if meta["kyc"] != "level3" {
		t.Fatalf("unexpected metadata: %+v", meta)
	}

	tok.RegisterHook("receiver", func(from string, amount uint64, data []byte) error {
		return errors.New("blocked")
	})
	if err := tok.TransferWithData("issuer", "receiver", 1, nil, "fail"); err == nil {
		t.Fatal("expected hook failure to abort transfer")
	}

	tok.AddToBlacklist("receiver")
	if err := tok.Transfer("issuer", "receiver", 1); err == nil {
		t.Fatal("expected blacklist rejection")
	}
}
