package p2p

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

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
	suite             noise.CipherSuite
	static            noise.DHKey
	mu                sync.RWMutex
	allowed           map[string]struct{}
	handshakeTimeout  time.Duration
	identityValidator func([]byte) error
}

// NewNoiseTransport creates a Noise-based transport with a random key pair.
func NewNoiseTransport() (*NoiseTransport, error) {
	suite := noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b)
	static, err := suite.GenerateKeypair(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &NoiseTransport{
		suite:            suite,
		static:           static,
		allowed:          make(map[string]struct{}),
		handshakeTimeout: 10 * time.Second,
	}, nil
}

// StaticPublicKey returns the long-term static key used for Noise handshakes.
func (t *NoiseTransport) StaticPublicKey() []byte {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return append([]byte(nil), t.static.Public...)
}

// AllowPeer whitelists a remote static key. When the allow list is empty all
// peers are accepted.
func (t *NoiseTransport) AllowPeer(pub []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.allowed[string(pub)] = struct{}{}
}

// DisallowPeer removes a remote key from the allowlist.
func (t *NoiseTransport) DisallowPeer(pub []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.allowed, string(pub))
}

// SetIdentityValidator registers a callback that inspects the remote static key.
func (t *NoiseTransport) SetIdentityValidator(fn func([]byte) error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.identityValidator = fn
}

// Dial connects to a remote peer and performs a Noise XX handshake.
func (t *NoiseTransport) Dial(ctx context.Context, addr string) (net.Conn, error) {
	d := &net.Dialer{}
	if deadline, ok := ctx.Deadline(); ok {
		d.Deadline = deadline
	} else {
		// enforce handshake timeout when caller does not provide one
		d.Timeout = t.handshakeTimeout
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
	if tcp, ok := l.Listener.(*net.TCPListener); ok {
		_ = tcp.SetDeadline(time.Now().Add(l.transport.handshakeTimeout))
	}
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
	enc          *noise.CipherState
	dec          *noise.CipherState
	remoteStatic []byte
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

// RemoteStatic returns the remote static key associated with the connection.
func (c *NoiseConn) RemoteStatic() []byte {
	return append([]byte(nil), c.remoteStatic...)
}

// handshake performs a Noise XX handshake on the given connection.
func (t *NoiseTransport) handshake(conn net.Conn, initiator bool) (net.Conn, error) {
	t.mu.RLock()
	cfg := noise.Config{
		CipherSuite:   t.suite,
		Pattern:       noise.HandshakeXX,
		Initiator:     initiator,
		StaticKeypair: t.static,
		Prologue:      []byte("synnergy"),
	}
	allowed := len(t.allowed)
	validator := t.identityValidator
	t.mu.RUnlock()

	hs, err := noise.NewHandshakeState(cfg)
	if err != nil {
		conn.Close()
		return nil, err
	}
	buf := make([]byte, 512)
	var remoteStatic []byte
	if initiator {
		msg, _, _, err := hs.WriteMessage(nil, nil)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if err := writeAll(conn, msg); err != nil {
			conn.Close()
			return nil, err
		}
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if _, _, _, err = hs.ReadMessage(nil, buf[:n]); err != nil {
			conn.Close()
			return nil, err
		}
		msg, tx, rx, err := hs.WriteMessage(nil, nil)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if err := writeAll(conn, msg); err != nil {
			conn.Close()
			return nil, err
		}
		remoteStatic = append([]byte(nil), hs.PeerStatic()...)
		if err := t.validateRemote(remoteStatic, allowed, validator); err != nil {
			conn.Close()
			return nil, err
		}
		return &NoiseConn{Conn: conn, enc: tx, dec: rx, remoteStatic: remoteStatic}, nil
	}
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
	if err := writeAll(conn, msg); err != nil {
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
		remoteStatic = append([]byte(nil), hs.PeerStatic()...)
		if err := t.validateRemote(remoteStatic, allowed, validator); err != nil {
			conn.Close()
			return nil, err
		}
		return &NoiseConn{Conn: conn, enc: tx, dec: rx, remoteStatic: remoteStatic}, nil
	}
}

func (t *NoiseTransport) validateRemote(remote []byte, allowCount int, validator func([]byte) error) error {
	if allowCount > 0 {
		t.mu.RLock()
		_, ok := t.allowed[string(remote)]
		t.mu.RUnlock()
		if !ok {
			return errors.New("p2p: remote key not authorised")
		}
	}
	if validator != nil {
		if err := validator(remote); err != nil {
			return err
		}
	}
	return nil
}

func writeAll(conn net.Conn, msg []byte) error {
	for len(msg) > 0 {
		n, err := conn.Write(msg)
		if err != nil {
			return err
		}
		msg = msg[n:]
	}
	return nil
}
