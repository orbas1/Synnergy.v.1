package tokens

import (
	"encoding/hex"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestBaseTokenMintTransferBurn(t *testing.T) {
	tok := NewBaseToken(1, "Test", "TST", 0)
	if err := tok.Mint("alice", 100); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := tok.Transfer("alice", "bob", 40); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if tok.BalanceOf("bob") != 40 {
		t.Fatalf("expected bob=40 got %d", tok.BalanceOf("bob"))
	}
	if err := tok.Burn("alice", 10); err != nil {
		t.Fatalf("burn: %v", err)
	}
	if tok.TotalSupply() != 90 {
		t.Fatalf("unexpected supply %d", tok.TotalSupply())
	}
}

func TestBaseTokenAllowance(t *testing.T) {
	tok := NewBaseToken(10, "Allow", "ALL", 0)
	if err := tok.Mint("alice", 100); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := tok.Approve("alice", "bob", 50); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if tok.Allowance("alice", "bob") != 50 {
		t.Fatalf("allowance wrong")
	}
	if err := tok.TransferFrom("alice", "bob", "carol", 30); err != nil {
		t.Fatalf("transferFrom: %v", err)
	}
	if tok.BalanceOf("carol") != 30 || tok.Allowance("alice", "bob") != 20 {
		t.Fatalf("post transfer state incorrect")
	}
}

func TestSYN10Info(t *testing.T) {
	tkn := NewSYN10Token(2, "CBDC", "CBD", "Central Bank", 1.0, 2)
	if err := tkn.Mint("alice", 50); err != nil {
		t.Fatalf("mint: %v", err)
	}
	info := tkn.Info()
	if info.Issuer != "Central Bank" || info.ExchangeRate != 1.0 {
		t.Fatalf("unexpected info %+v", info)
	}
}

func TestSYN1000ReserveValue(t *testing.T) {
	idx := NewSYN1000Index()
	id := idx.Create("Stable", "STBL", 2)
	if err := idx.AddReserve(id, "USD", 100); err != nil {
		t.Fatalf("add reserve: %v", err)
	}
	if err := idx.SetReservePrice(id, "USD", 1.0); err != nil {
		t.Fatalf("set price: %v", err)
	}
	val, err := idx.TotalValue(id)
	if err != nil || val != 100 {
		t.Fatalf("unexpected value %f err %v", val, err)
	}
}

func TestSYN12Metadata(t *testing.T) {
	meta := SYN12Metadata{BillID: "TB1", Issuer: "Treasury", IssueDate: time.Now(), Maturity: time.Now().Add(24 * time.Hour), Discount: 0.02, FaceValue: 1000}
	tkn := NewSYN12Token(3, "TBill", "TB", meta, 2)
	if tkn.Metadata.Issuer != "Treasury" || tkn.Metadata.FaceValue != 1000 {
		t.Fatalf("unexpected metadata: %+v", tkn.Metadata)
	}
}

func TestSYN20PauseFreeze(t *testing.T) {
	tok := NewSYN20Token(4, "Token", "S20", 0)
	if err := tok.Mint("alice", 10); err != nil {
		t.Fatalf("mint: %v", err)
	}
	tok.Freeze("alice")
	if err := tok.Transfer("alice", "bob", 1); err == nil {
		t.Fatalf("expected frozen error")
	}
	tok.Unfreeze("alice")
	tok.Pause()
	if err := tok.Transfer("alice", "bob", 1); err == nil {
		t.Fatalf("expected paused error")
	}
	tok.Unpause()
	if err := tok.Transfer("alice", "bob", 1); err != nil {
		t.Fatalf("transfer: %v", err)
	}
}

func TestSYN70AssetLifecycle(t *testing.T) {
	tok := NewSYN70Token(5, "Game", "G70", 0)
	if err := tok.RegisterAsset("a1", "alice", "Sword", "RPG"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := tok.SetAttribute("a1", "damage", "10"); err != nil {
		t.Fatalf("setattr: %v", err)
	}
	if err := tok.AddAchievement("a1", "FirstBlood"); err != nil {
		t.Fatalf("achievement: %v", err)
	}
	if err := tok.TransferAsset("a1", "bob"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	info, err := tok.AssetInfo("a1")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if info.Owner != "bob" || info.Attributes["damage"] != "10" || len(info.Achievements) != 1 {
		t.Fatalf("unexpected asset info: %+v", info)
	}
	if len(tok.ListAssets()) != 1 {
		t.Fatalf("expected one asset")
	}
}

func TestRegistryInfoList(t *testing.T) {
	reg := NewRegistry()
	tok := NewBaseToken(reg.NextID(), "Tkn", "TK", 0)
	reg.Register(tok)
	if len(reg.List()) != 1 {
		t.Fatalf("expected one token in list")
	}
	info, ok := reg.Info(tok.ID())
	if !ok || info.Symbol != "TK" {
		t.Fatalf("unexpected info %+v", info)
	}
}

func TestSYN1100Access(t *testing.T) {
	store := NewSYN1100Token()
	data, _ := hex.DecodeString("abcd")
	if err := store.AddRecord(1, "alice", data); err != nil {
		t.Fatalf("add record: %v", err)
	}
	if err := store.GrantAccess(1, "bob"); err != nil {
		t.Fatalf("grant: %v", err)
	}
	if _, err := store.GetRecord(1, "bob"); err != nil {
		t.Fatalf("bob should access: %v", err)
	}
	if err := store.RevokeAccess(1, "bob"); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := store.GetRecord(1, "bob"); err == nil {
		t.Fatalf("expected access denied")
	}
}

func TestBaseTokenTransferErrors(t *testing.T) {
	tok := NewBaseToken(1, "Err", "ERR", 0)
	if err := tok.Transfer("alice", "bob", 1); !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance got %v", err)
	}
	if err := tok.TransferFrom("alice", "bob", "carol", 1); !errors.Is(err, ErrAllowanceExceeded) {
		t.Fatalf("expected ErrAllowanceExceeded got %v", err)
	}
}

func TestBaseTokenConcurrent(t *testing.T) {
	tok := NewBaseToken(2, "Conc", "CON", 0)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := tok.Mint("alice", 1); err != nil {
				t.Errorf("mint: %v", err)
			}
		}()
	}
	wg.Wait()
	if tok.BalanceOf("alice") != 100 {
		t.Fatalf("expected balance 100 got %d", tok.BalanceOf("alice"))
	}
}
