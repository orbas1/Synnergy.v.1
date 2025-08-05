package synnergy

import "testing"

func TestIdentityService(t *testing.T) {
	svc := NewIdentityService()
	if err := svc.Register("addr1", "Alice", "2000-01-01", "US"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.Verify("addr1", "passport"); err != nil {
		t.Fatalf("verify: %v", err)
	}
	info, ok := svc.Info("addr1")
	if !ok || info.Name != "Alice" {
		t.Fatalf("unexpected info: %v %v", info, ok)
	}
	logs := svc.Logs("addr1")
	if len(logs) != 1 || logs[0].Method != "passport" {
		t.Fatalf("unexpected logs: %v", logs)
	}
}
