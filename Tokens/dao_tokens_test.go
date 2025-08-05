package tokens

import (
	"testing"
	"time"
)

func TestSYN223TokenTransfer(t *testing.T) {
	tok := NewSYN223Token("Test", "TST", "alice", 100)
	tok.AddToWhitelist("bob")
	if err := tok.Transfer("alice", "bob", 40); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if tok.BalanceOf("bob") != 40 {
		t.Fatalf("expected bob=40 got %d", tok.BalanceOf("bob"))
	}
	if err := tok.Transfer("alice", "carol", 10); err == nil {
		t.Fatalf("expected error for non-whitelisted recipient")
	}
}

func TestSYN300TokenGovernanceFlow(t *testing.T) {
	tok := NewSYN300Token(map[string]uint64{"alice": 100, "bob": 50})
	tok.Delegate("bob", "alice")
	id := tok.CreateProposal("alice", "test proposal")
	if err := tok.Vote(id, "alice", true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := tok.Execute(id, 120); err != nil {
		t.Fatalf("execute: %v", err)
	}
	prop, err := tok.ProposalStatus(id)
	if err != nil || !prop.Executed {
		t.Fatalf("proposal not executed: %v", err)
	}
}

func TestSyn2500Registry(t *testing.T) {
	reg := NewSyn2500Registry()
	m := NewSyn2500Member("1", "alice", 10, nil)
	reg.AddMember(m)
	got, ok := reg.GetMember("1")
	if !ok || got.Address != "alice" {
		t.Fatalf("member retrieval failed")
	}
	got.UpdateVotingPower(20)
	if got.VotingPower != 20 {
		t.Fatalf("update voting power failed")
	}
	reg.RemoveMember("1")
	if _, ok := reg.GetMember("1"); ok {
		t.Fatalf("expected member removed")
	}
}

func TestSYN3500TokenMintRedeem(t *testing.T) {
	tok := NewSYN3500Token("Currency", "CUR", "issuer", 1.0)
	tok.Mint("alice", 100)
	if err := tok.Redeem("alice", 40); err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if bal := tok.BalanceOf("alice"); bal != 60 {
		t.Fatalf("unexpected balance %d", bal)
	}
}

func TestSYN3700TokenValue(t *testing.T) {
	tok := NewSYN3700Token("Index", "IDX")
	tok.AddComponent("AAA", 1.0)
	tok.AddComponent("BBB", 2.0)
	prices := map[string]float64{"AAA": 2.0, "BBB": 3.0}
	if v := tok.Value(prices); v != 8.0 {
		t.Fatalf("unexpected value %f", v)
	}
}

func TestSYN4200TokenDonations(t *testing.T) {
	tok := NewSYN4200Token()
	tok.Donate("CHAR", "alice", 10, "help")
	if prog, ok := tok.CampaignProgress("CHAR"); !ok || prog != 10 {
		t.Fatalf("unexpected progress %d ok %v", prog, ok)
	}
	camp, ok := tok.Campaign("CHAR")
	if !ok || camp.Raised != 10 {
		t.Fatalf("unexpected campaign data")
	}
}

func TestLegalTokenWorkflow(t *testing.T) {
	lt := NewLegalToken("1", "Contract", "LGL", "doc", "hash", "owner", time.Now().Add(time.Hour), 1, []string{"alice", "bob"})
	if err := lt.Sign("alice", "sig1"); err != nil {
		t.Fatalf("sign: %v", err)
	}
	lt.UpdateStatus(LegalTokenStatusActive)
	lt.Dispute("challenge", "pending")
	if lt.Status != LegalTokenStatusDisputed || len(lt.Disputes) != 1 {
		t.Fatalf("dispute not recorded")
	}
}
