package p2p

import (
	"net"
	"testing"
)

// TestNoiseHandshake ensures that the Noise XX handshake completes between two peers.
func TestNoiseHandshake(t *testing.T) {
	server, err := NewNoiseTransport()
	if err != nil {
		t.Fatalf("server transport: %v", err)
	}
	client, err := NewNoiseTransport()
	if err != nil {
		t.Fatalf("client transport: %v", err)
	}
	server.AllowPeer(client.StaticPublicKey())
	client.AllowPeer(server.StaticPublicKey())
	serverConn, clientConn := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		conn, err := server.handshake(serverConn, false)
		if err == nil {
			defer conn.Close()
		}
		errCh <- err
	}()
	conn, err := client.handshake(clientConn, true)
	if err != nil {
		t.Fatalf("client handshake: %v", err)
	}
	noiseConn := conn.(*NoiseConn)
	if len(noiseConn.RemoteStatic()) == 0 {
		t.Fatalf("expected remote static key")
	}
	if err := <-errCh; err != nil {
		t.Fatalf("server handshake: %v", err)
	}
}

func TestNoiseHandshakeUnauthorized(t *testing.T) {
	server, err := NewNoiseTransport()
	if err != nil {
		t.Fatalf("server transport: %v", err)
	}
	client, err := NewNoiseTransport()
	if err != nil {
		t.Fatalf("client transport: %v", err)
	}
	// allow list mismatched key to trigger rejection
	server.AllowPeer(client.StaticPublicKey()[:16])
	serverConn, clientConn := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		_, err := server.handshake(serverConn, false)
		errCh <- err
	}()
	conn, err := client.handshake(clientConn, true)
	if err == nil {
		defer conn.Close()
		if _, err := conn.Write([]byte("test")); err == nil {
			t.Fatalf("expected write failure after server rejection")
		}
	}
	if err := <-errCh; err == nil {
		t.Fatalf("expected server failure")
	}
}
