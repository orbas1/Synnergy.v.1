package p2p

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"

	"github.com/flynn/noise"
)

// Transport defines the basic operations for a network transport.
type Transport interface {
	Dial(ctx context.Context, addr string) (net.Conn, error)
	Listen(ctx context.Context, addr string) (net.Listener, error)
}

// NoiseTransport implements the Noise protocol (XX handshake) for securing
// peer-to-peer connections.
type NoiseTransport struct {
	suite  noise.CipherSuite
	static noise.DHKey
}

// NewNoiseTransport creates a Noise-based transport with a random key pair.
func NewNoiseTransport() (*NoiseTransport, error) {
	suite := noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b)
	static, err := suite.GenerateKeypair(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &NoiseTransport{suite: suite, static: static}, nil
}

// Dial connects to a remote peer and performs a Noise XX handshake.
func (t *NoiseTransport) Dial(ctx context.Context, addr string) (net.Conn, error) {
	d := &net.Dialer{}
	if deadline, ok := ctx.Deadline(); ok {
		d.Deadline = deadline
	}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	return t.handshake(conn, true)
}

// Listen starts a TCP listener that performs Noise handshakes for each
// incoming connection.
func (t *NoiseTransport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &noiseListener{Listener: l, transport: t}, nil
}

// noiseListener wraps a net.Listener to perform Noise handshakes on Accept.
type noiseListener struct {
	net.Listener
	transport *NoiseTransport
}

// Accept waits for and returns the next connection, completing the Noise
// handshake as a responder.
func (l *noiseListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return l.transport.handshake(c, false)
}

// NoiseConn wraps a net.Conn and encrypts all traffic using Noise cipher
// states.
type NoiseConn struct {
	net.Conn
	enc *noise.CipherState
	dec *noise.CipherState
}

func (c *NoiseConn) Write(p []byte) (int, error) {
	msg, err := c.enc.Encrypt(nil, nil, p)
	if err != nil {
		return 0, err
	}
	var lenBuf [2]byte
	binary.BigEndian.PutUint16(lenBuf[:], uint16(len(msg)))
	if _, err := c.Conn.Write(lenBuf[:]); err != nil {
		return 0, err
	}
	if _, err := c.Conn.Write(msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *NoiseConn) Read(p []byte) (int, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(c.Conn, lenBuf[:]); err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint16(lenBuf[:])
	buf := make([]byte, n)
	if _, err := io.ReadFull(c.Conn, buf); err != nil {
		return 0, err
	}
	out, err := c.dec.Decrypt(nil, nil, buf)
	if err != nil {
		return 0, err
	}
	copy(p, out)
	return len(out), nil
}

// handshake performs a Noise XX handshake on the given connection.
func (t *NoiseTransport) handshake(conn net.Conn, initiator bool) (net.Conn, error) {
	cfg := noise.Config{
		CipherSuite:   t.suite,
		Pattern:       noise.HandshakeXX,
		Initiator:     initiator,
		StaticKeypair: t.static,
		Prologue:      []byte("synnergy"),
	}
	hs, err := noise.NewHandshakeState(cfg)
	if err != nil {
		conn.Close()
		return nil, err
	}
	buf := make([]byte, 512)
	if initiator {
		// -> e
		msg, _, _, err := hs.WriteMessage(nil, nil)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if _, err = conn.Write(msg); err != nil {
			conn.Close()
			return nil, err
		}
		// <- e, ee, s, es
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if _, _, _, err = hs.ReadMessage(nil, buf[:n]); err != nil {
			conn.Close()
			return nil, err
		}
		// -> s, se
		msg, tx, rx, err := hs.WriteMessage(nil, nil)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if _, err = conn.Write(msg); err != nil {
			conn.Close()
			return nil, err
		}
		return &NoiseConn{Conn: conn, enc: tx, dec: rx}, nil
	}
	// Responder path
	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, _, _, err = hs.ReadMessage(nil, buf[:n]); err != nil {
		conn.Close()
		return nil, err
	}
	msg, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, err = conn.Write(msg); err != nil {
		conn.Close()
		return nil, err
	}
	n, err = conn.Read(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, rx, tx, err := hs.ReadMessage(nil, buf[:n]); err != nil {
		conn.Close()
		return nil, err
	} else {
		return &NoiseConn{Conn: conn, enc: tx, dec: rx}, nil
	}
}
