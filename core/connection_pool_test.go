package core

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestConnectionPool(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	accepted := make(chan net.Conn, 4)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			accepted <- c
		}
	}()

	addr := ln.Addr().String()
	pool := NewConnectionPoolWithOptions(PoolOptions{
		Max:                 1,
		IdleTimeout:         10 * time.Millisecond,
		HealthCheckInterval: 5 * time.Millisecond,
	})
	defer pool.Close()

	c1, err := pool.Acquire(addr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c1.Conn == nil {
		t.Fatalf("expected net.Conn")
	}

	if other, err := pool.Acquire("other:0"); err == nil || other != nil {
		t.Fatalf("expected pool exhaustion error")
	}
	if pool.Size() != 1 {
		t.Fatalf("expected size 1, got %d", pool.Size())
	}

	if reused, err := pool.Acquire(addr); err != nil || reused != c1 {
		t.Fatalf("expected connection reuse, err=%v", err)
	}

	serverConn1 := <-accepted
	pool.ReportFailure(addr, errors.New("temporary"))

	c2, err := pool.Acquire(addr)
	if err != nil {
		t.Fatalf("reacquire after failure: %v", err)
	}
	if c2 == c1 {
		t.Fatalf("expected new connection after failure")
	}
	serverConn2 := <-accepted

	serverConn1.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	buf := make([]byte, 1)
	if _, err := serverConn1.Read(buf); err != io.EOF {
		t.Fatalf("expected EOF after replacement, got %v", err)
	}
	_ = serverConn1.Close()

	time.Sleep(30 * time.Millisecond)
	if pool.Size() != 0 {
		t.Fatalf("expected idle pruning to clear pool")
	}

	stats := pool.Stats()
	if stats.Capacity != 1 || stats.Created == 0 || stats.Reused == 0 {
		t.Fatalf("unexpected stats %+v", stats)
	}
	if stats.ClosedIdle == 0 {
		t.Fatalf("expected idle closures to be recorded")
	}

	_ = serverConn2.Close()
}
