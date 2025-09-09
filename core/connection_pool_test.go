package core

import (
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

	accepted := make(chan net.Conn, 1)
	go func() {
		c, err := ln.Accept()
		if err == nil {
			accepted <- c
		}
	}()

	addr := ln.Addr().String()
	p := NewConnectionPool(1)
	c1, err := p.Acquire(addr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c1.Conn == nil {
		t.Fatalf("expected net.Conn")
	}
	if _, err := p.Acquire("other:0"); err == nil {
		t.Fatalf("expected pool exhaustion error")
	}
	if p.Size() != 1 {
		t.Fatalf("expected size 1, got %d", p.Size())
	}

	p.Release(addr)
	if p.Size() != 0 {
		t.Fatalf("release failed")
	}
	serverConn := <-accepted
	defer serverConn.Close()
	serverConn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	buf := make([]byte, 1)
	if _, err := serverConn.Read(buf); err != io.EOF {
		t.Fatalf("expected EOF after release, got %v", err)
	}

	if stats := p.Stats(); stats.Capacity != 1 || stats.Active != 0 {
		t.Fatalf("unexpected stats %+v", stats)
	}
	p.Close()
}
