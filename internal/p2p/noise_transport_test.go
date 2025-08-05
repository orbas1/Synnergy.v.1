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
	serverConn, clientConn := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		_, err := server.handshake(serverConn, false)
		errCh <- err
	}()
	if _, err := client.handshake(clientConn, true); err != nil {
		t.Fatalf("client handshake: %v", err)
	}
	if err := <-errCh; err != nil {
		t.Fatalf("server handshake: %v", err)
	}
}
