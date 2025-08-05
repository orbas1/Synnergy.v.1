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
	tok := reg.Issue("AssetA", "alice", 100, expiry)
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
}

func TestSYN2800LifePolicyRegistry(t *testing.T) {
	reg := NewLifePolicyRegistry()
	start := time.Now()
	end := start.Add(365 * 24 * time.Hour)
	policy := reg.IssuePolicy("alice", "bob", 1000, 100, start, end)
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
}

func TestSYN2900InsuranceRegistry(t *testing.T) {
	reg := NewInsuranceRegistry()
	start := time.Now()
	end := start.Add(365 * 24 * time.Hour)
	policy := reg.IssuePolicy("alice", "fire", 100, 1000, 10, 1000, start, end)
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
}

func TestSYN3400ForexRegistry(t *testing.T) {
	reg := NewForexRegistry()
	pair := reg.Register("USD", "EUR", 1.1)
	if err := reg.UpdateRate(pair.PairID, 1.2); err != nil {
		t.Fatalf("update rate: %v", err)
	}
	got, ok := reg.Get(pair.PairID)
	if !ok || got.Rate != 1.2 {
		t.Fatalf("unexpected pair data: %+v", got)
	}
	pairs := reg.List()
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair got %d", len(pairs))
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
