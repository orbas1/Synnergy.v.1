package cli

import (
        "net"
        "strings"
        "testing"
)

// TestConnPoolLifecycle ensures the pool can dial, release and close connections.
func TestConnPoolLifecycle(t *testing.T) {
        ln, err := net.Listen("tcp", "127.0.0.1:0")
        if err != nil {
                t.Fatalf("listen: %v", err)
        }
        defer ln.Close()
        go func() {
                conn, err := ln.Accept()
                if err == nil {
                        conn.Close()
                }
        }()
        addr := ln.Addr().String()
        if _, err := execCommand("connpool", "dial", addr); err != nil {
                t.Fatalf("dial failed: %v", err)
        }
	out, err := execCommand("connpool", "stats")
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if !strings.Contains(out, "active: 1") {
		t.Fatalf("expected active 1, got %q", out)
	}
        if _, err := execCommand("connpool", "release", addr); err != nil {
                t.Fatalf("release failed: %v", err)
        }
	out, err = execCommand("connpool", "stats")
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if !strings.Contains(out, "active: 0") {
		t.Fatalf("expected active 0, got %q", out)
	}
	if _, err := execCommand("connpool", "close"); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}
