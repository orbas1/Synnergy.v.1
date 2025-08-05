package p2p

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
)

// TLSTransport implements the Transport interface using TLS for security.
type TLSTransport struct {
	config *tls.Config
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
		cfg.InsecureSkipVerify = false
	}
	return &TLSTransport{config: cfg}
}

// Dial connects to a remote peer over TLS.
func (t *TLSTransport) Dial(ctx context.Context, addr string) (net.Conn, error) {
	d := &net.Dialer{}
	if deadline, ok := ctx.Deadline(); ok {
		d.Deadline = deadline
	}
	return tls.DialWithDialer(d, "tcp", addr, t.config)
}

// Listen starts a TLS listener on the given address.
func (t *TLSTransport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	return tls.Listen("tcp", addr, t.config)
}
