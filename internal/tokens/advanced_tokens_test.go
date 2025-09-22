package tokens

import (
	"testing"
	"time"
)

func TestSYN200CarbonRegistry(t *testing.T) {
	reg := NewCarbonRegistry()
	proj := reg.Register("alice", "Forest", 100)
	if proj.ID == "" {
		t.Fatalf("expected project id")
	}
	if err := reg.Issue(proj.ID, "bob", 50); err != nil {
		t.Fatalf("issue: %v", err)
	}
	if err := reg.Retire(proj.ID, "bob", 20); err != nil {
		t.Fatalf("retire: %v", err)
	}
	if err := reg.AddVerification(proj.ID, "verifier", "V1", "ok"); err != nil {
		t.Fatalf("verify: %v", err)
	}
	ver, ok := reg.Verifications(proj.ID)
	if !ok || len(ver) != 1 {
		t.Fatalf("expected 1 verification")
	}
	info, ok := reg.ProjectInfo(proj.ID)
	if !ok {
		t.Fatalf("project info missing")
	}
	if info.IssuedCredits != 30 || info.RetiredCredits != 20 {
		t.Fatalf("unexpected credit counts: %+v", info)
	}
	projects := reg.ListProjects()
	if len(projects) != 1 {
		t.Fatalf("expected 1 project got %d", len(projects))
	}
}

func TestSYN2600InvestorRegistry(t *testing.T) {
	reg := NewInvestorRegistry()
	expiry := time.Now().Add(24 * time.Hour)
	tok, err := reg.Issue("AssetA", "alice", 100, expiry)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if err := reg.Transfer(tok.ID, "bob"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if err := reg.RecordReturn(tok.ID, 10); err != nil {
		t.Fatalf("record return: %v", err)
	}
	got, ok := reg.Get(tok.ID)
	if !ok || got.Owner != "bob" || len(got.Returns) != 1 {
		t.Fatalf("unexpected token data: %+v", got)
	}
	tokens := reg.List()
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token got %d", len(tokens))
	}
	if _, err := reg.Issue("", "", 0, time.Now()); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSYN2800LifePolicyRegistry(t *testing.T) {
	reg := NewLifePolicyRegistry()
	start := time.Now()
	end := start.Add(365 * 24 * time.Hour)
	policy, err := reg.IssuePolicy("alice", "bob", 1000, 100, start, end)
	if err != nil {
		t.Fatalf("issue policy: %v", err)
	}
	if err := reg.PayPremium(policy.PolicyID, 50); err != nil {
		t.Fatalf("pay premium: %v", err)
	}
	if _, err := reg.FileClaim(policy.PolicyID, 200); err != nil {
		t.Fatalf("file claim: %v", err)
	}
	info, ok := reg.GetPolicy(policy.PolicyID)
	if !ok || info.PaidPremium != 50 || len(info.Claims) != 1 {
		t.Fatalf("unexpected policy info: %+v", info)
	}
	policies := reg.ListPolicies()
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy got %d", len(policies))
	}
	if _, err := reg.IssuePolicy("", "", 0, 0, end, start); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSYN2900InsuranceRegistry(t *testing.T) {
	reg := NewInsuranceRegistry()
	start := time.Now()
	end := start.Add(365 * 24 * time.Hour)
	policy, err := reg.IssuePolicy("alice", "fire", 100, 1000, 10, 1000, start, end)
	if err != nil {
		t.Fatalf("issue policy: %v", err)
	}
	if _, err := reg.FileClaim(policy.PolicyID, "damage", 500); err != nil {
		t.Fatalf("file claim: %v", err)
	}
	info, ok := reg.GetPolicy(policy.PolicyID)
	if !ok || len(info.Claims) != 1 {
		t.Fatalf("unexpected policy info: %+v", info)
	}
	policies := reg.ListPolicies()
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy got %d", len(policies))
	}
	if _, err := reg.IssuePolicy("", "", 0, 0, 0, 0, end, start); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSYN3400ForexRegistry(t *testing.T) {
	reg := NewForexRegistry()
	pair, err := reg.Register("USD", "EUR", 1.1)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := reg.Register("USD", "EUR", 1.15); err != nil {
		t.Fatalf("re-register: %v", err)
	}
	if err := reg.UpdateRate(pair.PairID, 1.2); err != nil {
		t.Fatalf("update rate: %v", err)
	}
	got, err := reg.Get(pair.PairID)
	if err != nil || got.Rate != 1.2 {
		t.Fatalf("unexpected pair data: %+v err:%v", got, err)
	}
	symbol, err := reg.GetBySymbol("USD", "EUR")
	if err != nil || symbol.Rate != 1.2 {
		t.Fatalf("unexpected symbol lookup: %+v err:%v", symbol, err)
	}
	pairs := reg.List()
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair got %d", len(pairs))
	}
	if err := reg.Remove(pair.PairID); err != nil {
		t.Fatalf("remove pair: %v", err)
	}
}

