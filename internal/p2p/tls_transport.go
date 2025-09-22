package p2p

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"sync"
	"time"
)

// TLSTransport implements the Transport interface using TLS for security.
type TLSTransport struct {
	mu               sync.RWMutex
	config           *tls.Config
	handshakeTimeout time.Duration
	allowedSPKI      [][]byte
}

// NewTLSTransport builds a TLS transport. The same certificate may be used for
// both client and server in test environments.
func NewTLSTransport(cert tls.Certificate, caPool *x509.CertPool, server bool) *TLSTransport {
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}
	if server {
		cfg.ClientCAs = caPool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		cfg.RootCAs = caPool
	}
	t := &TLSTransport{config: cfg, handshakeTimeout: 10 * time.Second}
	t.configureVerifier()
	return t
}

// SetAllowedSPKI restricts connections to peers whose leaf certificates match
// one of the provided Subject Public Key Info hashes.
func (t *TLSTransport) SetAllowedSPKI(fingerprints [][]byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.allowedSPKI = make([][]byte, len(fingerprints))
	for i, fp := range fingerprints {
		t.allowedSPKI[i] = append([]byte(nil), fp...)
	}
}

// ReloadCertificates hot swaps the certificate chain.
func (t *TLSTransport) ReloadCertificates(cert tls.Certificate) {
	t.mu.Lock()
	t.config.Certificates = []tls.Certificate{cert}
	t.mu.Unlock()
}

// Dial connects to a remote peer over TLS.
func (t *TLSTransport) Dial(ctx context.Context, addr string) (net.Conn, error) {
	t.mu.RLock()
	cfg := t.config.Clone()
	timeout := t.handshakeTimeout
	t.mu.RUnlock()
	d := &net.Dialer{}
	if deadline, ok := ctx.Deadline(); ok {
		d.Deadline = deadline
	} else {
		d.Timeout = timeout
	}
	return tls.DialWithDialer(d, "tcp", addr, cfg)
}

// Listen starts a TLS listener on the given address.
func (t *TLSTransport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	t.mu.RLock()
	cfg := t.config.Clone()
	timeout := t.handshakeTimeout
	t.mu.RUnlock()
	ln, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}
	return &tlsListener{Listener: ln, timeout: timeout}, nil
}

type tlsListener struct {
	net.Listener
	timeout time.Duration
}

func (l *tlsListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if l.timeout > 0 {
		_ = conn.SetDeadline(time.Now().Add(l.timeout))
	}
	return conn, nil
}

func (t *TLSTransport) configureVerifier() {
	t.mu.Lock()
	cfg := t.config
	t.mu.Unlock()
	cfg.VerifyPeerCertificate = func(raw [][]byte, chains [][]*x509.Certificate) error {
		if len(raw) == 0 {
			return errors.New("tls: missing peer certificate")
		}
		if len(t.allowedSPKI) == 0 {
			return nil
		}
		cert, err := x509.ParseCertificate(raw[0])
		if err != nil {
			return err
		}
		fp := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
		t.mu.RLock()
		defer t.mu.RUnlock()
		for _, allowed := range t.allowedSPKI {
			if len(allowed) == len(fp) && hmacEqual(allowed, fp[:]) {
				return nil
			}
		}
		return errors.New("tls: peer fingerprint not authorised")
	}
}

func hmacEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := range a {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
