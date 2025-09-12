package synnergy

import (
        "fmt"
        "sync"
        "testing"
)

func TestAccessController(t *testing.T) {
	ac := NewAccessController()
	ac.Grant("admin", "addr1")
	if !ac.HasRole("admin", "addr1") {
		t.Fatalf("expected role granted")
	}
	roles := ac.List("addr1")
	if len(roles) != 1 || roles[0] != "admin" {
		t.Fatalf("unexpected roles: %v", roles)
	}
	ac.Revoke("admin", "addr1")
	if ac.HasRole("admin", "addr1") {
		t.Fatalf("role should be revoked")
	}
}

func TestAccessControllerAuditConcurrent(t *testing.T) {
        ac := NewAccessController()
        var wg sync.WaitGroup
        for i := 0; i < 10; i++ {
                wg.Add(1)
                go func(id int) {
                        defer wg.Done()
                        addr := fmt.Sprintf("addr%d", id)
                        for j := 0; j < 100; j++ {
                                ac.Grant("user", addr)
                                _ = ac.Audit()
                                ac.Revoke("user", addr)
                        }
                }(i)
        }
        wg.Wait()
        if snap := ac.Audit(); len(snap) != 0 {
                t.Fatalf("expected empty snapshot, got %v", snap)
        }
}