func TestSYN845DebtRegistry(t *testing.T) {
	reg := NewDebtRegistry()
	tokenID, _ := reg.CreateToken("DebtToken", "DBT", "alice", 1000)
	due := time.Now().Add(24 * time.Hour)
	if err := reg.IssueDebt(tokenID, "D1", "bob", 500, 0.05, 0.02, due); err != nil {
		t.Fatalf("issue debt: %v", err)
	}
	if err := reg.RecordPayment(tokenID, "D1", 100); err != nil {
		t.Fatalf("record payment: %v", err)
	}
	debt, err := reg.GetDebt(tokenID, "D1")
	if err != nil || debt.Paid != 100 {
		t.Fatalf("unexpected debt data: %+v err:%v", debt, err)
	}
}

func TestSYN2369ItemRegistry(t *testing.T) {
	reg := NewItemRegistry()
	item := reg.CreateItem("alice", "Sword", "desc", map[string]string{"atk": "10"})
	if err := reg.TransferItem(item.ItemID, "bob"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if err := reg.UpdateAttributes(item.ItemID, map[string]string{"atk": "15"}); err != nil {
		t.Fatalf("update attrs: %v", err)
	}
	got, ok := reg.GetItem(item.ItemID)
	if !ok || got.Owner != "bob" || got.Attributes["atk"] != "15" {
		t.Fatalf("unexpected item data: %+v", got)
	}
	items := reg.ListItems()
	if len(items) != 1 {
		t.Fatalf("expected 1 item got %d", len(items))
	}
}

func TestSYN223Token(t *testing.T) {
	tok := NewSYN223Token("Token", "TKN", "alice", 100)
	tok.AddToWhitelist("bob")
	if err := tok.Transfer("alice", "bob", 20); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if bal := tok.BalanceOf("bob"); bal != 20 {
		t.Fatalf("unexpected balance %d", bal)
	}
	tok.AddToBlacklist("bob")
	if err := tok.Transfer("alice", "bob", 10); err == nil {
		t.Fatalf("expected blacklist error")
	}
	tok.RemoveFromBlacklist("bob")
	if err := tok.Transfer("alice", "bob", 10); err != nil {
		t.Fatalf("transfer2: %v", err)
	}
	tok.RemoveFromWhitelist("bob")
	if err := tok.Transfer("alice", "bob", 1); err == nil {
		t.Fatalf("expected whitelist error")
	}
}

func TestSyn2500RegistryOperations(t *testing.T) {
	reg := NewSyn2500Registry()
	m := NewSyn2500Member("1", "alice", 10, map[string]string{"role": "admin"})
	reg.AddMember(m)
	got, ok := reg.GetMember("1")
	if !ok || got.Address != "alice" || got.Metadata["role"] != "admin" {
		t.Fatalf("unexpected member data: %+v", got)
	}
	m.UpdateVotingPower(20)
	got, _ = reg.GetMember("1")
	if got.VotingPower != 20 {
		t.Fatalf("unexpected voting power %d", got.VotingPower)
	}
	members := reg.ListMembers()
	if len(members) != 1 {
		t.Fatalf("expected 1 member got %d", len(members))
	}
	reg.RemoveMember("1")
	if _, ok := reg.GetMember("1"); ok {
		t.Fatalf("member not removed")
	}
}

func TestSYN300Token(t *testing.T) {
	tok := NewSYN300Token(map[string]uint64{"alice": 100, "bob": 50})
	tok.Delegate("bob", "alice")
	if pow := tok.VotingPower("alice"); pow != 150 {
		t.Fatalf("unexpected voting power %d", pow)
	}
	tok.RevokeDelegation("bob")
	if pow := tok.VotingPower("alice"); pow != 100 {
		t.Fatalf("unexpected power after revoke %d", pow)
	}
	tok.Delegate("bob", "alice")
	id, err := tok.CreateProposal("alice", "test")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := tok.Vote(id, "alice", true); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := tok.Execute(id, 150); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if err := tok.Execute(id, 150); err == nil {
		t.Fatalf("expected execute error")
	}
	status, err := tok.ProposalStatus(id)
	if err != nil || !status.Executed {
		t.Fatalf("proposal not executed")
	}
	id2, err := tok.CreateProposal("alice", "second")
	if err != nil {
		t.Fatalf("create proposal: %v", err)
	}
	if err := tok.Vote(id2, "bob", true); err != nil {
		t.Fatalf("vote2: %v", err)
	}
	if err := tok.Execute(id2, 60); err == nil {
		t.Fatalf("expected quorum error")
	}
}

func TestSYN3500Token(t *testing.T) {
	tok := NewSYN3500Token("USD Token", "USDT", "issuer", 1.0)
	tok.Mint("alice", 100)
	tok.SetRate(1.2)
	sym, issuer, rate := tok.Info()
	if sym != "USDT" || issuer != "issuer" || rate != 1.2 {
		t.Fatalf("unexpected info %s %s %f", sym, issuer, rate)
	}
	if err := tok.Redeem("alice", 40); err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if bal := tok.BalanceOf("alice"); bal != 60 {
		t.Fatalf("unexpected balance %d", bal)
	}
	if err := tok.Redeem("alice", 100); err == nil {
		t.Fatalf("expected redeem error")
	}
}

func TestSYN3700Token(t *testing.T) {
	tok := NewSYN3700Token("Index", "INDX")
	tok.AddComponent("BTC", 0.5)
	tok.AddComponent("ETH", 0.5)
	comps := tok.ListComponents()
	if len(comps) != 2 {
		t.Fatalf("expected 2 components got %d", len(comps))
	}
	val := tok.Value(map[string]float64{"BTC": 10000, "ETH": 2000})
	if val != 6000 {
		t.Fatalf("unexpected value %f", val)
	}
	if err := tok.RemoveComponent("BTC"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if err := tok.RemoveComponent("BTC"); err == nil {
		t.Fatalf("expected remove error")
	}
}

func TestSYN4200Token(t *testing.T) {
	tok := NewSYN4200Token()
	tok.Donate("HELP", "alice", 100, "Food")
	tok.Donate("HELP", "bob", 50, "Food")
	raised, ok := tok.CampaignProgress("HELP")
	if !ok || raised != 150 {
		t.Fatalf("unexpected progress %d", raised)
	}
	camp, ok := tok.Campaign("HELP")
	if !ok || camp.Raised != 150 || camp.Donations["alice"] != 100 {
		t.Fatalf("unexpected campaign %+v", camp)
	}
}

func TestSYN4700LegalToken(t *testing.T) {
	expiry := time.Now().Add(24 * time.Hour)
	tok := NewLegalToken("L1", "Doc", "DOC", "contract", "hash", "alice", expiry, 100, []string{"alice", "bob"})
	if err := tok.Sign("alice", "sigA"); err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := tok.Sign("charlie", "sigC"); err == nil {
		t.Fatalf("expected unknown party error")
	}
	if err := tok.Sign("bob", "sigB"); err != nil {
		t.Fatalf("sign bob: %v", err)
	}
	tok.RevokeSignature("bob")
	tok.UpdateStatus(LegalTokenStatusActive)
	tok.Dispute("breach", "resolved")
	reg := NewLegalTokenRegistry()
	reg.Add(tok)
	got, ok := reg.Get("L1")
	if !ok || got.Status != LegalTokenStatusDisputed || len(got.Signatures) != 1 {
		t.Fatalf("unexpected token %+v", got)
	}
	list := reg.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 token got %d", len(list))
	}
	reg.Remove("L1")
	if _, ok := reg.Get("L1"); ok {
		t.Fatalf("token not removed")
	}
}

func TestSYN70Token(t *testing.T) {
	tok := NewSYN70Token(1, "Game", "GM", 0)
	if err := tok.RegisterAsset("A1", "alice", "Sword", "GameX"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if bal := tok.BalanceOf("alice"); bal != 1 {
		t.Fatalf("unexpected balance %d", bal)
	}
	if err := tok.TransferAsset("A1", "bob"); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if bal := tok.BalanceOf("bob"); bal != 1 {
		t.Fatalf("bob balance %d", bal)
	}
	if err := tok.SetAttribute("A1", "atk", "10"); err != nil {
		t.Fatalf("set attr: %v", err)
	}
	if err := tok.AddAchievement("A1", "First"); err != nil {
		t.Fatalf("add achievement: %v", err)
	}
	info, err := tok.AssetInfo("A1")
	if err != nil || info.Owner != "bob" || info.Attributes["atk"] != "10" || len(info.Achievements) != 1 {
		t.Fatalf("unexpected asset info %+v err %v", info, err)
	}
	assets := tok.ListAssets()
	if len(assets) != 1 {
		t.Fatalf("expected 1 asset got %d", len(assets))
	}
}
